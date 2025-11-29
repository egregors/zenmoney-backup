// Package srv provides backup server functionality.
package srv

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	log "github.com/go-pkgz/lgr"
	"github.com/nemirlev/zenmoney-go-sdk/v2/api"
)

// Saver is interface for file storage.
type Saver interface {
	Save(filename string, bs []byte) error
}

// Notifier is an interface for sending notifications.
type Notifier interface {
	Notify(title, message string) error
}

// Server is backup server.
type Server struct {
	token     string
	sleepTime time.Duration
	timeout   time.Duration
	store     Saver
	client    *api.Client
	notifier  Notifier
}

// NewServer makes Server from options.
func NewServer(token string, sleepTime time.Duration, timeout time.Duration, storage Saver, notifier Notifier) *Server {
	return &Server{
		token:     token,
		sleepTime: sleepTime,
		timeout:   timeout,
		store:     storage,
		notifier:  notifier,
	}
}

// Run starts Server.
func (srv *Server) Run(ctx context.Context) {
	log.Printf("[INFO] login...")

	// Configure HTTP transport with proper timeouts to avoid TLS handshake timeout issues
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   30 * time.Second,
		ResponseHeaderTimeout: srv.timeout,
		IdleConnTimeout:       90 * time.Second,
		MaxIdleConns:          10,
		MaxIdleConnsPerHost:   5,
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   srv.timeout,
	}

	log.Printf("[DEBUG] creating API client with timeout=%s, TLS handshake timeout=30s", srv.timeout)

	// Note: We pass both WithHTTPClient and WithTimeout because:
	// - httpClient.Timeout is the effective timeout enforced by Go's http package
	// - WithTimeout passes the timeout to the SDK for internal configuration
	client, err := api.NewClient(
		srv.token,
		api.WithHTTPClient(httpClient),
		api.WithTimeout(srv.timeout),
	)
	if err != nil {
		log.Printf("[ERROR] failed to create client: %s", err)
		srv.sendNotification("Client Creation Error", err.Error())
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
		srv.sendNotification("Backup Export Error", err.Error())
		return
	}

	fileName := srv.genFileName(time.Now())
	err = srv.store.Save(fileName, bs)
	if err != nil {
		log.Printf("[ERROR] downloading failed: %s", err)
		srv.sendNotification("Backup Save Error", err.Error())
		return
	}
	log.Printf("[INFO] %s saved", fileName)
	log.Printf("[INFO] sleep for %s", srv.sleepTime.String())
}

func (srv *Server) export(ctx context.Context) ([]byte, error) {
	log.Printf("[DEBUG] downloading data with timeout=%s ...", srv.timeout)
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(ctx, srv.timeout)
	defer cancel()

	resp, err := srv.client.FullSync(ctx)
	if err != nil {
		elapsed := time.Since(startTime)
		log.Printf("[ERROR] failed to download data after %s: %s", elapsed, err)
		return nil, err
	}
	_ = resp

	elapsed := time.Since(startTime)
	log.Printf("[DEBUG] API request completed in %s", elapsed)

	bs, err := json.Marshal(resp)
	if err != nil {
		log.Printf("[ERROR] failed to marshal data: %s", err)
		return nil, err
	}

	log.Printf("[DEBUG] downloaded")
	return bs, nil
}

func (srv *Server) sendNotification(title, message string) {
	if srv.notifier != nil {
		if err := srv.notifier.Notify(title, message); err != nil {
			log.Printf("[WARN] failed to send notification: %s", err)
		}
	}
}

func (srv *Server) genFileName(t time.Time) string {
	return fmt.Sprintf("zen_%s.json", t.Format("2006-01-02_15-04-05"))
}
