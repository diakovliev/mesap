package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/diakovliev/mesap/backend/controllers"
	"github.com/diakovliev/mesap/backend/fake_database"
)

const (
	//defaultListenAddressTLS  = ":8443"
	defaultListenAddress     = ":8080"
	defaultKeyFile           = ""
	defaultCertFile          = ""
	defaultStaticContent     = ""
	defaultStaticContentRoot = "/"
)

var (
	listenAddress *string
	certFile      *string
	keyFile       *string
	enableTls     *bool

	staticContent     *string
	staticContentRoot *string
)

func init() {

	staticContent = flag.String("static", defaultStaticContent, "Directory path with static content to serve")
	staticContentRoot = flag.String("root", defaultStaticContentRoot, "Root of the served static content folder")
	enableTls = flag.Bool("tls", false, "Enable TLS support")
	listenAddress = flag.String("listen", defaultListenAddress, "Listen address")
	certFile = flag.String("cert", defaultCertFile, "Server certificate")
	keyFile = flag.String("key", defaultKeyFile, "Server certificate key")

	flag.Parse()

	if *staticContent != "" {
		log.Printf("Static content directory: '%s'", *staticContent)
		log.Printf("Static content root: '%s'", *staticContentRoot)
	} else {
		log.Print("Static content: OFF")
	}

	log.Printf("Listen address: '%s'", *listenAddress)

	if !*enableTls {
		log.Println("TLS: OFF")
	} else {
		log.Println("TLS: ON")
		log.Printf("Certificate: '%s'", *certFile)
		log.Printf("Key: '%s'", *keyFile)
	}

}

func main() {

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Mount("/auth", controllers.NewAuthController(fake_database.NewDatabase()).Controller())
	})

	FileServer(r)

	if *enableTls {
		if err := http.ListenAndServeTLS(*listenAddress, *certFile, *keyFile, r); err != nil {
			log.Panicf("Fatal: %s", err)
		}
	} else {
		if err := http.ListenAndServe(*listenAddress, r); err != nil {
			log.Panicf("Fatal: %s", err)
		}
	}
}

// FileServer is serving static files.
func FileServer(router *chi.Mux) {
	root := *staticContent
	fs := http.FileServer(http.Dir(root))

	router.Get(*staticContentRoot+"*", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request URL: %s", r.RequestURI)
		if _, err := os.Stat(filepath.Join(root, r.RequestURI)); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}
