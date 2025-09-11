package store

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/joshua-zingale/activity-grading-server/internal/activity"
)

type ActivityStore struct {
	activityMap map[string]activity.Activity
}

// Loads all activities contained in a file into the ActivityStore.
func NewActivityStore(activityFilePath string) ActivityStore {
	activity_data_file, err := os.Open(activityFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer activity_data_file.Close()

	activityStore := ActivityStore{make(map[string]activity.Activity)}

	scanner := bufio.NewScanner(activity_data_file)
	for scanner.Scan() {
		var activity activity.Activity
		json.Unmarshal(scanner.Bytes(), &activity)
		if activity.Validate() != nil {
			panic(err)
		}
		err = activityStore.Add(activity)
		if err != nil {
			panic(err)
		}
	}

	return activityStore
}

func (store *ActivityStore) Add(activity activity.Activity) error {
	_, exists := store.activityMap[activity.Identifier]
	if exists {
		return fmt.Errorf("cannot add activity with duplicate identifier '%s'", activity.Identifier)
	}
	store.activityMap[activity.Identifier] = activity
	return nil
}

func (store *ActivityStore) Get(identifier string) (activity.Activity, error) {
	activity, exists := store.activityMap[identifier]
	if !exists {
		return activity, fmt.Errorf("invalid activity identifier '%s'", activity.Identifier)
	}
	return activity, nil
}
