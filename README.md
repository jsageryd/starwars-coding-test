# starwars-coding-test

## Usage
Install [Go](https://go.dev) if you haven't.

```
$ brew install go
```

### Tests
Use `go test` to run tests.

```
$ go test -v ./...
```

### Run
Compile and run using for example `go run`.

```
$ go run main.go
2023/06/13 11:16:05 Listening at :8080...
```

Run the UI in a browser...

```
$ open http://localhost:8080/
```

...or send an API request (example uses [HTTPie](https://httpie.io/)).

```
$ http get :8080/top-fat-characters
HTTP/1.1 200 OK
Content-Length: 989
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Jun 2023 09:17:07 GMT

[
    {
        "birth_year": "52BBY",
        "height": "178",
        "mass": "120",
        "name": "Owen Lars"
    },
    {
        "birth_year": "15BBY",
        "height": "200",
        "mass": "140",
        "name": "IG-88"
    },
[...]
```

```
$ http get :8080/top-old-characters
HTTP/1.1 200 OK
Content-Length: 1284
Content-Type: application/json
Date: Tue, 13 Jun 2023 13:57:38 GMT

[
    {
        "birth_year": "112BBY",
        "height": "167",
        "mass": "75",
        "name": "C-3PO"
    },
    {
        "birth_year": "92BBY",
        "height": "198",
        "mass": "82",
        "name": "Ki-Adi-Mundi"
    },
[...]
```
