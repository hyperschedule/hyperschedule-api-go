package main

import (
  "log"
  "net/http"
  "encoding/json"
)

func apiV3Handler(resp http.ResponseWriter, req *http.Request) {
  output, err := json.Marshal(state.GetData())
  if err != nil {
    log.Printf("api: failed, %s", err)
    resp.WriteHeader(http.StatusInternalServerError)
  }
  resp.Header().Add("Content-Type", "application/json")
  resp.Write(output)
}
