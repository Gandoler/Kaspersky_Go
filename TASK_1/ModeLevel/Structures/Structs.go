package Structures

type JobState string

const (
	StateQueued  JobState = "queued"
	StateRunning JobState = "running"
	StateDone    JobState = "done"
	StateFailed  JobState = "failed"
)

type EnqueueRequest struct {
	ID         string `json:"id"`
	Payload    string `json:"payload"`
	MaxRetries int    `json:"max_retries"`
}

type Job struct {
	ID         string
	Payload    string
	MaxRetries int
	Attempts   int
}

type JobStatus struct {
	State    JobState
	Attempts int
	Error    string
}
