//go:build include_web

package serve

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

func init() {
	Handle = handle
}

func handle(w http.ResponseWriter, r *http.Request) {
	f := r.URL.Path[1:]
	if f == "" {
		f = "index.html"
	}
	s, ok := base64EncodedFiles[f]
	if !ok {
		http.Error(w, fmt.Sprintf("file %q not found", f), http.StatusNotFound)
		return
	}
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		http.Error(w, "error base64 decoding", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
