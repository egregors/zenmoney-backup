package notifier

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNoop(t *testing.T) {
	n := NewNoop()
	assert.NotNil(t, n)
}

func TestNoop_Notify(t *testing.T) {
	n := NewNoop()
	err := n.Notify("test title", "test message")
	assert.NoError(t, err)
}

func TestNewNtfy(t *testing.T) {
	url := "https://ntfy.sh/mytopic"
	n := NewNtfy(url)
	assert.NotNil(t, n)
	assert.Equal(t, url, n.url)
}

func TestNtfy_Notify(t *testing.T) {
	t.Run("successful notification", func(t *testing.T) {
		// Create a test server to mock ntfy.sh
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "Test Error", r.Header.Get("Title"))
			assert.Equal(t, "warning,zenmoney-backup", r.Header.Get("Tags"))
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		n := NewNtfy(server.URL)
		err := n.Notify("Test Error", "test error message")
		assert.NoError(t, err)
	})

	t.Run("server error", func(t *testing.T) {
		// Create a test server that returns an error
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		n := NewNtfy(server.URL)
		// Should not return error because we don't check response status code
		err := n.Notify("Test Error", "test error message")
		assert.NoError(t, err)
	})
}
