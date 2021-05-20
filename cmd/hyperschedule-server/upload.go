package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MuddCreates/hyperschedule-api-go/internal/lingk"
	"github.com/kr/pretty"
	"log"
	"mime/multipart"
	"net/http"
)

type LingkEmail struct {
	Envelope   Envelope
	Attachment *multipart.FileHeader
}

type Envelope struct {
	From string   `json:"from"`
	To   []string `json:"to"`
}

func parseEmail(req *http.Request) (*LingkEmail, error) {
	if err := req.ParseMultipartForm(0); err != nil {
		return nil, err
	}

	envelopes := req.MultipartForm.Value["envelope"]
	if len(envelopes) == 0 {
		return nil, errors.New("missing envelope")
	}
	var envelope Envelope
	if err := json.Unmarshal([]byte(envelopes[0]), &envelope); err != nil {
		return nil, errors.New("failed to parse envelope json")
	}

	var attachment *multipart.FileHeader
	for _, fhs := range req.MultipartForm.File {
		for _, fh := range fhs {
			if fh.Filename != "HMCarchive.zip" {
				return nil, errors.New(fmt.Sprintf("unrecognized attachment filename %#v", fh.Filename))
			}
			attachment = fh
		}
	}
	if attachment == nil {
		return nil, errors.New("missing HMCarchive.zip attachment")
	}

	return &LingkEmail{
		Envelope:   envelope,
		Attachment: attachment,
	}, nil
}

func validateEmail(email *LingkEmail) error {
	if len(email.Envelope.To) != 1 {
		return errors.New("wrong number of email tos")
	}

	to := email.Envelope.To[0]
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(to)))
	if hash != uploadEmailHash {
		log.Printf("expected %s, got %s (pre-hash %s)", uploadEmailHash, hash[:], to)
		return errors.New("hash mismatch, get rekt")
	}

	return nil
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

	if err := validateEmail(email); err != nil {
		log.Printf("UPLOAD: wrong target email (unauthorized), %v", err)
		resp.WriteHeader(http.StatusUnauthorized)
		resp.Write([]byte("nice try"))
		return
	}

	data, err := lingk.FromAttachment(email.Attachment)
	if err != nil {
		log.Printf("UPLOAD: failed to parse, %s", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	updateInfo, err := dbConn.Update(req.Context(), data)
	if err != nil {
		log.Printf("UPLOAD: failed to update DB, %v", err)
		return
	}
	pretty.Logf("update info: %# v", updateInfo)

	state.SetData(data)

	log.Printf("UPLOAD: successfully parsed email")
	resp.WriteHeader(http.StatusOK)
}
