package main

import (
	"github.com/gdamore/tcell/v2"
)

var app App

// TODO: rewrite have a type/struct for the entire page and then generalize the search function and the grid function and all those functions to work with the entire page

func main() {
	app = *createApp()
	app.showProjects()
	go app.getProjectsAndIssuesRoutine()

	app.createProjectsView()

	app.tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			app.tviewApp.Stop()
		}
		return event
	})

	if err := app.tviewApp.SetRoot(app.pages, true).SetFocus(app.pages).Run(); err != nil {
		panic(err)
		app.tviewApp.Stop()
	}
}
