package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jsageryd/starwars-coding-test/mock"
	"github.com/jsageryd/starwars-coding-test/starwars"
)

func TestAPI_TopFatCharacters(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		a := New(
			&mock.Core{
				TopFatCharactersFunc: func() ([]starwars.Character, error) {
					return []starwars.Character{
						{Name: "R2-D2", Height: 96, Mass: 32},
						{Name: "C-3PO", Height: 167, Mass: 75},
						{Name: "Luke Skywalker", Height: 172, Mass: 77},
					}, nil
				},
			},
		)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		a.topFatCharacters(w, r)

		if got, want := w.Code, http.StatusOK; got != want {
			t.Errorf("got HTTP %d, want %d", got, want)
		}

		wantBody := `[{"name":"R2-D2","height":"96","mass":"32"},{"name":"C-3PO","height":"167","mass":"75"},{"name":"Luke Skywalker","height":"172","mass":"77"}]` + "\n"

		if got, want := w.Body.String(), wantBody; got != want {
			t.Errorf("got body:\n%s\nwant:\n%s", got, want)
		}
	})

	t.Run("Error from core", func(t *testing.T) {
		a := New(
			&mock.Core{
				TopFatCharactersFunc: func() ([]starwars.Character, error) {
					return nil, errors.New("foo error")
				},
			},
		)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		a.topFatCharacters(w, r)

		if got, want := w.Code, http.StatusInternalServerError; got != want {
			t.Errorf("got HTTP %d, want %d", got, want)
		}
	})

	t.Run("Wrong method", func(t *testing.T) {
		a := New(nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/", nil)

		a.topFatCharacters(w, r)

		if got, want := w.Code, http.StatusMethodNotAllowed; got != want {
			t.Errorf("got HTTP %d, want %d", got, want)
		}
	})
}

func TestAPI_TopOldCharacters(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		a := New(
			&mock.Core{
				TopOldCharactersFunc: func() ([]starwars.Character, error) {
					return []starwars.Character{
						{Name: "Darth Vader", BirthYear: "41.9BBY"},
						{Name: "Luke Skywalker", BirthYear: "19BBY"},
					}, nil
				},
			},
		)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		a.topOldCharacters(w, r)

		if got, want := w.Code, http.StatusOK; got != want {
			t.Errorf("got HTTP %d, want %d", got, want)
		}

		wantBody := `[{"name":"Darth Vader","birth_year":"41.9BBY"},{"name":"Luke Skywalker","birth_year":"19BBY"}]` + "\n"

		if got, want := w.Body.String(), wantBody; got != want {
			t.Errorf("got body:\n%s\nwant:\n%s", got, want)
		}
	})

	t.Run("Error from core", func(t *testing.T) {
		a := New(
			&mock.Core{
				TopOldCharactersFunc: func() ([]starwars.Character, error) {
					return nil, errors.New("foo error")
				},
			},
		)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		a.topOldCharacters(w, r)

		if got, want := w.Code, http.StatusInternalServerError; got != want {
			t.Errorf("got HTTP %d, want %d", got, want)
		}
	})

	t.Run("Wrong method", func(t *testing.T) {
		a := New(nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/", nil)

		a.topOldCharacters(w, r)

		if got, want := w.Code, http.StatusMethodNotAllowed; got != want {
			t.Errorf("got HTTP %d, want %d", got, want)
		}
	})
}
