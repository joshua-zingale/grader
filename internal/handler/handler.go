package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/joshua-zingale/activity-grading-server/internal/activity"
	"github.com/joshua-zingale/activity-grading-server/internal/store"
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

func (h *SubmissionHandler) Post(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsuported media type: Expected Content-Type: application/json", http.StatusUnsupportedMediaType)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var submission activity.Submission
	json.NewDecoder(req.Body).Decode(&submission)

	activity, err := h.activityStore.Get(submission.Identifier)
	if err != nil {
		http.Error(w, "Invalid activity identifier.", 404)
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
