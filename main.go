package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	root, address, useTLS := parseArgs()
	shutdown := serve(root, address, useTLS)
	wait()
	shutdown()
}

// wait blocks and waits for incoming OS signals so we can perform
// a graceful server shutdown.
func wait() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Kill, os.Interrupt)
	log.Println("received signal:", <-signals)
}

// serve fires up the server.
// Returns a function which performs a graceful shutdown when called.
func serve(root, address string, useTLS bool) func() {
	server := http.Server{
		Addr:         address,
		Handler:      logger(http.FileServer(http.Dir(root))),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	ln, err := Listen(address)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("serving %s", root)
	log.Printf("listening on %s", ln.Addr())

	if !useTLS {
		go func() {
			if err := server.Serve(ln); err != nil {
				log.Fatal(err)
			}
		}()
	} else {
		certCache := filepath.Join(os.TempDir(), "srv-certs")
		certManager := autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache(certCache),
		}

		server.TLSConfig = &tls.Config{
			PreferServerCipherSuites: true,
			NextProtos:               []string{"h2", "http/1.1"},
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP256, tls.X25519},
			GetCertificate:           certManager.GetCertificate,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		}

		go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
		go func() {
			if err := server.ServeTLS(ln, "", ""); err != nil {
				log.Fatal(err)
			}
		}()
	}

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		if err := server.Shutdown(ctx); err != nil {
			log.Println(err)
		}
		cancel()
	}
}

func logger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

// parseArgs parses command line arguments and returns the serve root directory,
// listener address and whether or not to use TLS.
func parseArgs() (string, string, bool) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options] <directory>\n", os.Args[0])
		flag.PrintDefaults()
	}

	useTLS := flag.Bool("tls", false, "Use TLS and Let's Encrypt.")
	address := flag.String("addr", "", "address to listen on.")
	version := flag.Bool("version", false, "Display version information.")
	flag.Parse()

	if *version {
		fmt.Fprintln(os.Stderr, Version())
		os.Exit(0)
	}

	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	if len(*address) == 0 {
		*address = ":0"
	}

	*address = strings.ToLower(*address)
	if *useTLS && strings.HasPrefix(*address, "http://") {
		*address = "https://" + (*address)[7:]
	}

	return root, *address, *useTLS
}
