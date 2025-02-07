package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func projectsGrid(menu *tview.List) *tview.Flex {
	menu.SetChangedFunc(handleProjectSelect)

	menu.SetTitle("Projects").SetBorder(true)
	projectsTextView.SetTitle("Issues").SetBorder(true)

	mainUI := tview.NewFlex().
		AddItem(menu, 45, 1, true).
		AddItem(projectsTextView, 0, 3, false)

	return mainUI
}

func handleProjectSelect(index int, mainText string, secondaryText string, shortcut rune) {
	currentProject = projects[index]
	issues := listProjectIssues(currentProject)
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
	// TODO: This is repetative copy paste code, figure out how to generalize it
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

func createProjectsView(projects []*gitlab.Project) *tview.Flex {
	projectsUI := showProjects(projects)
	if projectsTextView == nil {
		projectsTextView = createPrimitive("")
	}
	handleProjectSelect(0, "", "", 'a')
	return projectsGrid(projectsUI)
}

func createIssueView(issues []*gitlab.Issue) (*tview.Flex, *tview.TextView) {
	issueView := showAllIssues(issues)
	textView := createPrimitive("")
	if len(issues) > 0 {
		issueText := getIssueDetails(issues[0])
		textView.SetText(issueText)
	}
	return IssueGrid(issueView, textView), textView
}

func handleIssueSelect(index int, mainText string, secondaryText string, shortcut rune) {
	issues := projectIssues[currentProject]
	issueText := getIssueDetails(issues[index])
	textView, _ := issueViews[currentProject]
	textView.SetText(issueText)
}

func IssueGrid(menu *tview.List, main *tview.TextView) *tview.Flex {
	menu.SetChangedFunc(handleIssueSelect)

	menu.SetTitle("Issues").SetBorder(true)
	main.SetTitle("Issue #").SetBorder(true)

	mainUI := tview.NewFlex().
		AddItem(menu, 45, 1, true).
		AddItem(main, 0, 3, false)

	return mainUI
}
