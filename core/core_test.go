package core

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jsageryd/starwars-coding-test/starwars"
)

func TestCore_TopFatCharacters(t *testing.T) {
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
    {"name":"Luke Skywalker","height":"172","mass":"77"},
    {"name":"R2-D2","height":"96","mass":"32"}
  ]
}
`))
				case "/foo-next-page":
					w.Write([]byte(`
{
  "results": [
    {"name":"C-3PO","height":"167","mass":"75"}
  ]
}
`))
				default:
					t.Fatalf("unexpected path: %s", path)
				}
			},
		))

		nextURL = ts.URL + "/foo-next-page"

		c := New(ts.URL)

		gotCharacters, err := c.TopFatCharacters()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		wantCharacters := []starwars.Character{
			{Name: "R2-D2", Height: 96, Mass: 32},
			{Name: "C-3PO", Height: 167, Mass: 75},
			{Name: "Luke Skywalker", Height: 172, Mass: 77},
		}

		if fmt.Sprint(gotCharacters) != fmt.Sprint(wantCharacters) {
			t.Errorf("got %v, want %v", gotCharacters, wantCharacters)
		}
	})
}

func TestTopFat(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		input := []starwars.Character{
			{Name: "Burly Bob", Height: 1.65, Mass: 118},         // BMI 24
			{Name: "Hearty Hank", Height: 1.70, Mass: 102},       // BMI 22
			{Name: "Middleweight Mitch", Height: 1.75, Mass: 88}, // BMI 20
			{Name: "Plump Paul", Height: 1.68, Mass: 110},        // BMI 23
			{Name: "Stocky Steve", Height: 1.63, Mass: 127},      // BMI 25
			{Name: "Sturdy Stan", Height: 1.73, Mass: 95},        // BMI 21
		}

		gotOutput := topFat(input, 3)

		wantOutput := []starwars.Character{
			{Name: "Stocky Steve", Height: 1.63, Mass: 127}, // BMI 25
			{Name: "Burly Bob", Height: 1.65, Mass: 118},    // BMI 24
			{Name: "Plump Paul", Height: 1.68, Mass: 110},   // BMI 23
		}

		if fmt.Sprint(gotOutput) != fmt.Sprint(wantOutput) {
			t.Errorf("got %v, want %v", gotOutput, wantOutput)
		}
	})

	t.Run("Character count less than limit", func(t *testing.T) {
		input := []starwars.Character{
			{Name: "Burly Bob", Height: 1.65, Mass: 118}, // BMI 24
		}

		gotOutput := topFat(input, 3)

		wantOutput := []starwars.Character{
			{Name: "Burly Bob", Height: 1.65, Mass: 118}, // BMI 24
		}

		if fmt.Sprint(gotOutput) != fmt.Sprint(wantOutput) {
			t.Errorf("got %v, want %v", gotOutput, wantOutput)
		}
	})
}
