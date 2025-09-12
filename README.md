# Grader

Grader is an activity-grading web server that receives JSON-formatted submissions to activities and responds with graded feedback in JSON format. Grader is intended to facilitate the embedding of activities into web pages in a way that submission can be tracked, e.g. for research or graded school work.

## Compilation

Use `go build` in the repository to get the executable binary for the web server on your machine.

## Usage

To run the server with the demonstration activities, run

```bash
./grader -activities demo/activity-data.jsonl
```

To see all flag arguments available for this web server, use

```bash
./grader -h
```

The following command runs the web serving using TLS (for handling HTTPS connections) at 127.0.0.1 on port 80, assuming the existence of any referenced files:


```bash
./grader -tls -certificate certs/server.crt -private-key certs/server.key -activities demo/activity-data.jsonl -host 127.0.0.1 -port 80 
```

Logistical information is printed to standard error while the server is running and submissions are printed to standard output.

And the following use of the `curl` utility will get a response from the web server:

```bash
curl \
  --cacert certs/server.crt \
  https://127.0.0.1:80/submissions \
  -X POST \
  -H "Content-Type: application/json" \
  --data '{"identifier": "activity-2", "answer": "False", "session": "c01b760jf5d9s"}'
```

After sending the request, you should get a response that looks something like

```json
{"grade":1,"hint":"Oh yeah, this is the stuff."}
```

and the web server should have printed the submission, including a timestamp, to the standard output:

```json
{"identifier":"activity-2","answer":"False","session":"c01b760jf5d9s","timestamp":"2025-09-12T19:22:47.156046Z"}
```

## Schemata

Initialization of and communication with the web server both happen through JSON formatted data.

### Activity

```javascript
{
  identifier: string; // A unique identifier for the activity (e.g., 'question-1', 'activity-3').
  options: Option[]; // An array of possible answers for the activity. Must contain at least one Option with a grade of 1.0.
  hint?: string; // A general hint provided if a submission does not match any option.
}
```

The web server is initialized with `.jsonl` file containing all activities to be supported.
Specifically, the file must contain one `Activity` in JSON format per line.
See `demo/activity-data.jsonl` for an example.

### Option

```javascript
{
  answer: string; // The text of the answer.
  grade: number; // The grade for this answer (0.0 <= grade <= 1.0).
  hint?: string; // A specific hint to provide if this option is chosen.
}
```

Each `Option` represents a possible answer to an `Activity`.


### Submission
```javascript
{
  identifier: string; // The identifier of the activity being submitted.
  answer: string; // The user's submitted answer.
  session: string; // A unique session ID for the user.
}
```

A `Submission` is made to `/submissions` via a post request with 
JSON data and a `SubmissionFeedback` is sent back as a response.

### SubmissionFeedback
```javascript
{
  grade: number; // The evaluation for the submission (0.0 <= grade <= 1.0).
  hint: string; // A hint provided for the submitted answer, or an empty string if there is no hint.
}
```


## JavaScript

Inside `demo/grader.js` is a function, `gradeSubmission(...)` for working with the grading server from JavaScript. The function is documented above its definition.