package webapp

import (
	"html/template"
	"net/http"
	
	jd "github.com/josephburnett/jd/lib"
)

type State struct {
	A jd.JsonNode
	B jd.JsonNode
	Diff jd.Diff
	Error error
}

var index *template.Template

func init() {
	http.HandleFunc("/", handler)
	index, _ = template.ParseFiles("index.html")
}

func handler(w http.ResponseWriter, r *http.Request) {
	var err error = nil
	switch r.Method {
	case http.MethodGet:
		err = index.Execute(w, &State{})
	case http.MethodPost:
		err = index.Execute(w, &State{})
	default:
		http.Error(w, "Unsupported method.", 405)
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
