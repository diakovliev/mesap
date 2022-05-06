package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	testToken = "ZZ.YY.XX"
)

type AuthTestServer struct {
	r  chi.Router
	a  *Auth
	ts *httptest.Server
}

func NewAuthTestServer() *AuthTestServer {
	ret := AuthTestServer{
		r: chi.NewRouter(),
		a: NewAuthController(),
	}
	ret.r.Use(middleware.Logger)
	ret.r.Mount("/", ret.a.Controller())
	ret.ts = httptest.NewServer(ret.r)
	return &ret
}

func (s *AuthTestServer) Close() {
	s.ts.Close()
}

type TestTransport struct {
	Token  string
	Server *AuthTestServer
}

func (t *TestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if len(t.Token) > 0 {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.Token))
	}
	return http.DefaultTransport.RoundTrip(req)
}

func (t *TestTransport) _Do(method string, query string, body io.Reader) (*http.Response, error) {
	client := &http.Client{Transport: t}
	q := fmt.Sprintf("%s/%s", t.Server.ts.URL, query)
	req, err := http.NewRequest(method, q, body)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

func (t *TestTransport) _Get(query string) (*http.Response, error) {
	return t._Do(http.MethodGet, query, nil)
}

func (t *TestTransport) _Post(query string, body io.Reader) (*http.Response, error) {
	return t._Do(http.MethodPost, query, body)
}

func (s *AuthTestServer) NewClient(token string) *TestTransport {
	return &TestTransport{Token: token, Server: s}
}

func TestRegister(t *testing.T) {
	testServer := NewAuthTestServer()
	defer testServer.Close()

	registerData := RegisterRequestData{
		Login:   "Test login",
		Mail:    "test@mail.com",
		Secret0: "1234567656",
	}

	body := bytes.NewBufferString("")
	encoder := json.NewEncoder(body)
	if err := encoder.Encode(registerData); err != nil {
		t.Fatalf("Request encoding error: %s", err)
	}

	testClient := testServer.NewClient("")
	testClient._Post("register", body)
}
