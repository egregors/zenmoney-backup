// Package srv provides backup server functionality.
package srv

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/go-pkgz/lgr"
	"github.com/nemirlev/zenmoney-go-sdk/v2/api"
)

// Saver is interface for file storage.
type Saver interface {
	Save(filename string, bs []byte) error
}

// Server is backup server.
type Server struct {
	token     string
	sleepTime time.Duration
	timeout   time.Duration
	store     Saver
	client    *api.Client
}

// NewServer makes Server from options.
func NewServer(token string, sleepTime time.Duration, timeout time.Duration, storage Saver) *Server {
	return &Server{
		token:     token,
		sleepTime: sleepTime,
		timeout:   timeout,
		store:     storage,
	}
}

// Run starts Server.
func (srv *Server) Run(ctx context.Context) {
	log.Printf("[INFO] login...")

	client, err := api.NewClient(srv.token)
	if err != nil {
		log.Printf("[ERROR] failed to create client: %s", err)
		return
	}
	srv.client = client

	srv.saveExport(ctx)

	ticker := time.NewTicker(srv.sleepTime)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			srv.saveExport(ctx)
		}
	}
}

func (srv *Server) saveExport(ctx context.Context) {
	log.Printf("[INFO] downloading...")
	bs, err := srv.export(ctx)
	if err != nil {
		log.Printf("[ERROR] failed: %s", err)
		return
	}

	fileName := srv.genFileName(time.Now())
	err = srv.store.Save(fileName, bs)
	if err != nil {
		log.Printf("[ERROR] downloading failed: %s", err)
		return
	}
	log.Printf("[INFO] %s saved", fileName)
	log.Printf("[INFO] sleep for %s", srv.sleepTime.String())
}

func (srv *Server) export(ctx context.Context) ([]byte, error) {
	log.Printf("[DEBUG] downloading data ...")
	ctx, cancel := context.WithTimeout(ctx, srv.timeout)
	defer cancel()

	resp, err := srv.client.FullSync(ctx)
	if err != nil {
		log.Printf("[ERROR] failed to download data: %s", err)
		return nil, err
	}
	_ = resp

	bs, err := json.Marshal(resp)
	if err != nil {
		log.Printf("[ERROR] failed to marshal data: %s", err)
		return nil, err
	}

	log.Printf("[DEBUG] downloaded")
	return bs, nil
}

func (srv *Server) genFileName(t time.Time) string {
	return fmt.Sprintf("zen_%s.json", t.Format("2006-01-02_15-04-05"))
}
