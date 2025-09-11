# Grader

Grader is an activity-grading web server server receives JSON-formatted submissions to activities and responds with graded feedback in JSON format.


## Compilation

Use `go build` to get the executable binary for the web server on your machine.

## Usage

To run the server with the demonstration activities, run

```bash
./grader -activities demo/activity-data.jsonl
```