package main

import (
	"context"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rivo/tview"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func (a *App) SetFocus(primitive tview.Primitive) {
	app.tviewApp.SetFocus(primitive)
}

func (a *App) findItem(text string, items []ListItem) ListItem {
	for _, item := range items {
		if item.Name() == text {
			return item
		}
	}
	return nil
}

func (a *App) handleProjectSelect(index int, mainText string, secondaryText string, shortcut rune) {
	a.projectsPage.currentItem = a.findItem(mainText, a.projectsPage.listItems)
	issues := listProjectIssues(a.projectsPage.currentItem)
	var text string = ""
	for _, issue := range issues {
		text += "# " + issue.Title + "\n"
	}
	a.projectsPage.textView.SetText(text)
}

func (a *App) createPrimitive(title string) *tview.TextView {
	textView := tview.NewTextView().SetDynamicColors(true).SetRegions(true).SetChangedFunc(func() { a.tviewApp.Draw() })
	textView.SetTitle(title)
	textView.SetBorder(true)
	return textView
}

func (a *App) switchToPageFunc(pageName string) *tview.Pages {
	return a.pages.SwitchToPage(pageName)
}

func addToArray(items []ListItem, m map[string]ListItem, values []fuzzy.Rank) []ListItem {
	for _, value := range values {
		item := m[value.Target]
		if !contains(items, item) {
			items = append(items, item)
		}
	}
	return items
}

func contains(s []ListItem, value ListItem) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}
	return false
}

func (a *App) searchProjectsViewList() []int {
	searchString := a.projectSearchField.GetText()
	var ignoreCase bool = strings.ToLower(searchString) == searchString
	indices := a.projectsViewList.FindItems(searchString, searchString, false, ignoreCase)
	return indices
}

func (a *App) createProjectListView() {
	a.projectsPage.listView = tview.NewList().
		ShowSecondaryText(false)
	a.projectsPage.listView.SetTitle("Issues")

	a.projectsPage.populateProjectsViewList(context.Background(), a.switchToPageFunc)
	if len(a.projectsPage.listItems) > 0 {
		// TODO: change this to use pages
		a.handleProjectSelect(0, "", "", 'a')
	}

	a.projectsPage.listView = setNavigation(a.projectsPage.listView, func(event *tcell.EventKey) {
		if event.Key() == tcell.KeyTab && a.projectsPage.listView.HasFocus() {
			a.tviewApp.SetFocus(a.projectsPage.searchField)
		}
	})
	a.projectsPage.listView.SetChangedFunc(a.handleProjectSelect)
}

func (a *App) createProjectsPage() {
	a.projectsPage.CreateSearchField(a.SetFocus, a.switchToPageFunc)
	a.createProjectListView()
	a.projectsPage.CreatePageGrid()
	a.pages.AddPage("projects", a.projectsPage.gridView, true, true)
}

func setNavigation(list *tview.List, extraHandler func(event *tcell.EventKey)) *tview.List {
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
		if extraHandler != nil {
			extraHandler(event)
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

	list = setNavigation(list, nil)

	return list
}

func createIssueView(issues []*gitlab.Issue) (*tview.Flex, *tview.TextView) {
	issueView := showAllIssues(issues)
	textView := app.createPrimitive("")
	if len(issues) > 0 {
		issueText := getIssueDetails(issues[0])
		textView.SetText(issueText)
	}
	return IssueGrid(issueView, textView), textView
}

func handleIssueSelect(index int, mainText string, secondaryText string, shortcut rune) {
	issues, ok := app.projectIssues[app.currentProject]
	if !ok {
		return
	}

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
