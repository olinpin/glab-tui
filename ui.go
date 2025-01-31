package main

import (
	"github.com/rivo/tview"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func grid(menu *tview.List, main *tview.TextView) *tview.Flex {

	menu.SetChangedFunc(handleListSelect)

	menu.SetTitle("Projects").SetBorder(true)

	mainUI := tview.NewFlex().
		// SetRows(3, 0, 3).
		// SetColumns(30, 0, 30).
		// SetBorders(true).
		AddItem(menu, 45, 1, true).
		AddItem(main, 0, 3, false)

	return mainUI
}

func handleListSelect(index int, mainText string, secondaryText string, shortcut rune) {
	issues := listProjectIssues(projects[index])
    var text string = ""
    for _, issue := range issues {
        text += "# " + issue.Title + "\n"
    }
    textView.SetText(text)
}

func createPrimitive(text string) *tview.TextView {
	textView := tview.NewTextView().SetDynamicColors(true).SetRegions(true).SetChangedFunc(func() { app.Draw() })
	textView.SetText(text)
	textView.SetBorder(true)
	return textView
}

func help() *tview.List {
	list := tview.NewList().
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		}).
		AddItem("List item 2", "Explain", 'j', func() {
		}).
		AddItem("List item 3", "Explain", 'k', func() {
		})
	return list
}

func showProjects(projects []*gitlab.Project) *tview.List {
	list := tview.NewList().
		ShowSecondaryText(false)

	for _, project := range projects {
		list.AddItem(project.Name, string(project.ID), 'a', nil)
	}
	// AddItem("Quit", "Press to exit", 'q', func() {
	// 	app.Stop()
	// }).
	// AddItem("List item 3", "Explain", 'k', func() {
	// })
	return list
}
