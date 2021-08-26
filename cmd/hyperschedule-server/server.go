package main

import (
	"context"
	"fmt"
	"github.com/MuddCreates/hyperschedule-api-go/internal/db"
	"github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
	"log"
	"net/http"
	"time"
)

type Server struct {
	*http.Server
}

type Context struct {
	uploaderHash string
	dbConn       *db.Connection
	oldState     *OldState
	apiCache     *cache.Client
}

func (c *Cmd) NewServer() (*Server, error) {
	conn, err := db.New(context.Background(), c.DbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}

	cacheAdapter, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.MFU),
		memory.AdapterWithCapacity(8),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cache adapter: %w", err)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(cacheAdapter),
		cache.ClientWithTTL(time.Minute),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cache middleware: %w", err)
	}

	mux := http.NewServeMux()
	ctx := &Context{
		dbConn:       conn,
		uploaderHash: c.UploadEmailHash,
		oldState:     &OldState{},
		apiCache:     cacheClient,
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
