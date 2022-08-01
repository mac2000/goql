# goql

small utility tries to mimique curl for graphql services

instead of:

```bash
curl -X POST http://localhost/graphql -H "Content-Type: application/json" -d '{ "query": "{ user(id: \"1\") { id name } }" }'
```

we can:

```bash
goql --url=http://localhost/graphql --query='{ user(id: "1") { id name } }'
```

