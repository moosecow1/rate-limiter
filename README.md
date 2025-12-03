# Custom Rate Limiter
### Rate Limit
If you send at least X amount of requests in Y time span, the API doesn't complete your requests.

It will tell you how long until you can request again.
```json
{
  "try_in": 0
}
```

### API
`/get-quote` returns a random quote:
```json
{
  "text": "Example quote.",
  "author": "Example Name"
}
```
`/get-quote?Demo Name` returns only quotes by a specific person:
```json
{
  "text": "Example quote.",
  "author": "Demo Name"
}
```
