package actions_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sethetter/go-web-starter/pkg/config"
	"github.com/sethetter/go-web-starter/pkg/server"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/publicsuffix"
)

// Test Helpers ------------------------------

type email struct {
	from, recipient, subject, body string
}

type mockService struct {
	emails []email
}

func (svc *mockService) SendEmail(from, recipient, subject, body string) error {
	svc.emails = append(svc.emails, email{from, recipient, subject, body})
	return nil
}
func makeServer(t *testing.T) (*httptest.Server, *mockService, sqlmock.Sqlmock, *config.Config) {
	db, dbmock, err := sqlmock.New()
	assert.NoError(t, err)

	conf := &config.Config{AppSecret: "sup", Env: "debug"}
	svc := &mockService{}

	s, err := server.NewServer(
		&server.ServerConfig{
			Config:       conf,
			DB:           db,
			EmailService: svc,
			TemplatePath: "../../templates",
		},
	)
	assert.NoError(t, err)

	testServer := httptest.NewServer(s.Handler)
	conf.URL = testServer.URL

	return testServer, svc, dbmock, conf
}

func sendRequest(t *testing.T, path string, postBody []byte) (string, *http.Response) {
	var resp *http.Response
	var err error

	// We need a cookie jar so cookies are retained between redirects
	cookieJar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	assert.NoError(t, err)

	client := http.Client{Jar: cookieJar}

	if postBody == nil {
		resp, err = client.Get(path)
	} else {
		// TODO: switch this to client.PostForm to simplify
		resp, err = client.Post(path, "application/x-www-form-urlencoded", bytes.NewReader(postBody))
	}

	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	time.Sleep(1 * time.Millisecond) // Makes sure the resp object is populated properly
	return string(body), resp
}

// func getDbFields(thing interface{}) []string {
// 	dbFields := make([]string, 0)

// 	t := reflect.TypeOf(thing)

// 	for i := 0; i < t.NumField(); i++ {
// 		dbTag := t.Field(i).Tag.Get("db")
// 		if dbTag != "" {
// 			dbFields = append(dbFields, dbTag)
// 		}
// 	}

// 	return dbFields
// }
