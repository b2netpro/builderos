package models

import "time"

type JournalEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Text      string    `json:"text"`
	Tags      []string  `json:"tags,omitempty"`
}
