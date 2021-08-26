package main

import (
	"encoding/json"
	"github.com/MuddCreates/hyperschedule-api-go/internal/api"
	"log"
	"net/http"
)

func (ctx *Context) apiHandler() http.Handler {
	mux := http.NewServeMux()

	stripQueryCache := func(
		f func(http.ResponseWriter, *http.Request),
	) http.Handler {
		return stripQuery(ctx.apiCache.Middleware(http.HandlerFunc(f)))
	}

	mux.Handle("/v3/courses", stripQueryCache(ctx.apiV3Handler))
	mux.Handle("/v3-new/courses", stripQueryCache(ctx.apiV3NewHandler))

	return mux
}

func stripQuery(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		req.URL.RawQuery = ""
		handler.ServeHTTP(resp, req)
	})
}

func getIp(req *http.Request) string {
	f := req.Header.Get("X-Forwarded-For")
	if len(f) != 0 {
		return f
	}
	return req.RemoteAddr
}

func (ctx *Context) apiV3Handler(resp http.ResponseWriter, req *http.Request) {
	data := ctx.oldState.GetData()
	if data == nil {
		log.Printf("received api request before loaded")
		resp.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	output, err := json.Marshal(api.MakeV3(data))
	if err != nil {
		log.Printf("api: failed to jsonify, %s", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.Header().Add("Content-Type", "application/json")
	resp.Header().Add("Access-Control-Allow-Origin", "*")
	resp.Write(output)
}

func (ctx *Context) apiV3NewHandler(resp http.ResponseWriter, req *http.Request) {
	data, err := api.FetchV3(req.Context(), ctx.dbConn)
	if err != nil {
		log.Printf("API: failed to retrieve data from db: %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	output, err := json.Marshal(data)
	if err != nil {
		log.Printf("api: failed to jsonify, %s", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.Header().Add("Content-Type", "application/json")
	resp.Header().Add("Access-Control-Allow-Origin", "*")
	resp.Write(output)
}
