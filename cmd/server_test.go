package cmd

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getCredsFromENV() (username, password string, ok bool) {
	username, uOk := os.LookupEnv("ZEN_USERNAME")
	password, pOk := os.LookupEnv("ZEN_PASSWORD")
	if uOk && pOk && username != "" && password != "" {
		return username, password, true
	}
	return
}

type saverMock struct{}

func (s saverMock) Save(_ string, _ []byte) error {
	return nil
}

func TestNewServer(t *testing.T) {
	login, pass := "login", "pass"
	dur := time.Minute * 30

	s := NewServer(login, pass, dur, saverMock{})
	assert.Equal(t, login, s.username)
	assert.Equal(t, pass, s.password)
	assert.Equal(t, dur, dur)
}

func TestServer_genFileName(t *testing.T) {
	s := Server{}
	bT, _ := time.Parse("2006-01-02_15-04-05", "2022-03-12_21-48-00")
	assert.Equal(
		t,
		"zen_2022-03-12_21-48-00.csv",
		s.genFileName(bT),
	)
}

func TestServer_makeLoginRequest(t *testing.T) {
	s := Server{username: "McDuck", password: "ilikemoney"}
	req := s.makeLoginRequest()
	assert.Equal(t, req.URL.String(), "https://zenmoney.ru/login/enter/")
	params, _ := io.ReadAll(req.Body)
	assert.Equal(t, "check=true&login=McDuck&password=ilikemoney&remember=true", string(params))
}

func TestServer_makeExportRequest(t *testing.T) {
	s := Server{}
	req := s.makeExportRequest()
	assert.Equal(t, req.URL.String(), "https://zenmoney.ru/export/download2011/")
	params, _ := io.ReadAll(req.Body)
	assert.Equal(t, "account_filter=0&month_filter=0", string(params))
}

func TestServer_login(t *testing.T) {
	if u, p, ok := getCredsFromENV(); ok {
		// wrong login \ pass
		s := Server{username: u, password: "foobar"}
		err := s.login()
		assert.EqualError(t, err, "Неправильное имя пользователя или пароль.")

		// success
		s = Server{username: u, password: p}
		err = s.login()
		assert.NoError(t, err)
	} else {
		t.Fatalf("can't get test username or passworn from ENV")
	}
}

func TestServer_export(t *testing.T) {
	if u, p, ok := getCredsFromENV(); ok {
		s := Server{username: u, password: p}
		err := s.login()
		assert.NoError(t, err)

		bs, _ := s.export()
		rs := csv.NewReader(bytes.NewReader(bs))

		_, _ = rs.Read() // skipp first line
		_, _ = rs.Read() // skipp second line

		fLine, _ := rs.Read()
		sLine, _ := rs.Read()

		assert.Equal(
			t,
			[]string{
				"2022-03-12", "Продукты", "Пятерочка", "Это демо-операция. Удалите её, когда посмотрите отчёты.",
				"Карта", "1450,00", "RUB", "", "", "", "2022-03-13 23:53:28", "2022-03-13 23:53:28",
			},
			fLine,
		)
		assert.Equal(
			t,
			[]string{"2022-03-13", "Кафе и рестораны", "Любимая кофейня", "Это демо-операция. Удалите её, когда посмотрите отчёты.",
				"Карта", "160,00", "RUB", "", "", "", "2022-03-13 23:53:28", "2022-03-13 23:53:28",
			},
			sLine,
		)

	} else {
		t.Fatalf("can't get test username or passworn from ENV")
	}
}
