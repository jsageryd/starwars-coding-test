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
		ts := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				t.Logf("page is %q", r.URL.Query().Get("page"))
				switch page := r.URL.Query().Get("page"); page {
				case "1":
					w.Write([]byte(`
{
  "count": 3,
  "results": [
    {"name":"Luke Skywalker","height":"172","mass":"77","birth_year":"19BBY"},
    {"name":"R2-D2","height":"96","mass":"32","birth_year":"33BBY"}
  ]
}
`))
				case "2":
					w.Write([]byte(`
{
  "results": [
    {"name":"C-3PO","height":"167","mass":"75","birth_year":"112BBY"}
  ]
}
`))
				default:
					t.Fatalf("unexpected page: %s", page)
				}
			},
		))

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

	t.Run("Non-OK response from API", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusTeapot)
			},
		))

		c := NewClient(ts.URL)

		_, err := c.People()

		if err == nil {
			t.Fatal("error is nil")
		}

		if got, want := err.Error(), "SWAPI returned HTTP 418"; got != want {
			t.Errorf("error is %q, want %q", got, want)
		}
	})
}

func TestClient_FetchPeople(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var gotPath string
		var gotRawQuery string

		ts := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				gotRawQuery = r.URL.RawQuery

				w.Write([]byte(`
{
  "count": 2,
  "results": [
    {"name":"Luke Skywalker","height":"172","mass":"77","birth_year":"19BBY"},
    {"name":"R2-D2","height":"96","mass":"32","birth_year":"33BBY"}
  ]
}
`))
			},
		))

		c := NewClient(ts.URL)

		gotResp, err := c.fetchPeople(3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		wantResp := peopleResponse{
			Count: 2,
			Results: []starwars.Character{
				{Name: "Luke Skywalker", Height: 172, Mass: 77, BirthYear: "19BBY"},
				{Name: "R2-D2", Height: 96, Mass: 32, BirthYear: "33BBY"},
			},
		}

		if got, want := gotPath, "/people/"; got != want {
			t.Errorf("got path %q, want %q", got, want)
		}

		if got, want := gotRawQuery, "page=3"; got != want {
			t.Errorf("got raw query %q, want %q", got, want)
		}

		if fmt.Sprint(gotResp) != fmt.Sprint(wantResp) {
			t.Errorf("got %v, want %v", gotResp, wantResp)
		}
	})

	t.Run("Non-OK response from API", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusTeapot)
			},
		))

		c := NewClient(ts.URL)

		_, err := c.fetchPeople(3)

		if err == nil {
			t.Fatal("error is nil")
		}

		if got, want := err.Error(), "SWAPI returned HTTP 418"; got != want {
			t.Errorf("error is %q, want %q", got, want)
		}
	})
}
