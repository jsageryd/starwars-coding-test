# starwars-coding-test

## Usage
Install [Go](https://go.dev) if you haven't.

```
$ brew install go
```

Run the server.

```
$ go run main.go
2023/06/13 11:16:05 Listening at :8080...
```

Send a request (example uses [HTTPie](https://httpie.io/)).

```
$ http get :8080/top-fat-characters
HTTP/1.1 200 OK
Content-Length: 989
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Jun 2023 09:17:07 GMT

[
    {
        "height": "178",
        "mass": "120",
        "name": "Owen Lars"
    },
    {
        "height": "200",
        "mass": "140",
        "name": "IG-88"
    },
[...]
```
