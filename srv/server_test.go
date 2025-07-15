package srv

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type saverMock struct{}

func (s saverMock) Save(_ string, _ []byte) error {
	return nil
}

func TestNewServer(t *testing.T) {
	token := "test_token"
	dur := time.Minute * 30
	timeout := time.Second * 15

	s := NewServer(token, dur, timeout, saverMock{})
	assert.NotNil(t, s)
	assert.Equal(t, token, s.token)
	assert.Equal(t, dur, s.sleepTime)
	assert.Equal(t, timeout, s.timeout)
}

func TestServer_TimeoutConfiguration(t *testing.T) {
	// Test that different timeout values are properly set
	tests := []struct {
		name            string
		timeoutSeconds  int
		expectedTimeout time.Duration
	}{
		{"default timeout", 10, 10 * time.Second},
		{"custom timeout", 30, 30 * time.Second},
		{"short timeout", 5, 5 * time.Second},
		{"long timeout", 60, 60 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer("test_token", time.Hour, tt.expectedTimeout, saverMock{})
			assert.Equal(t, tt.expectedTimeout, s.timeout)
		})
	}
}

func TestServer_genFileName(t *testing.T) {
	s := Server{}
	bT, _ := time.Parse("2006-01-02_15-04-05", "2022-03-12_21-48-00")
	assert.Equal(
		t,
		"zen_2022-03-12_21-48-00.json",
		s.genFileName(bT),
	)
}
