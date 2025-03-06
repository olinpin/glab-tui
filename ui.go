package main

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rivo/tview"
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
		text += "# " + issue.Name() + "\n"
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

func (p *Page) CreateListView(a *App, listTitle string, handleSelect func(index int, mainText string, secondaryText string, shortcut rune)) {
	p.listView = tview.NewList().
		ShowSecondaryText(false)
	p.listView.SetTitle(listTitle)
	p.listView.SetSelectedFocusOnly(true)

	p.PopulateListView(context.Background(), a.switchToPageFunc)
	if len(p.listItems) > 0 {
		handleSelect(0, "", "", 'a')
	}

	p.listView = setNavigation(p.listView, func(event *tcell.EventKey) {
		if event.Key() == tcell.KeyTab && p.listView.HasFocus() {
			a.tviewApp.SetFocus(p.searchField)
		}
	})
	p.listView.SetChangedFunc(handleSelect)
}

func (a *App) createProjectsPage() {
	a.projectsPage.CreateSearchField(a.SetFocus, a.switchToPageFunc)
	a.projectsPage.CreateListView(a, "Issues", a.handleProjectSelect)
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

func handleIssueSelect(index int, mainText string, secondaryText string, shortcut rune) {
	currentProject := app.projectsPage.currentItem
	issues, ok := app.projectIssues[currentProject]
	if !ok || len(issues) == 0 {
		return
	}

	issueText := getIssueDetails(issues[index])
	app.safeIssueViews.mu.Lock()
	page, _ := app.safeIssueViews.issueViews[currentProject]
	textView := page.textView
	textView.SetText(issueText)
	app.safeIssueViews.mu.Unlock()
}
