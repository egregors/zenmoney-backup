package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMakeServer(t *testing.T) {
	tests := []struct {
		name        string
		opts        Opts
		shouldError bool
		errorMsg    string
	}{
		{
			name: "valid options",
			opts: Opts{
				Token:     "test_token",
				SleepTime: "24h",
				Timeout:   10,
			},
			shouldError: false,
		},
		{
			name: "zero timeout",
			opts: Opts{
				Token:     "test_token",
				SleepTime: "24h",
				Timeout:   0,
			},
			shouldError: true,
			errorMsg:    "timeout must be a positive integer, got 0",
		},
		{
			name: "negative timeout",
			opts: Opts{
				Token:     "test_token",
				SleepTime: "24h",
				Timeout:   -5,
			},
			shouldError: true,
			errorMsg:    "timeout must be a positive integer, got -5",
		},
		{
			name: "invalid sleep time",
			opts: Opts{
				Token:     "test_token",
				SleepTime: "invalid",
				Timeout:   10,
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := makeServer(tt.opts)
			
			if tt.shouldError {
				assert.Error(t, err)
				assert.Nil(t, server)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, server)
			}
		})
	}
}