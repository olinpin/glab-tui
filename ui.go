package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func projectsGrid(menu *tview.List) *tview.Flex {
	menu.SetChangedFunc(handleListSelect)

	menu.SetTitle("Projects").SetBorder(true)
	if projectsTextView == nil {
		projectsTextView = createPrimitive("")
	}
	projectsTextView.SetTitle("Issues").SetBorder(true)

	mainUI := tview.NewFlex().
		AddItem(menu, 45, 1, true).
		AddItem(projectsTextView, 0, 3, false)

	return mainUI
}

func handleListSelect(index int, mainText string, secondaryText string, shortcut rune) {
	issues := listProjectIssues(projects[index])
	var text string = ""
	for _, issue := range issues {
		text += "# " + issue.Title + "\n"
	}
	projectsTextView.SetText(text)
}

func createPrimitive(text string) *tview.TextView {
	textView := tview.NewTextView().SetDynamicColors(true).SetRegions(true).SetChangedFunc(func() { app.Draw() })
	textView.SetText(text)
	textView.SetBorder(true)
	return textView
}

func showProjects(projects []*gitlab.Project) *tview.List {
	list := tview.NewList().
		ShowSecondaryText(false)

	for _, project := range projects {
		list.AddItem(project.Name, string(project.ID), rune(0), func() {
			pages.SwitchToPage("issues" + project.Name)
		})
	}
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Rune()
		currentItem := list.GetCurrentItem()
		switch k {
		case 'j':
			if currentItem < list.GetItemCount() {
				list.SetCurrentItem(currentItem + 1)
			}
		case 'k':
			if currentItem > 0 {
				list.SetCurrentItem(currentItem + -1)
			}
		case 'g':
			list.SetCurrentItem(0)
		case 'G':
			list.SetCurrentItem(list.GetItemCount() - 1)
		}
		return event
	})
	return list
}

func showAllIssues(issues []*gitlab.Issue) *tview.List {
	list := tview.NewList().
		ShowSecondaryText(false)

	for _, issue := range issues {
		list.AddItem(issue.Title, string(issue.ID), 0, nil)
	}

	return list
}

func createProjectsView(projects []*gitlab.Project) *tview.Flex {
	projectsUI := showProjects(projects)
	return projectsGrid(projectsUI)
}

func createIssueView(issues []*gitlab.Issue) *tview.Flex {
	issueView := showAllIssues(issues)
	textView := createPrimitive("")
	return IssueGrid(issueView, textView)
}

func IssueGrid(menu *tview.List, main *tview.TextView) *tview.Flex {
	menu.SetChangedFunc(handleListSelect)

	menu.SetTitle("Issues").SetBorder(true)
	main.SetTitle("Issue #").SetBorder(true)

	mainUI := tview.NewFlex().
		AddItem(menu, 45, 1, true).
		AddItem(main, 0, 3, false)

	return mainUI
}
