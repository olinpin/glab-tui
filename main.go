package main

import (
	"github.com/gdamore/tcell/v2"
)

var app App

func main() {
	app = *createApp()
	app.showProjects()
	go app.getProjectsAndIssuesRoutine()

	app.createProjectsView()

	app.tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()
		if k == tcell.KeyEsc {
			app.pages.SwitchToPage("projects")
		} else if event.Rune() == 'q' {
			app.tviewApp.Stop()
		}
		return event
	})

	if err := app.tviewApp.SetRoot(app.pages, true).SetFocus(app.pages).Run(); err != nil {
		panic(err)
		app.tviewApp.Stop()
	}
}
