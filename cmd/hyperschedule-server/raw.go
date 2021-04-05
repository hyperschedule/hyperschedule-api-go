package main

import (
	"encoding/json"
	"net/http"
)

func rawHandler(rw http.ResponseWriter, req *http.Request) {

	data := state.GetData()
	raw, err := json.Marshal(data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(raw)
}
