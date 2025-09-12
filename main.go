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
	flag.Parse()

	activityStore := store.NewActivityStore(*activityFilePath)
	submissionHandler := handler.NewSubmissionHandler(&activityStore)
	go submissionHandler.StoreRecords()

	http.HandleFunc("POST /submissions", submissionHandler.Post)
	http.HandleFunc("OPTIONS /submissions", submissionHandler.Options)

	log.Printf("Activity Grading Server running at %s:%s", *webHost, *port)
	http.ListenAndServe(fmt.Sprintf("%s:%s", *webHost, *port), nil)
}
