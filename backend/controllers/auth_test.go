package controllers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kong/go-srp"

	"github.com/diakovliev/mesap/backend/fake_database"
	"github.com/diakovliev/mesap/backend/ifaces"
)

var (
	testSalt     = []byte("test salt")
	testLogin    = []byte("bob")
	testPassword = []byte("1234245asdf")
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

func ensureResponse(t *testing.T, resp *http.Response, err error) {
	if err != nil {
		t.Fatalf("Login request error: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Bad status code: %d", resp.StatusCode)
	}
}

func TestRegisterLogin(t *testing.T) {

	testVerifier := srp.ComputeVerifier(SRP_PARAMS, testSalt, testLogin, testPassword)

	testDatabase := fake_database.NewDatabase()

	testServer := NewAuthTestServer(testDatabase)
	defer testServer.Close()

	testClient := testServer.NewClient("")

	registerData := RegisterRequestData{
		Login:    AuthEncodeBytes(testLogin),
		Salt:     AuthEncodeBytes(testSalt),
		Verifier: AuthEncodeBytes(testVerifier),
	}

	testClient._Post("register", AuthEncodeJson(registerData))

	srpClient := srp.NewClient(SRP_PARAMS, testSalt, testLogin, testPassword, srp.GenKey())

	loginData := LoginRequestData{
		Login:   registerData.Login,
		Secret1: AuthEncodeBytes(srpClient.ComputeA()),
	}

	resp, err := testClient._Post("login", AuthEncodeJson(loginData))
	ensureResponse(t, resp, err)

	loginResponse := AuthDecodeJson[LoginResponseData](resp.Body, func(err error) {
		t.Fatalf("Can't decode login responce! Error: %s", err)
	})
	if loginResponse == nil {
		return
	}

	t.Logf("Login response: %s", loginResponse.String())

	srpClient.SetB(AuthDecodeString(loginResponse.Secret2))

	login2Data := Login2RequestData{
		Server:  loginResponse.Server,
		Secret3: AuthEncodeBytes(srpClient.ComputeM1()),
	}

	resp, err = testClient._Post("login2", AuthEncodeJson(login2Data))
	ensureResponse(t, resp, err)

	login2Response := AuthDecodeJson[Login2ResponseData](resp.Body, func(err error) {
		t.Fatalf("Can't decode login responce! Error: %s", err)
	})
	if login2Response == nil {
		return
	}

	err = srpClient.CheckM2(AuthDecodeString(login2Response.Secret4))
	if err != nil {
		t.Fatalf("Client check M2 err: %s", err)
	}

	t.Logf("Client K: '%s'", AuthEncodeBytes(srpClient.ComputeK()))
}
