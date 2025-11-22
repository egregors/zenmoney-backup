// Package notifier provides notification functionality for errors.
package notifier

import (
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
	req, err := http.NewRequest("POST", n.url, strings.NewReader(message))
	if err != nil {
		return err
	}
	req.Header.Set("Title", title)
	req.Header.Set("Tags", "warning,zenmoney-backup")
	
	// Use a client with timeout to prevent hanging
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
