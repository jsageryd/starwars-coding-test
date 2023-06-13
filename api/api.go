package api

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"text/template"

	"github.com/jsageryd/starwars-coding-test/starwars"
)

//go:embed index.tmpl
var index string

var tmpl *template.Template

func init() {
	var err error

	tmpl, err = template.New("index").Parse(index)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}
}

type API struct {
	core starwars.Core
}

func New(core starwars.Core) *API {
	return &API{
		core: core,
	}
}

func (a *API) Register(mux *http.ServeMux) {
	mux.HandleFunc("/", a.ui)
	mux.HandleFunc("/top-fat-characters", a.topFatCharacters)
}

func (a *API) ui(w http.ResponseWriter, r *http.Request) {
	characters, err := a.core.TopFatCharacters()
	if err != nil {
		http.Error(w, "unknown error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	data := struct {
		Characters []starwars.Character
	}{
		Characters: characters,
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, data); err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Printf("Error rendering page: %v", err)
		return
	}

	buf.WriteTo(w)
}

func (a *API) topFatCharacters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	characters, err := a.core.TopFatCharacters()
	if err != nil {
		http.Error(w, "unknown error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("content-type", "application/json")

	json.NewEncoder(w).Encode(&characters)
}
