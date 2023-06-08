package dto

type RunDto struct {
	ID        uint   `json:"id"`
	JobID     uint   `json:"job_id"`
	Timestamp string `json:"timestamp"`
	Container string `json:"container"`
	Output    string `json:"output"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}
