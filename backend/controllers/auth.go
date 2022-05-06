package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Auth struct {
}

type RegisterRequestData struct {
	Login   string
	Mail    string
	Secret0 string
}

func (rrd *RegisterRequestData) String() string {
	return fmt.Sprintf("Login: '%s' Mail: '%s' Secret0: '%s'", rrd.Login, rrd.Mail, rrd.Secret0)
}

type RegisterResponseData struct {
	Secret1 string
}

type LoginRequestData struct {
	Login   string
	Secret2 string
}

type LoginResponseData struct {
	Token string
}

type LogoutRequestData struct {
	Token string
}

func NewAuthController() *Auth {
	return &Auth{}
}

func (a *Auth) Controller() chi.Router {
	r := chi.NewRouter()
	r.Post("/register", a.PostRegister)
	r.Post("/login", a.PostLogin)
	r.Post("/logout", a.PostLogout)
	return r
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

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(responseData); err != nil {
		log.Printf("Login response encoding error: %s", err)
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
