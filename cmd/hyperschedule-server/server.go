package main

import (
	"context"
	"fmt"
	"github.com/MuddCreates/hyperschedule-api-go/internal/data"
	"github.com/MuddCreates/hyperschedule-api-go/internal/db"
	"github.com/MuddCreates/hyperschedule-api-go/internal/update"
	"github.com/kr/pretty"
	"github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
	"log"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	server *http.Server
	ctx    *Context
}

type Context struct {
	uploaderHash    string
	uploadTokenHash string
	updateChan      chan *data.Data
	dbConn          *db.Connection
	oldState        *OldState
	apiCache        *cache.Client
	apiV3CacheData  []byte
	apiV3CacheTime  time.Time
	apiV3CacheMutex *sync.RWMutex
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
		dbConn:          conn,
		updateChan:      make(chan *data.Data),
		uploaderHash:    c.UploadEmailHash,
		uploadTokenHash: c.UploadTokenHash,
		oldState:        &OldState{},
		apiCache:        cacheClient,
		apiV3CacheData:  nil,
		apiV3CacheMutex: &sync.RWMutex{},
		apiV3CacheTime:  time.Unix(0, 0),
	}
	mux.HandleFunc("/upload/", ctx.inboundHandler)
	mux.HandleFunc("/ifttt/", ctx.iftttHandler)
	mux.HandleFunc("/raw/", ctx.rawHandler)
	mux.HandleFunc("/raw/staff", ctx.rawStaffHandler)

	mux.Handle("/api/", http.StripPrefix("/api", ctx.apiHandler()))

	s := &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", c.Port),
			Handler: mux,
		},
		ctx: ctx,
	}
	return s, nil
}

func (s *Server) Run() error {
	log.Printf("listening on %#v", s.server.Addr)
	go s.ctx.listenUpdates()
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (ctx *Context) listenUpdates() {
	for {
		data := <-ctx.updateChan

		summaries, err := update.Update(context.Background(), ctx.dbConn, data)
		if err != nil {
			log.Printf("UPLOAD: failed to update DB: %v", err)
			return
		}
		ctx.apiV3CacheMutex.Lock()
		ctx.apiV3CacheData = nil
		ctx.apiV3CacheMutex.Unlock()
		pretty.Logf("update info: %# v", summaries)
	}
}
