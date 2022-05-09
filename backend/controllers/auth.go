package controllers

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/opencoff/go-srp"

	"github.com/diakovliev/mesap/backend/ifaces"
	"github.com/diakovliev/mesap/backend/models"
)

const (
	N_BITS = 2048
)

type AuthServer struct {
	server string
	user   models.User
}

type Auth struct {
	sync.Mutex
	db      ifaces.Database
	servers map[string]AuthServer
}

type RegisterRequestData struct {
	Login   string
	Mail    string
	Secret0 string // User password
}

func (rrd *RegisterRequestData) String() string {
	return fmt.Sprintf("Login: '%s' Mail: '%s' Secret0: '%s'", rrd.Login, rrd.Mail, rrd.Secret0)
}

type RegisterResponseData struct {
	UserId models.IdData
}

type LoginRequestData struct {
	Login   string
	Secret1 string
}

func (lrd *LoginRequestData) String() string {
	return fmt.Sprintf("Login: '%s' Secret1: '%s'", lrd.Login, lrd.Secret1)
}

type LoginResponseData struct {
	Server  string
	Secret1 string
	Secret2 string
}

func (lrd *LoginResponseData) String() string {
	return fmt.Sprintf("Server: '%s' Secret1: '%s' Secret2: '%s'", lrd.Server, lrd.Secret1, lrd.Secret2)
}

type Login2RequestData struct {
	Server  string
	Secret1 string
	Secret2 string
	Secret3 string
}

func (l2rd *Login2RequestData) String() string {
	return fmt.Sprintf("Server: '%s' Secret1: '%s' Secret2: '%s' Secret3: '%s'", l2rd.Server, l2rd.Secret1, l2rd.Secret2, l2rd.Secret3)
}

type Login2ResponseData struct {
	Token string
}

func (l2r *Login2ResponseData) String() string {
	return fmt.Sprintf("Token: '%s'", l2r.Token)
}

type LogoutRequestData struct {
	Token string
}

func NewAuthController(db ifaces.Database) *Auth {
	return &Auth{db: db, servers: make(map[string]AuthServer)}
}

func (a *Auth) Controller() chi.Router {
	r := chi.NewRouter()
	r.Post("/register", a.PostRegister)
	r.Post("/login", a.PostLogin)
	r.Post("/login2", a.PostLogin2)
	r.Post("/logout", a.PostLogout)
	return r
}

func (a *Auth) addServer(content string, record models.User) string {
	a.Lock()
	defer a.Unlock()

	checksum := sha1.New()
	checksum.Write([]byte(content))

	key := base64.StdEncoding.EncodeToString(checksum.Sum(nil))
	a.servers[key] = AuthServer{
		server: content,
		user:   record,
	}
	return key
}

func (a *Auth) forgetServer(key string) {
	a.Lock()
	defer a.Unlock()

	_, ok := a.servers[key]
	if !ok {
		return
	}

	delete(a.servers, key)
}

func (a *Auth) getServer(key string) (AuthServer, bool) {
	a.Lock()
	defer a.Unlock()

	content, ok := a.servers[key]
	return content, ok
}

func (a *Auth) PostRegister(w http.ResponseWriter, r *http.Request) {

	var requestData RegisterRequestData
	var responseData RegisterResponseData

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestData); err != nil {
		log.Printf("Register request decoding error: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// TODO: Register
	log.Printf("Register data: %s", requestData.String())

	s, err := srp.New(N_BITS)
	if err != nil {
		log.Printf("Srp instance not created. Error: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	v, err := s.Verifier([]byte(requestData.Login), []byte(requestData.Secret0))
	if err != nil {
		log.Printf("Verifier not created. Error: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	oi, ov := v.Encode()

	//log.Printf("V: i: %s v: %s", oi, ov)

	users, err := a.db.Users()
	if err != nil {
		log.Printf("Can't access to 'users' table: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	userId, err := users.Insert(models.User{
		Login: requestData.Login,
		Mail:  requestData.Mail,
		Ii:    oi,
		Iv:    ov,
	})
	if err != nil {
		log.Printf("Can't insert data into database: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	responseData.UserId = userId

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(responseData); err != nil {
		log.Printf("Register response encoding error: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (a *Auth) PostLogin(w http.ResponseWriter, r *http.Request) {

	var requestData LoginRequestData
	var responseData LoginResponseData

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestData); err != nil {
		log.Printf("Login request decoding error: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// TODO: Login
	log.Printf("Login data: %s", requestData.String())

	_, A, err := srp.ServerBegin(requestData.Secret1)

	users, err := a.db.Users()
	if err != nil {
		log.Printf("Can't access to 'users' table: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//log.Printf("id: %s", id)

	record, err := users.Find(func(u models.User) bool {
		return u.Login == requestData.Login
	})
	if err != nil {
		log.Printf("Can't find user record: %s", err)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	//log.Printf("record.Iv: %s", record.Iv)

	s, v, err := srp.MakeSRPVerifier(record.Iv)
	if err != nil {
		log.Printf("Can't make srp verifier: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	srv, err := s.NewServer(v, A)
	if err != nil {
		log.Printf("Can't create srp server instance: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//responseData.Secret1 = requestData.Secret1
	responseData.Secret2 = srv.Credentials()
	responseData.Server = a.addServer(srv.Marshal(), record)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(responseData); err != nil {
		log.Printf("Login response encoding error: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (a *Auth) PostLogin2(w http.ResponseWriter, r *http.Request) {

	var requestData Login2RequestData
	var responseData Login2ResponseData

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestData); err != nil {
		log.Printf("Login2 request decoding error: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	log.Printf("Login2 data: %s", requestData.String())

	server, ok := a.getServer(requestData.Server)
	if !ok {
		log.Printf("Unknown srp server id: %s", requestData.Server)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// Forget SRP context
	defer a.forgetServer(requestData.Server)

	srv, err := srp.UnmarshalServer(server.server)
	if err != nil {
		log.Printf("Can't unmarshal srp server instance: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	proof, ok := srv.ClientOk(requestData.Secret3)
	if !ok {
		log.Printf("Authentication failed!")
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// TODO: Create session for user. Autorized user info in server.user!
	responseData.Token = proof

	serverKey := base64.StdEncoding.EncodeToString(srv.RawKey())
	log.Printf("server key: %s", serverKey)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(responseData); err != nil {
		log.Printf("Login2 response encoding error: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (a *Auth) PostLogout(w http.ResponseWriter, r *http.Request) {
	var requestData LogoutRequestData

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestData); err != nil {
		log.Printf("Logout request decoding error: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// TODO: Logout
}
