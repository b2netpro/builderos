package models

type Idea struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
	Feedback  string `json:"feedback,omitempty"`
}
