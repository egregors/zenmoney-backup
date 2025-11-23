// Package notifier provides notification functionality for errors.
package notifier

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"
)

// Notifier is an interface for sending notifications.
type Notifier interface {
	Notify(title, message string) error
}

// Noop is a no-op notifier that does nothing.
type Noop struct{}

// NewNoop creates a new Noop notifier.
func NewNoop() Noop {
	return Noop{}
}

// Notify does nothing for Noop notifier.
func (n Noop) Notify(_, _ string) error {
	return nil
}

// Ntfy is a notifier that sends notifications to ntfy.sh.
type Ntfy struct {
	url string
}

// NewNtfy creates a new Ntfy notifier with the given URL.
func NewNtfy(url string) *Ntfy {
	return &Ntfy{url: url}
}

// Notify sends a notification to ntfy.sh with the given title and message.
func (n Ntfy) Notify(title, message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	req, err := http.NewRequestWithContext(ctx, "POST", n.url, strings.NewReader(message))
	if err != nil {
		return err
	}
	req.Header.Set("Title", title)
	req.Header.Set("Tags", "warning,zenmoney-backup")
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	// Discard response body to ensure connection reuse
	_, _ = io.Copy(io.Discard, resp.Body)

	return nil
}
