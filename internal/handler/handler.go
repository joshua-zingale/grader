package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joshua-zingale/grader/internal/activity"
	"github.com/joshua-zingale/grader/internal/store"
)

type SubmissionHandler struct {
	activityStore *store.ActivityStore
	recordChan    chan SubmissionRecord
}

func NewSubmissionHandler(activityStore *store.ActivityStore) SubmissionHandler {
	return SubmissionHandler{
		activityStore: activityStore,
		recordChan:    make(chan SubmissionRecord, 100),
	}
}

func setCorsHeaders(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func (h *SubmissionHandler) Options(w http.ResponseWriter, req *http.Request) {
	setCorsHeaders(&w, req)
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.WriteHeader(http.StatusOK)
}

func (h *SubmissionHandler) Post(w http.ResponseWriter, req *http.Request) {
	setCorsHeaders(&w, req)
	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsuported media type: Expected Content-Type: application/json", http.StatusUnsupportedMediaType)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var submission activity.Submission
	json.NewDecoder(req.Body).Decode(&submission)

	activity, err := h.activityStore.Get(submission.Identifier)
	if err != nil {
		http.Error(w, "Invalid activity identifier", 404)
		log.Printf("Invalid activity identifier '%s'", activity.Identifier)
		return
	}

	go func() {
		h.recordChan <- SubmissionRecord{
			Submission: submission,
			Timestamp:  time.Now().UTC(),
		}
	}()

	feedback := activity.Grade(submission)
	feedbackJson, _ := json.Marshal(feedback)

	w.Write(feedbackJson)
}

// A non-terminating process that logs all submissions.
func (h *SubmissionHandler) StoreRecords() {
	for record := range h.recordChan {
		recordJson, _ := json.Marshal(record)
		fmt.Printf("%s\n", recordJson)
	}
}

type SubmissionRecord struct {
	activity.Submission
	Timestamp time.Time `json:"timestamp"`
}
