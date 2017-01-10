package webapp

import (
	"html/template"
	"net/http"
	
	jd "github.com/josephburnett/jd/lib"
)

type Mode string

const (
	DiffMode = Mode("diff")
	PatchMode = Mode("patch")
)

type State struct {
	InputA string
	NodeA jd.JsonNode
	NodeAError error
	InputB string
	NodeB jd.JsonNode
	NodeBError error
	InputDiff string
	Diff jd.Diff
	DiffError error
	Mode Mode
	Debug string
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
		err = index.Execute(w, &State{
			InputA: `{"foo":"bar"}`,
			InputB: `{"foo":"baz"}`,
			Mode: DiffMode,			
		})
	case http.MethodPost:
		s := getState(r)
		if s.Mode == DiffMode {
			s = diff(s)
		}
		if s.Mode == PatchMode {
			s = patch(s)
		}
		err = index.Execute(w, s)
	default:
		http.Error(w, "Unsupported method.", 405)
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func getState(r *http.Request) *State {
	state := &State{
		InputA: r.FormValue("a"),
		InputB: r.FormValue("b"),
		InputDiff: r.FormValue("diff"),
		Mode: Mode(r.FormValue("mode")),
	}
	state.NodeA, state.NodeAError = jd.ReadJsonString(state.InputA)
	switch state.Mode {
	case DiffMode:
		state.NodeB, state.NodeBError = jd.ReadJsonString(state.InputB)
	case PatchMode:
		state.Diff, state.DiffError = jd.ReadDiffString(state.InputDiff)
	}
	return state
}

func diff(s *State) *State {
	if s.NodeA != nil && s.NodeB != nil {
		d := s.NodeA.Diff(s.NodeB)
	    s.Diff = d
		s.InputDiff = d.Render()
	}
	return s
}

func patch(s *State) *State {
	if s.NodeA != nil && s.Diff != nil {
		s.NodeB, s.DiffError = s.NodeA.Patch(s.Diff)
		if s.DiffError == nil {
			s.InputB = s.NodeB.Json()
		}
	}
	return s
}