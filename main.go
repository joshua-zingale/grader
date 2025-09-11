package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var activityStore ActivityStore
var recordChan = make(chan SubmissionRecord, 100)

func main() {

	activityFilePath := flag.String("activities", "activity-data.jsonl", "path to a .jsonl file containing the activity data.")
	webHost := flag.String("host", "127.0.0.1", "the host for this web server")
	port := flag.String("port", "8080", "the port for this web server")
	flag.Parse()

	activityStore = loadActivities(*activityFilePath)

	go storeRecords(recordChan)
	http.HandleFunc("POST /submissions", postSubmission)
	log.Printf("Activity Grading Server running at %s:%s", *webHost, *port)
	http.ListenAndServe(fmt.Sprintf("%s:%s", *webHost, *port), nil)
}

func loadActivities(activityFilePath string) ActivityStore {
	activity_data_file, err := os.Open(activityFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer activity_data_file.Close()

	store := NewActivityStore()

	scanner := bufio.NewScanner(activity_data_file)
	for scanner.Scan() {
		var activity Activity
		json.Unmarshal(scanner.Bytes(), &activity)
		if activity.Validate() != nil {
			panic(err)
		}
		err = store.Add(activity)
		if err != nil {
			panic(err)
		}
	}

	return store
}

func storeRecords(records <-chan SubmissionRecord) {
	for record := range records {
		recordJson, _ := json.Marshal(record)
		fmt.Printf("%s\n", recordJson)
	}
}

func postSubmission(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsuported media type: Expected Content-Type: application/json", http.StatusUnsupportedMediaType)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var submission Submission
	json.NewDecoder(req.Body).Decode(&submission)

	activity, err := activityStore.Get(submission.Identifier)
	if err != nil {
		http.Error(w, "Invalid activity identifier.", 404)
		return
	}

	recordChan <- SubmissionRecord{
		submission,
		time.Now().UTC(),
	}

	feedback := activity.Grade(submission)
	feedbackJson, _ := json.Marshal(feedback)

	w.Write(feedbackJson)
}

type Submission struct {
	Identifier string `json:"identifier"`
	Answer     string `json:"answer"`
	Session    string `json:"session"`
}

type SubmissionRecord struct {
	Submission
	Timestamp time.Time `json:"timestamp"`
}

type SubmissionFeedback struct {
	Grade float64 `json:"grade"`
	Hint  string  `json:"hint"`
}

type Activity struct {
	Identifier string   `json:"identifier"`
	Options    []Option `json:"options"`
	Hint       string   `json:"hint"`
}

func (activity Activity) Validate() error {
	if activity.Identifier == "" {
		return fmt.Errorf("missing identifier for activity")
	}

	maximum_grade := 0.0
	minimum_grade := 1.0
	num_options := 0
	for _, option := range activity.Options {
		maximum_grade = max(maximum_grade, option.Grade)
		minimum_grade = min(minimum_grade, option.Grade)
		num_options += 1
		if maximum_grade > 1 || minimum_grade < 0 {
			return fmt.Errorf("activity '%s' cannot have a grade of %f", activity.Identifier, option.Grade)
		}
	}

	if num_options == 0 {
		return fmt.Errorf("there must be at least one option for activity '%s'", activity.Identifier)
	}
	if maximum_grade != 1.0 {
		return fmt.Errorf("there must be at least one option with a grade of 1.0 for activity '%s'", activity.Identifier)
	}

	return nil

}

func (activity Activity) Grade(submission Submission) SubmissionFeedback {
	for _, option := range activity.Options {
		if option.Answer == submission.Answer {
			return SubmissionFeedback{
				Grade: option.Grade,
				Hint:  option.Hint,
			}
		}
	}

	return SubmissionFeedback{
		Grade: 0.0,
		Hint:  activity.Hint,
	}
}

type Option struct {
	Answer string  `json:"answer"`
	Grade  float64 `json:"grade"`
	Hint   string  `json:"hint"`
}

type ActivityStore struct {
	activityMap map[string]Activity
}

func NewActivityStore() ActivityStore {
	return ActivityStore{make(map[string]Activity)}
}

func (store *ActivityStore) Add(activity Activity) error {
	_, exists := store.activityMap[activity.Identifier]
	if exists {
		return fmt.Errorf("cannot add activity with duplicate identifier '%s'", activity.Identifier)
	}
	store.activityMap[activity.Identifier] = activity
	return nil
}

func (store *ActivityStore) Get(identifier string) (Activity, error) {
	activity, exists := store.activityMap[identifier]
	if !exists {
		return activity, fmt.Errorf("invalid activity identifier '%s'", activity.Identifier)
	}
	return activity, nil
}
