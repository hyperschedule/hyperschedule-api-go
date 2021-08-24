package main

import (
	"github.com/davecgh/go-spew/spew"
	"net/http"
)

func (ctx *Context) rawHandler(rw http.ResponseWriter, req *http.Request) {
	data := ctx.oldState.GetData()
	rw.Header().Add("Content-Type", "text/plain")
	spew.Fdump(rw, data)
}

func (ctx *Context) rawStaffHandler(rw http.ResponseWriter, req *http.Request) {
	data := ctx.oldState.GetData()
	rw.Header().Add("Content-Type", "text/plain")
	spew.Fdump(rw, data.Staff)
}
