package main

import (
  "log"
  "net/http"
  "github.com/davecgh/go-spew/spew"
  "io/ioutil"
  //"github.com/sendgrid/sendgrid-go/helpers/inbound"
)

func inboundHandler(resp http.ResponseWriter, req *http.Request) {
  log.Printf("got request from %s (forwarded from %s)", req.RemoteAddr, req.Header["X-Forwarded-For"])
  spew.Dump(req.Header)
  err := req.ParseMultipartForm(0)
  if err != nil {
    log.Printf("invalid request")
    return
  }

  spew.Dump(req.MultipartForm)

  for key := range req.MultipartForm.Value {
	  log.Printf("value key: %#v", key)
  }
  for key := range req.MultipartForm.File {
	  log.Printf("file key: %#v", key)
  }
  a, err := ioutil.ReadDir("/run/user/1000")
  if err != nil { log.Fatal(err) }
  spew.Dump(a)

  //email := inbound.Parse(req)
  //log.Printf("got email from %#v", email.Headers["From"])

  //for f, _ := range email.Attachments {
  //  log.Printf("has attachments %#v", f)
  //}

  //for sec, _ := range email.Body {
  //  log.Printf("has body %#v", sec);
  //}

  resp.WriteHeader(http.StatusOK)
}

func main() {
  http.HandleFunc("/upload/", inboundHandler)
  if err := http.ListenAndServe(":9123", nil); err != nil {
    log.Fatalf("bad %v", err)
  }
}
