package journal

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type JournalEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Text      string    `json:"text"`
	Tags      []string  `json:"tags,omitempty"`
}

const journalFile = "journal.json"

func AddEntry(text string, tagString string) error {
	tags := []string{}
	if tagString != "" {
		tags = strings.Split(tagString, ",")
	}

	entry := JournalEntry{
		Timestamp: time.Now(),
		Text:      text,
		Tags:      tags,
	}

	entries, _ := LoadEntries()
	entries = append(entries, entry)

	file, err := os.Create(journalFile)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	return enc.Encode(entries)
}

func LoadEntries() ([]JournalEntry, error) {
	var entries []JournalEntry
	data, err := os.ReadFile(journalFile)
	if err != nil {
		return entries, nil
	}
	json.Unmarshal(data, &entries)
	return entries, nil
}

func ListEntries(filterLastDays int, markdown bool) error {
	entries, err := LoadEntries()
	if err != nil {
		return err
	}

	cutoff := time.Time{}
	if filterLastDays > 0 {
		cutoff = time.Now().AddDate(0, 0, -filterLastDays)
	}

	for _, entry := range entries {
		if !cutoff.IsZero() && entry.Timestamp.Before(cutoff) {
			continue
		}

		if markdown {
			fmt.Printf("### %s\n\n%s\n\n", entry.Timestamp.Format("2006-01-02"), entry.Text)
			if len(entry.Tags) > 0 {
				fmt.Printf("_Tags: %s_\n\n", strings.Join(entry.Tags, ", "))
			}
		} else {
			fmt.Printf("[%s] %s", entry.Timestamp.Format(time.RFC822), entry.Text)
			if len(entry.Tags) > 0 {
				fmt.Printf(" — Tags: %v", entry.Tags)
			}
			fmt.Println()
		}
	}

	return nil
}
