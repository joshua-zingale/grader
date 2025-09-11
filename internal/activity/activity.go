package activity

import (
	"fmt"
)

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

type Submission struct {
	Identifier string `json:"identifier"`
	Answer     string `json:"answer"`
	Session    string `json:"session"`
}

type SubmissionFeedback struct {
	Grade float64 `json:"grade"`
	Hint  string  `json:"hint"`
}
