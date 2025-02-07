package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func projectsGrid(menu *tview.List) *tview.Flex {
	menu.SetChangedFunc(handleProjectSelect)

	menu.SetTitle("Projects").SetBorder(true)
	app.projectsTextView.SetTitle("Issues").SetBorder(true)

	mainUI := tview.NewFlex().
		AddItem(menu, 45, 1, true).
		AddItem(app.projectsTextView, 0, 3, false)

	return mainUI
}

func handleProjectSelect(index int, mainText string, secondaryText string, shortcut rune) {
	app.currentProject = app.projects[index]
	issues := listProjectIssues(app.currentProject)
	var text string = ""
	for _, issue := range issues {
		text += "# " + issue.Title + "\n"
	}
	app.projectsTextView.SetText(text)
}

func createPrimitive(text string) *tview.TextView {
	textView := tview.NewTextView().SetDynamicColors(true).SetRegions(true).SetChangedFunc(func() { app.tviewApp.Draw() })
	textView.SetText(text)
	textView.SetBorder(true)
	return textView
}

func (a *App) populateProjectsViewList() {
	for _, project := range a.projects {
		a.projectsViewList.AddItem(project.Name, string(project.ID), rune(0), func() {
			a.pages.SwitchToPage("issues" + project.Name)
		})
	}
}

func (a *App) showProjects() {
	a.projectsViewList = tview.NewList().
		ShowSecondaryText(false)
	a.populateProjectsViewList()
	if len(a.projects) > 0 {
		handleProjectSelect(0, "", "", 'a')
	}

	a.projectsViewList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Rune()
		currentItem := a.projectsViewList.GetCurrentItem()
		switch k {
		case 'j':
			if currentItem < a.projectsViewList.GetItemCount() {
				a.projectsViewList.SetCurrentItem(currentItem + 1)
			}
		case 'k':
			if currentItem > 0 {
				a.projectsViewList.SetCurrentItem(currentItem + -1)
			}
		case 'g':
			a.projectsViewList.SetCurrentItem(0)
		case 'G':
			a.projectsViewList.SetCurrentItem(a.projectsViewList.GetItemCount() - 1)
		}
		return event
	})
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

func (a *App) createProjectsView() {
	a.projectsView = projectsGrid(a.projectsViewList)
	a.pages.AddPage("projects", app.projectsView, true, true)
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
	issues := app.projectIssues[app.currentProject]
	issueText := getIssueDetails(issues[index])
	app.safeIssueViews.mu.Lock()
	textView, _ := app.safeIssueViews.issueViews[app.currentProject]
	app.safeIssueViews.mu.Unlock()
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
