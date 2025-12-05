# Custom Rate Limiter
This is a simple REST API made with Gin to learn Go and implement a custom rate limiter.

The demo API has one endpoint that returns a random quote. 

### Run Instructions

1. Install [Docker Desktop](https://www.docker.com/products/docker-desktop/)
2. Open Docker Desktop
3. Run `docker-compose up`

### Rate Limit
If a client sends too many requests in too short of a time span, new requests will be blocked.

The limits can be configured. 

It will tell you how long until you can request again.
```json
{
  "try_in": 5
}
```

### API
`GET /get-quote` returns a random quote:
```json
{
  "text": "Example quote.",
  "author": "Example Name"
}
```
`GET /get-quote?Demo Name` returns only quotes by a specific person:
```json
{
  "text": "Example quote.",
  "author": "Demo Name"
}
```
