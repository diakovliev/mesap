package controllers

import (
	"crypto"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/kong/go-srp"

	"github.com/diakovliev/mesap/backend/ifaces"
	"github.com/diakovliev/mesap/backend/models"
)

const (
	HASH = crypto.SHA1
)

var (
	SRP_PARAMS = srp.GetParams(4096)
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
	Login    string `json: "login"`
	Salt     string `json: "salt"`
	Verifier string `json: "verifier"`

	// TODO: User extended info
}

func (rrd *RegisterRequestData) String() string {
	return fmt.Sprintf("Login: '%s' Salt: '%s' Verifier: '%s'", rrd.Login, rrd.Salt, rrd.Verifier)
}

type RegisterResponseData struct {
	UserId models.IdData
}

func (rrd *RegisterResponseData) String() string {
	return fmt.Sprintf("UserId: '%d'", rrd.UserId)
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
	Secret2 string
}

func (lrd *LoginResponseData) String() string {
	return fmt.Sprintf("Server: '%s' Secret2: '%s'", lrd.Server, lrd.Secret2)
}

type Login2RequestData struct {
	Server  string
	Secret3 string
}

func (l2rd *Login2RequestData) String() string {
	return fmt.Sprintf("Server: '%s' Secret3: '%s'", l2rd.Server, l2rd.Secret3)
}

type Login2ResponseData struct {
	Secret4 string
}

func (l2r *Login2ResponseData) String() string {
	return fmt.Sprintf("Secret4: '%s'", l2r.Secret4)
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

	checksum := HASH.New()
	checksum.Write([]byte(content))

	key := base64.StdEncoding.EncodeToString(checksum.Sum(nil))
	a.servers[key] = AuthServer{
		server: content,
		user:   record,
	}
	return key
}

func encodeServer(key []byte, verifier []byte, A []byte) string {
	return fmt.Sprintf(
		"%s:%s:%s",
		AuthEncodeHexBytes(key),
		AuthEncodeHexBytes(verifier),
		AuthEncodeHexBytes(A),
	)
}

func decodeServer(input string) *srp.SRPServer {
	e := strings.Split(input, ":")
	if len(e) < 3 {
		panic("Not expected elements count!")
	}

	key := AuthDecodeHexString(e[0])
	verifier := AuthDecodeHexString(e[1])
	A := AuthDecodeHexString(e[2])

	srv := srp.NewServer(SRP_PARAMS, verifier, key)
	srv.SetA(A)

	return srv
}

func LogStringChecksum(name string, content string) {
	checksum := HASH.New()
	checksum.Write([]byte(content))

	log.Printf("%s: %x", name, checksum.Sum(nil))
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

	var responseData RegisterResponseData

	requestData := AuthDecodeJson[RegisterRequestData](r.Body, func(err error) {
		log.Printf("Register request decoding error: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	})
	if requestData == nil {
		return
	}

	users, err := a.db.Users()
	if err != nil {
		log.Printf("Can't access to 'users' table: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, err = users.Find(func(record models.User) bool {
		return record.Login == requestData.Login
	})
	if err != nil && (err != ifaces.ErrEmptyTable && err != ifaces.ErrNoSuchRecord) {
		log.Printf("User with login '%s' already registered!", requestData.Login)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	log.Printf("Register data: %s", requestData.String())

	// Table contains base64 encoded data
	user := models.User{
		Login:    requestData.Login,
		Salt:     requestData.Salt,
		Verifier: requestData.Verifier,
	}

	log.Printf("Register user: Login: '%s' Salt: '%s' Verifier: '%s'", user.Login, user.Salt, user.Verifier)

	userId, err := users.Insert(user)
	if err != nil {
		log.Printf("Can't insert data into database: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	responseData.UserId = userId

	log.Printf("Register response: %s", responseData.String())

	AuthEncodeAndWriteJson(w, responseData)
}

func (a *Auth) PostLogin(w http.ResponseWriter, r *http.Request) {

	var responseData LoginResponseData

	requestData := AuthDecodeJson[LoginRequestData](r.Body, func(err error) {
		log.Printf("Login request decoding error: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	})
	if requestData == nil {
		return
	}

	log.Printf("Login data: %s", requestData.String())

	users, err := a.db.Users()
	if err != nil {
		log.Printf("Can't access to 'users' table: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	record, err := users.Find(func(u models.User) bool {
		return u.Login == requestData.Login
	})
	if err != nil {
		log.Printf("Can't find user record: %s", err)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	log.Printf("record Salt: '%s' Verifier: '%s'", record.Salt, record.Verifier)

	//salt := AuthDecodeString(record.Salt)
	verifier := AuthDecodeString(record.Verifier)

	//log.Printf("salt: '%s' verifier: '%s'", salt, verifier)

	key := srp.GenKey()
	A := AuthDecodeString(requestData.Secret1)

	srv := srp.NewServer(SRP_PARAMS, verifier, key)

	responseData.Server = a.addServer(encodeServer(key, verifier, A), record)
	responseData.Secret2 = AuthEncodeBytes(srv.ComputeB())

	log.Printf("Login response: %s", responseData.String())

	AuthEncodeAndWriteJson(w, responseData)
}

func (a *Auth) PostLogin2(w http.ResponseWriter, r *http.Request) {

	var responseData Login2ResponseData

	requestData := AuthDecodeJson[Login2RequestData](r.Body, func(err error) {
		log.Printf("Login request decoding error: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	})
	if requestData == nil {
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

	//log.Printf("Server to decode: %s", server.server)

	srv := decodeServer(server.server)

	serverM2, err := srv.CheckM1(AuthDecodeString(requestData.Secret3))
	if err != nil {
		log.Printf("Server M1 err: %s", err)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	responseData.Secret4 = AuthEncodeBytes(serverM2)

	log.Printf("Login2 response: %s", responseData.String())

	log.Printf("Server K: '%s'", AuthEncodeBytes(srv.ComputeK()))

	AuthEncodeAndWriteJson(w, responseData)
}

func (a *Auth) PostLogout(w http.ResponseWriter, r *http.Request) {
	// 	var requestData LogoutRequestData

	// 	decoder := json.NewDecoder(r.Body)
	// 	if err := decoder.Decode(&requestData); err != nil {
	// 		log.Printf("Logout request decoding error: %s", err)
	// 		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	// 		return
	// 	}

	// 	// TODO: Logout
}
