package main

import (
	"errors"
	"log"
	"net/http"
  "mime/multipart"
)

type LingkEmail struct {
	From string
  To string
  Attachment *multipart.FileHeader
}

func parseEmail(req *http.Request) (*LingkEmail, error) {
	if err := req.ParseMultipartForm(0); err != nil {
		return nil, err
	}

	from := req.MultipartForm.Value["from"]
	if len(from) == 0 {
		return nil, errors.New("missing from")
	}

	to := req.MultipartForm.Value["from"]
	if len(to) == 0 {
		return nil, errors.New("missing from")
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

  log.Printf("UPLOAD: successfully parsed email, from = %s, to = %s", email.From, email.To)
	resp.WriteHeader(http.StatusOK)
}
