package json

import (
	"encoding/json"
	"net/http"
)

func ServeJson(w http.ResponseWriter, status int, v interface{}) error {

	// don't use Header().Set() or Header().Add(): https://github.com/golang/go/issues/5022
	w.Header()["Content-Type"] = []string{"application/json"}
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	return nil
}

func ServeHalJson(w http.ResponseWriter, status int, v interface{}) error {

	// don't use Header().Set() or Header().Add(): https://github.com/golang/go/issues/5022
	w.Header()["Content-Type"] = []string{"application/hal+json"}
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	return nil
}
