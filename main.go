package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alecthomas/kong"
)

type GraphQLRequest struct {
	Query string `json:"query"`
}

type QueryCmd struct {
	Fail    bool              `short:"f" help:"Fail silently (no output at all) on HTTP errors"`
	Include bool              `short:"i" help:"Include protocol response headers in the output"`
	URL     string            `default:"http://localhost/graphql" required:"" help:"URL to send query to"`
	Query   string            `required:"" help:"GraphQL query"`
	Headers map[string]string `short:"H" help:"Additional headers to pass"`
}

func (r *QueryCmd) Run() error {
	payload, err := json.Marshal(GraphQLRequest{
		Query: r.Query,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", r.URL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if r.Include {
		fmt.Printf("%s %d\n", resp.Proto, resp.StatusCode)
		for k, a := range resp.Header {
			for _, v := range a {
				fmt.Printf("%s: %s\n", k, v)
			}
		}
		fmt.Println()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	if r.Fail && resp.StatusCode != 200 {
		return fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	_, hasErrors := data["errors"]
	if hasErrors {
		return fmt.Errorf("response has errors")
	}
	return nil
}

var CLI struct {
	Query QueryCmd `cmd:"" default:"withargs" help:"Query GraphQL Service"`
}

func main() {
	ctx := kong.Parse(&CLI,
		kong.Name("goql"),
		kong.Description("GraphQL CLI Client"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": "0.0.1",
		},
	)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
