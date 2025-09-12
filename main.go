package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/joshua-zingale/grader/internal/handler"
	"github.com/joshua-zingale/grader/internal/store"
)

func main() {

	activityFilePath := flag.String("activities", "activity-data.jsonl", "path to a .jsonl file containing the activity data.")
	webHost := flag.String("host", "127.0.0.1", "the host for this web server")
	port := flag.String("port", "8080", "the port for this web server")
	useTls := flag.Bool("tls", false, "if set, uses TLS")
	certFilePath := flag.String("certificate", "server.crt", "the TLS certificate")
	privateKeyFilePath := flag.String("private-key", "server.key", "the TLS private key")

	flag.Parse()

	activityStore := store.NewActivityStore(*activityFilePath)
	submissionHandler := handler.NewSubmissionHandler(&activityStore)
	go submissionHandler.StoreRecords()

	http.HandleFunc("POST /submissions", submissionHandler.Post)
	http.HandleFunc("OPTIONS /submissions", submissionHandler.Options)

	url := fmt.Sprintf("%s:%s", *webHost, *port)

	var err error
	if *useTls {
		log.Printf("Starting HTTPS server at %s", url)
		err = http.ListenAndServeTLS(url, *certFilePath, *privateKeyFilePath, nil)
	} else {
		log.Printf("Starting HTTP server at %s", url)
		err = http.ListenAndServe(url, nil)
	}

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
