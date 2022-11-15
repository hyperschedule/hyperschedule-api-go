package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/MuddCreates/hyperschedule-api-go/internal/lingk"
	"io"
	"log"
	"net/http"
)

func (ctx *Context) iftttHandler(resp http.ResponseWriter, req *http.Request) {

	log.Printf("IFTTT upload")

	token := req.Header.Get("Hyperschedule-Upload-Token")
	if ctx.uploadTokenHash != fmt.Sprintf("%x", sha256.Sum256([]byte(token))) {
		resp.WriteHeader(http.StatusUnauthorized)
		resp.Write([]byte("missing or invalid auth token"))
		return
	}

	attachmentURL := req.URL.Query().Get("attachment")
	attachmentResp, err := http.Get(attachmentURL)
	if err != nil {
		log.Printf("IFTTT: attachment url request failed: %v", err)
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte(fmt.Sprintf("request to attachment url failed: %v", err)))
		return
	}
	defer attachmentResp.Body.Close()

	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, attachmentResp.Body); err != nil {
		log.Printf("IFTTT: failed to download attachment: %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(fmt.Sprintf("IFTTT: failed to download attachment: %v", err)))
		return
	}

	data, err := lingk.FromZipBuffer(buf)
	if err != nil {
		log.Printf("IFTTT UPLOAD: failed to parse, %v", err)
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte(fmt.Sprintf("IFTTT UPLOAD: failed to parse, %v", err)))
		return
	}

	select {
	case ctx.updateChan <- data:
		log.Printf("IFTTT UPLOAD: successfully parsed, uploading to db")
		resp.WriteHeader(http.StatusOK)
	default:
		log.Printf("received request while busy")
		resp.WriteHeader(http.StatusServiceUnavailable)
	}

}
