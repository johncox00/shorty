# shorty
Go + MongoDB exercise

Run it:
`$ go run *.go`

You're now running on localhost:5100.

Then:

`POST /shorturls` with

```
{ "longurl" : "http://somereallylongurl/with/a/path?and=query&string=params"}
```

to get

```
{ "longurl" : "http://somereallylongurl/with/a/path?and=query&string=params", "shorturl" : "J7ger30k"}
```

so that you can

`GET /J7ger30k`

and get the return

```
{ "longurl" : "http://somereallylongurl/with/a/path?and=query&string=params", "shorturl" : "J7ger30k"}
```

You can also `GET /J7ger30k/click_count` to get the number of times somebody has done a `GET` on your shorturl.
