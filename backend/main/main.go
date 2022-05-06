package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/diakovliev/mesap/backend/controllers"
)

var (
	defaultListenAddressTLS = ":8443"
	defaultListenAddress    = ":8080"
	listenAddress           = defaultListenAddress
	certFile                = ""
	keyFile                 = ""
	enableTls               = false

	// Controllers
	//auth *controllers.Auth = nil
)

func init() {

	enableTls = *flag.Bool("tls", enableTls, "Enable TLS support")
	if !enableTls {
		log.Println("TLS: OFF")

		listenAddress = *flag.String("listen", defaultListenAddress, "Listen address")
		log.Printf("Listen address: %s", listenAddress)

	} else {
		log.Println("TLS: ON")

		certFile = *flag.String("cert", certFile, "Server certificate")
		log.Printf("Certificate: %s", certFile)

		keyFile = *flag.String("key", keyFile, "Server certificate key")
		log.Printf("Key: %s", keyFile)

		listenAddress = *flag.String("listen", defaultListenAddressTLS, "Listen address")
		log.Printf("Listen address: %s", listenAddress)
	}
}

func main() {

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	auth := controllers.NewAuthController()
	r.Mount("/auth", auth.Controller())

	if enableTls {
		if err := http.ListenAndServeTLS(listenAddress, certFile, keyFile, r); err != nil {
			log.Panicf("Fatal: %s", err)
		}
	} else {
		if err := http.ListenAndServe(listenAddress, r); err != nil {
			log.Panicf("Fatal: %s", err)
		}
	}
}
