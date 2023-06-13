package core

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jsageryd/starwars-coding-test/starwars"
	"github.com/jsageryd/starwars-coding-test/swapi"
)

func TestCore_TopFatCharacters(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`
{
  "results": [
    {"name":"Luke Skywalker","height":"172","mass":"77"},
    {"name":"R2-D2","height":"96","mass":"32"},
    {"name":"C-3PO","height":"167","mass":"75"}
  ]
}
`))
			},
		))

		c := New(
			swapi.NewClient(ts.URL),
		)

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

	t.Run("Non-OK response from API", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusTeapot)
			},
		))

		c := New(
			swapi.NewClient(ts.URL),
		)

		_, err := c.TopFatCharacters()

		if err == nil {
			t.Fatal("error is nil")
		}

		wantErrStr := "error fetching characters from SWAPI: SWAPI returned HTTP 418"

		if got, want := err.Error(), wantErrStr; got != want {
			t.Errorf("err is %q, want %q", got, want)
		}
	})
}

func TestCore_TopOldCharacters(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`
{
  "results": [
    {"name":"Luke Skywalker","birth_year":"19BBY"},
    {"name":"Darth Vader","birth_year":"41.9BBY"},
    {"name":"R5-D4","birth_year":"unknown"}
  ]
}
`))
			},
		))

		c := New(
			swapi.NewClient(ts.URL),
		)

		gotCharacters, err := c.TopOldCharacters()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		wantCharacters := []starwars.Character{
			{Name: "Darth Vader", BirthYear: "41.9BBY"},
			{Name: "Luke Skywalker", BirthYear: "19BBY"},
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

		c := New(
			swapi.NewClient(ts.URL),
		)

		_, err := c.TopOldCharacters()

		if err == nil {
			t.Fatal("error is nil")
		}

		wantErrStr := "error fetching characters from SWAPI: SWAPI returned HTTP 418"

		if got, want := err.Error(), wantErrStr; got != want {
			t.Errorf("err is %q, want %q", got, want)
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

func TestTopOld(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		input := []starwars.Character{
			{Name: "Ben Solo", BirthYear: "5ABY"},
			{Name: "Darth Vader", BirthYear: "41.9BBY"},
			{Name: "Leia Skywalker", BirthYear: "19BBY"},
			{Name: "Luke Skywalker", BirthYear: "19BBY"},
			{Name: "R5-D4", BirthYear: "unknown"},
			{Name: "Rey", BirthYear: "15ABY"},
		}

		gotOutput := topOld(input, 4)

		wantOutput := []starwars.Character{
			{Name: "Darth Vader", BirthYear: "41.9BBY"},
			{Name: "Leia Skywalker", BirthYear: "19BBY"},
			{Name: "Luke Skywalker", BirthYear: "19BBY"},
			{Name: "Ben Solo", BirthYear: "5ABY"},
		}

		if fmt.Sprint(gotOutput) != fmt.Sprint(wantOutput) {
			t.Errorf("got %v, want %v", gotOutput, wantOutput)
		}
	})

	t.Run("Character count less than limit", func(t *testing.T) {
		input := []starwars.Character{
			{Name: "Ben Solo", BirthYear: "5ABY"},
		}

		gotOutput := topOld(input, 3)

		wantOutput := []starwars.Character{
			{Name: "Ben Solo", BirthYear: "5ABY"},
		}

		if fmt.Sprint(gotOutput) != fmt.Sprint(wantOutput) {
			t.Errorf("got %v, want %v", gotOutput, wantOutput)
		}
	})
}

func TestAbsBirthYear(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		for n, tc := range []struct {
			input string
			want  float64
		}{
			{"20BBY", -20},
			{"10.5BBY", -10.5},
			{"1BBY", -1},
			{"0BBY", 0},
			{"0ABY", 0},
			{"1ABY", 1},
			{"10.5ABY", 10.5},
			{"20ABY", 20},
		} {
			gotYear, err := absBirthYear(tc.input)
			if err != nil {
				t.Errorf("[%d] unexpected error: %v", n, err)
				continue
			}

			if gotYear != tc.want {
				t.Errorf("[%d] absBirthYear(%q) = %f, want %f", n, tc.input, gotYear, tc.want)
			}
		}
	})

	t.Run("Unknown year format", func(t *testing.T) {
		for n, input := range []string{
			"",
			"foo",
			"3foo",
			"unknown",
		} {
			_, err := absBirthYear(input)

			if err == nil {
				t.Fatal("error is nil")
			}

			if gotErrStr, wantPrefix := err.Error(), "unknown birth date format"; !strings.HasPrefix(gotErrStr, wantPrefix) {
				t.Errorf("[%d] error is %q, want prefix %q", n, gotErrStr, wantPrefix)
			}
		}
	})
}
