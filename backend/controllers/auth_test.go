package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/opencoff/go-srp"

	"github.com/diakovliev/mesap/backend/fake_database"
	"github.com/diakovliev/mesap/backend/ifaces"
)

var (
	testToken = "ZZ.YY.XX"
)

type AuthTestServer struct {
	r  chi.Router
	a  *Auth
	ts *httptest.Server
}

func NewAuthTestServer(db ifaces.Database) *AuthTestServer {
	ret := AuthTestServer{
		r: chi.NewRouter(),
		a: NewAuthController(db),
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

func TestRegisterLogin(t *testing.T) {
	testDatabase := fake_database.NewDatabase()

	testServer := NewAuthTestServer(testDatabase)
	defer testServer.Close()

	testClient := testServer.NewClient("")

	registerData := RegisterRequestData{
		Login:   "Test login",
		Mail:    "test@mail.com",
		Secret0: "1234567656",
	}

	body := bytes.NewBufferString("")
	encoder := json.NewEncoder(body)
	if err := encoder.Encode(registerData); err != nil {
		t.Fatalf("Register request encoding error: %s", err)
	}

	testClient._Post("register", body)

	s, err := srp.New(N_BITS)
	if err != nil {
		t.Fatalf("Can't create srp instance! Error: %s", err)
	}

	c, err := s.NewClient([]byte(registerData.Login), []byte(registerData.Secret0))
	if err != nil {
		t.Fatalf("Can't create srp client instance! Error: %s", err)
	}

	loginData := LoginRequestData{
		Login:   registerData.Login,
		Secret1: c.Credentials(),
	}

	body = bytes.NewBufferString("")
	encoder = json.NewEncoder(body)
	if err := encoder.Encode(loginData); err != nil {
		t.Fatalf("Login request encoding error: %s", err)
	}

	resp, err := testClient._Post("login", body)
	if err != nil {
		t.Fatalf("Login request error: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Bad status code: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var loginResponse LoginResponseData

	err = decoder.Decode(&loginResponse)
	if err != nil {
		t.Fatalf("Can't decode login responce! Error: %s", err)
	}

	t.Logf("Login response: %s", loginResponse.String())

	clientAuth, err := c.Generate(loginResponse.Secret2)
	if err != nil {
		t.Fatalf("Can't calculate client auth! Error: %s", err)
	}

	login2Data := Login2RequestData{
		//Secret1: loginResponse.Secret1,
		//Secret2: loginResponse.Secret2,
		Server:  loginResponse.Server,
		Secret3: clientAuth,
	}

	body = bytes.NewBufferString("")
	encoder = json.NewEncoder(body)
	if err := encoder.Encode(login2Data); err != nil {
		t.Fatalf("Login2 request encoding error: %s", err)
	}

	resp, err = testClient._Post("login2", body)
	if err != nil {
		t.Fatalf("Login2 request error: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Login2 bad status code: %d", resp.StatusCode)
	}

	decoder = json.NewDecoder(resp.Body)
	var login2Response Login2ResponseData

	err = decoder.Decode(&login2Response)
	if err != nil {
		t.Fatalf("Can't decode login2 response! Error: %s", err)
	}

	t.Logf("Login2 response: %s", login2Response.String())

	if !c.ServerOk(login2Response.Token) {
		t.Fatalf("Sever proof validation error!")
	}

	clientKey := base64.StdEncoding.EncodeToString(c.RawKey())
	t.Logf("client key: %s", clientKey)
}
