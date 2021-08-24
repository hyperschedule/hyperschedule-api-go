package main

import (
	"context"
	"fmt"
	"github.com/MuddCreates/hyperschedule-api-go/internal/db"
	"log"
	"net/http"
)

type Server struct {
	*http.Server
}

type Context struct {
	uploaderHash string
	dbConn       *db.Connection
	oldState     *OldState
}

func (c *Cmd) NewServer() (*Server, error) {
	conn, err := db.New(context.Background(), c.DbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}

	// TODO get rid of wack recursion
	mux := http.NewServeMux()
	ctx := &Context{
		dbConn:       conn,
		uploaderHash: c.UploadEmailHash,
		oldState:     &OldState{},
	}
	mux.HandleFunc("/upload/", ctx.inboundHandler)
	mux.HandleFunc("/raw/", ctx.rawHandler)
	mux.HandleFunc("/raw/staff", ctx.rawStaffHandler)

	mux.Handle("/api/", http.StripPrefix("/api", ctx.apiHandler()))

	s := &Server{
		&http.Server{
			Addr:    fmt.Sprintf(":%d", c.Port),
			Handler: mux,
		},
	}
	return s, nil
}

func (s *Server) Run() error {
	log.Printf("listening on %#v", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
