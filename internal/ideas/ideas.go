// ideas.go
package ideas

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Idea struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	Feedback  string    `json:"feedback,omitempty"`
}

const dataFile = "ideas.json"

func SaveIdea(text string, feedback string) error {
	ideas, err := LoadIdeas()
	if err != nil {
		return err
	}

	idea := Idea{
		ID:        len(ideas) + 1,
		Text:      text,
		CreatedAt: time.Now(),
		Feedback:  feedback,
	}
	ideas = append(ideas, idea)

	file, err := os.Create(dataFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(ideas)
}

func ListIdeas() error {
	ideas, err := LoadIdeas()
	if err != nil {
		return err
	}

	for _, idea := range ideas {
		feedbackStatus := "❌ No Feedback"
		if idea.Feedback != "" {
			feedbackStatus = "✅ Feedback"
		}

		fmt.Printf("[%d] %s (added %s) %s\n", idea.ID, idea.Text, idea.CreatedAt.Format(time.RFC822), feedbackStatus)
	}
	return nil
}

func LoadIdeas() ([]Idea, error) {
	var ideas []Idea

	file, err := os.Open(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return ideas, nil // return empty slice
		}
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ideas)
	if err != nil {
		return nil, err
	}

	return ideas, nil
}
