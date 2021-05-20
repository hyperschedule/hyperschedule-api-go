package main

import (
	"github.com/davecgh/go-spew/spew"
	"net/http"
)

func rawHandler(rw http.ResponseWriter, req *http.Request) {
	data := state.GetData()
	rw.Header().Add("Content-Type", "text/plain")
	spew.Fdump(rw, data)
}

func rawStaffHandler(rw http.ResponseWriter, req *http.Request) {
	data := state.GetData()
	rw.Header().Add("Content-Type", "text/plain")
	spew.Fdump(rw, data.Staff)
}
