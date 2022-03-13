package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/go-pkgz/lgr"
)

const (
	sessionIDKey     = "PHPSESSID"
	loginResultError = "error"
	loginResultDone  = "done"
)

// Saver is interface for file storage
type Saver interface {
	Save(filename string, bs []byte) error
}

// Server is backup server
type Server struct {
	username, password string
	sessionID          string
	sleepTime          time.Duration
	store              Saver
}

// NewServer makes Server from options
func NewServer(login, password string, sleepTime time.Duration, storage Saver) *Server {
	return &Server{
		username:  login,
		password:  password,
		sleepTime: sleepTime,
		store:     storage,
	}
}

// Run starts Server
func (srv *Server) Run(ctx context.Context) {
	log.Printf("[INFO] login...")
	err := srv.login()
	if err != nil {
		log.Printf("[FATAL] can't login: %s", err)
		return
	}

	srv.saveExport()

	ticker := time.NewTicker(srv.sleepTime)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			// todo: need a way how to figure out
			// 	that session id is expired and export is not
			// 	valid any more.
			srv.saveExport()
		}
	}
}

func (srv *Server) saveExport() {
	log.Printf("[INFO] downloading...")
	bs, err := srv.export()
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

func (srv *Server) makeLoginRequest() *http.Request {
	params := url.Values{}
	params.Add("login", srv.username)
	params.Add("password", srv.password)
	params.Add("remember", `true`)
	params.Add("check", `true`)
	body := strings.NewReader(params.Encode())

	req, _ := http.NewRequest("POST", "https://zenmoney.ru/login/enter/", body)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:98.0) Gecko/20100101 Firefox/98.0")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", "https://zenmoney.ru")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://zenmoney.ru/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Te", "trailers")

	return req
}

func (srv *Server) makeExportRequest() *http.Request {
	params := url.Values{}
	params.Add("month_filter", `0`)
	params.Add("account_filter", `0`)
	body := strings.NewReader(params.Encode())

	// todo: export URI
	req, _ := http.NewRequest("POST", "https://zenmoney.ru/export/download2011/", body)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:98.0) Gecko/20100101 Firefox/98.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://zenmoney.ru")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://zenmoney.ru/a/")
	req.Header.Set("Cookie", fmt.Sprintf("PHPSESSID=%s; uechat_1_mode=0", srv.sessionID))
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Te", "trailers")

	return req
}

func (srv *Server) login() error {
	// login response is always 200, need to check JSON body
	resp, err := http.DefaultClient.Do(srv.makeLoginRequest())
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	login := struct {
		Result  string `json:"result"`
		Login   string `json:"login"`
		Message string `json:"message"`
	}{}

	// todo: error?
	bs, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(bs, &login)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] got login: %s", login)

	switch login.Result {
	case loginResultError:
		return fmt.Errorf("%s", login.Message)
	case loginResultDone:
		log.Printf("[DEBUG] logged in")
	default:
		return fmt.Errorf("got unexpected answer: %s", login.Result)
	}

	for _, c := range resp.Cookies() {
		if c.Name == sessionIDKey {
			srv.sessionID = c.Value
			return nil
		}
	}

	return nil
}

func (srv *Server) export() ([]byte, error) {
	log.Printf("[DEBUG] downloading export file ...")
	resp, err := http.DefaultClient.Do(srv.makeExportRequest())
	if err != nil {
		return nil, fmt.Errorf("can't export login: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	bs, _ := io.ReadAll(resp.Body)
	log.Printf("[DEBUG] downloaded")
	return bs, nil
}

func (srv *Server) genFileName(t time.Time) string {
	return fmt.Sprintf("zen_%s.csv", t.Format("2006-01-02_15-04-05"))
}
