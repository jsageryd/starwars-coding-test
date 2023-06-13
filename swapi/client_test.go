package swapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jsageryd/starwars-coding-test/starwars"
)

func TestClient_People(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var nextURL string

		ts := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				switch path := r.URL.Path; path {
				case "/people/":
					w.Write([]byte(`
{
  "next": "` + nextURL + `",
  "results": [
    {"name":"Luke Skywalker","height":"172","mass":"77","birth_year":"19BBY"},
    {"name":"R2-D2","height":"96","mass":"32","birth_year":"33BBY"}
  ]
}
`))
				case "/foo-next-page":
					w.Write([]byte(`
{
  "results": [
    {"name":"C-3PO","height":"167","mass":"75","birth_year":"112BBY"}
  ]
}
`))
				default:
					t.Fatalf("unexpected path: %s", path)
				}
			},
		))

		nextURL = ts.URL + "/foo-next-page"

		c := NewClient(ts.URL)

		gotCharacters, err := c.People()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		wantCharacters := []starwars.Character{
			{Name: "Luke Skywalker", Height: 172, Mass: 77, BirthYear: "19BBY"},
			{Name: "R2-D2", Height: 96, Mass: 32, BirthYear: "33BBY"},
			{Name: "C-3PO", Height: 167, Mass: 75, BirthYear: "112BBY"},
		}

		if fmt.Sprint(gotCharacters) != fmt.Sprint(wantCharacters) {
			t.Errorf("got %v, want %v", gotCharacters, wantCharacters)
		}
	})

	t.Run("Caching", func(t *testing.T) {
		var gotReqCount int

		ts := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				gotReqCount++

				w.Write([]byte(`
{
  "results": [
    {"name":"C-3PO","height":"167","mass":"75","birth_year":"112BBY"}
  ]
}
`))
			},
		))

		c := NewClient(ts.URL)

		var gotCharacters []starwars.Character
		var err error

		for n := 0; n < 2; n++ {
			gotCharacters, err = c.People()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}

		if got, want := gotReqCount, 1; got != want {
			t.Errorf("set %d requests, want %d", got, want)
		}

		wantCharacters := []starwars.Character{
			{Name: "C-3PO", Height: 167, Mass: 75, BirthYear: "112BBY"},
		}

		if fmt.Sprint(gotCharacters) != fmt.Sprint(wantCharacters) {
			t.Errorf("got %v, want %v", gotCharacters, wantCharacters)
		}
	})
}
