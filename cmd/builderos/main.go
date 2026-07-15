// main.go
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/b2netpro/builderos/internal/ideas"
	"github.com/b2netpro/builderos/internal/journal"
	"github.com/b2netpro/builderos/internal/timebox"
	"github.com/sashabaranov/go-openai"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: builderos [command]")
		return
	}

	switch os.Args[1] {
	case "idea:add":
		handleIdeaAdd(os.Args[2:])
	case "idea:feedback":
		if len(os.Args) < 3 {
			fmt.Println("Usage: idea:feedback <id>")
			return
		}
		id, _ := strconv.Atoi(os.Args[2])
		allIdeas, _ := ideas.LoadIdeas()
		for _, idea := range allIdeas {
			if idea.ID == id {
				fmt.Printf("Feedback for Idea [%d]:\n%s\n", id, idea.Feedback)
				return
			}
		}
		fmt.Println("Idea not found.")
	case "idea:list":
		handleIdeaList()
	case "time:start":
		if len(os.Args) < 3 {
			fmt.Println("Usage: builderos time:start [ideaID] \"optional note\"")
			return
		}
		ideaID, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid idea ID")
			return
		}
		fmt.Println(ideas.LoadIdeas())
		note := ""
		if len(os.Args) > 3 {
			note = os.Args[3]
		}
		handleTimeStart(ideaID, note)
	case "journal:add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: journal:add \"your note\" [--tags tag1,tag2]")
			return
		}
		note := os.Args[2]
		tags := ""
		if len(os.Args) > 4 && os.Args[3] == "--tags" {
			tags = os.Args[4]
		}
		err := journal.AddEntry(note, tags)
		if err != nil {
			fmt.Println("Failed to save journal entry:", err)
		}
	case "journal:list":
		days := 0
		markdown := false

		for i, arg := range os.Args {
			if arg == "--last" && i+1 < len(os.Args) {
				days, _ = strconv.Atoi(os.Args[i+1])
			}
			if arg == "--md" {
				markdown = true
			}
		}

		err := journal.ListEntries(days, markdown)
		if err != nil {
			fmt.Println("Error loading journal:", err)
		}

	case "report":
		handleReport()
	default:
		fmt.Println("Unknown command")
	}
}

func handleIdeaAdd(args []string) {
	ideaText := args[0]
	useAI := len(args) > 1 && args[1] == "--ai"

	feedback := ""
	if useAI {
		feedback = getAIFeedback(ideaText)
	}

	err := ideas.SaveIdea(ideaText, feedback)
	if err != nil {
		fmt.Println("Failed to save idea:", err)
		return
	}

	fmt.Println("Captured Idea:", ideaText)
	if feedback != "" {
		fmt.Println("AI Feedback:", feedback)
	}
}

func handleIdeaList() {
	err := ideas.ListIdeas()
	if err != nil {
		fmt.Println("Error listing ideas:", err)
	}
}

func handleTimeStart(ideaID int, note string) {
	err := timebox.StartSession(ideaID, note)
	if err != nil {
		fmt.Println("Error starting time session:", err)
	}
}

func handleReport() {
	err := timebox.Report()
	if err != nil {
		fmt.Println("Error generating report:", err)
	}
}

func getAIFeedback(idea string) string {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	ctx := context.Background()

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "You're Chuck, a no-BS AI dev mentor. Be smart, blunt, and funny. Help improve ideas.",
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("Here’s the idea: %s\nWhat do you think?", idea),
			},
		},
	})

	if err != nil {
		return fmt.Sprintf("AI Error: %v", err)
	}

	return resp.Choices[0].Message.Content
}
