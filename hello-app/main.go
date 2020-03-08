package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	ver = "v1.0.0"
)

func main() {
	// register hello function to handle all requests
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/version", version)

	// use PORT environment variable, or default to 8080
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	svr := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	serveAndGracefulShutdown(svr)
}

// hello responds with a message.
func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	msg := os.Getenv("MESSAGE")
	if len(msg) == 0 {
		msg = "Bienvini nan mond teknologi a!"
	}
	fmt.Fprintf(w, "Álo tout moun!\n")
	fmt.Fprintf(w, "Version: %s\n", ver)
	fmt.Fprintf(w, "Message: %s\n", msg)
}

// endpoint to test the health of the app
func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

// version returns the service version
func version(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, ver)
}

// Starts server with gracefully shutdown semantics
func serveAndGracefulShutdown(svr *http.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// wait for requests and serve
	serveAndWait := make(chan error)
	go func() {
		log.Printf("Server listening on port %s", svr.Addr)
		serveAndWait <- svr.ListenAndServe()
	}()

	// block until either an error or OS-level signals
	// to shutdown gracefully
	select {
	case err := <-serveAndWait:
		log.Fatal(err)
	case <-sigChan:
		log.Printf("Shutdown signal received... closing server")
		svr.Shutdown(context.TODO())
	}
}
