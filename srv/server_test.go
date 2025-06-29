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

	s := NewServer(token, dur, saverMock{})
	assert.NotNil(t, s)
	assert.Equal(t, token, s.token)
	assert.Equal(t, dur, s.sleepTime)
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
