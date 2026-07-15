// time.go
package timebox

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/b2netpro/builderos/internal/ideas"
)

type TimeEntry struct {
	IdeaID   int       `json:"idea_id"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Notes    string    `json:"notes"`
	Duration string    `json:"duration"`
}

const logFile = "time_logs.json"

func StartSession(ideaID int, notes string) error {
	fmt.Println("Session started. Press ENTER to stop.")
	start := time.Now()
	fmt.Scanln()
	end := time.Now()
	duration := end.Sub(start)

	entry := TimeEntry{
		IdeaID:   ideaID,
		Start:    start,
		End:      end,
		Notes:    notes,
		Duration: duration.String(),
	}

	logs, err := loadLogs()
	if err != nil {
		return err
	}
	logs = append(logs, entry)

	file, err := os.Create(logFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(logs)
}

func Report() error {
	logs, err := loadLogs()
	if err != nil {
		return err
	}

	ideaMap := make(map[int]string)
	storedIdeas, err := ideas.LoadIdeas()
	if err == nil {
		for _, idea := range storedIdeas {
			ideaMap[idea.ID] = idea.Text
		}
	}

	totals := make(map[int]time.Duration)
	for _, entry := range logs {
		dur, err := time.ParseDuration(entry.Duration)
		if err != nil {
			continue
		}
		totals[entry.IdeaID] += dur
	}

	fmt.Println("\n=== Time Report by Idea ===")
	ids := make([]int, 0, len(totals))
	for id := range totals {
		ids = append(ids, id)
	}
	sort.Ints(ids)
	for _, id := range ids {
		title := ideaMap[id]
		if title == "" {
			title = "(unknown idea)"
		}
		fmt.Printf("[%d] %s — %s total\n", id, title, totals[id])
	}

	return nil
}

func loadLogs() ([]TimeEntry, error) {
	var logs []TimeEntry

	file, err := os.Open(logFile)
	if err != nil {
		if os.IsNotExist(err) {
			return logs, nil
		}
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&logs)
	if err != nil {
		return nil, err
	}

	return logs, nil
}
