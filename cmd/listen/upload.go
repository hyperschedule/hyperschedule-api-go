package main

import (
  "encoding/json"
	"errors"
	"log"
	"net/http"
  "mime/multipart"
)

type LingkEmail struct {
	From string
  To string
  Envelope *Envelope
  Attachment *multipart.FileHeader
}

type Envelope struct {
  From string `json:"from"`
  To string `json:"to"`
}

func parseEmail(req *http.Request) (*LingkEmail, error) {
	if err := req.ParseMultipartForm(0); err != nil {
		return nil, err
	}

	from := req.MultipartForm.Value["from"]
	if len(from) == 0 {
		return nil, errors.New("missing from")
	}

	to := req.MultipartForm.Value["to"]
	if len(to) == 0 {
		return nil, errors.New("missing to")
	}

  envelopes := req.MultipartForm.Value["envelope"]
  if len(envelopes) == 0 {
    return nil, errors.New("missing envelope")
  }
  var envelope *Envelope
  if err := json.Unmarshal([]byte(envelopes[0]), envelope); err != nil {
    return nil, errors.New("failed to parse envelope json")
  }

  var attachment *multipart.FileHeader
  for _, fhs := range req.MultipartForm.File {
    for _, fh := range fhs {
      if fh.Filename != "HMCarchive.zip" {
        return nil, errors.New("unrecognized attachment filename")
      }
      attachment = fh
    }
  }

	return &LingkEmail{
		From: from[0],
    To: to[0],
    Envelope: envelope,
    Attachment: attachment,
	}, nil
}

func inboundHandler(resp http.ResponseWriter, req *http.Request) {
	log.Printf(
		"UPLOAD: request from %s (forwarded from %s)",
		req.RemoteAddr,
		req.Header["X-Forwarded-For"],
	)

	email, err := parseEmail(req)
	if err != nil {
		log.Printf("UPLOAD: Failed to parse email from request: %v", err)

		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("no way jose"))
    return
	}

  log.Printf("UPLOAD: successfully parsed email, from = %s, to = %s", email.Envelope.From, email.Envelope.To)
	resp.WriteHeader(http.StatusOK)
}
