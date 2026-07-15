package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/b2netpro/builderos/internal/ideas"
)

func main() {
	a := app.New()
	w := a.NewWindow("BuilderOS GUI")
	w.Resize(fyne.NewSize(600, 400))

	// List of ideas
	ideaList := widget.NewList(
		func() int {
			data, _ := ideas.LoadIdeas()
			return len(data)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Idea")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			data, _ := ideas.LoadIdeas()
			o.(*widget.Label).SetText(fmt.Sprintf("[%d] %s", data[i].ID, data[i].Text))
		},
	)

	// Entry for new idea
	ideaEntry := widget.NewEntry()
	ideaEntry.SetPlaceHolder("Enter your idea")

	// Save button
	saveBtn := widget.NewButton("Add Idea", func() {
		text := ideaEntry.Text
		if text == "" {
			dialog.ShowError(fmt.Errorf("Idea can't be empty"), w)
			return
		}
		ideas.SaveIdea(text, "")
		ideaEntry.SetText("")
		ideaList.Refresh()
	})

	// Layout
	inputRow := container.NewHBox(ideaEntry, layout.NewSpacer(), saveBtn)
	ui := container.NewBorder(inputRow, nil, nil, nil, ideaList)
	w.SetContent(ui)

	w.ShowAndRun()
}
