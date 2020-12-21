package main

import (
  "fmt"
  "log"
  "net/http"
  "github.com/sendgrid/sendgrid-go/helpers/inbound"
)

func inboundHandler(resp http.ResponseWrite, req *http.Request) {
  email := inbound.Parse(req)
  log.Printf("got email from %#v", email.Headers["From"])

  for f, _ := range email.Attachments {
    log.Printf("has attachments %#v", f)
  }

  for sec, _ := range email.Body {
    log.Printf("has body %#v", sec);
  }

  resp.WriteHeader(http.StatusOK)
}

func main() {
  http.HandleFunc("/upload", inboundHandler)
  if err := http.ListenAndServe(":9123", nil); err != nil {
    log.Fatalf("bad %v", err)
  }
}
