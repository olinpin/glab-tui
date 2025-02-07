package main

import (
	"context"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rivo/tview"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func projectsGrid(menu *tview.List) *tview.Flex {
	menu.SetChangedFunc(handleProjectSelect)

	menu.SetTitle("Projects").SetBorder(true)
	app.projectsTextView.SetTitle("Issues").SetBorder(true)

	searchField := app.CreateSearchField()
	grid := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(searchField, 1, 0, true).
		AddItem(menu, 0, 1, false)

	mainUI := tview.NewFlex().
		AddItem(grid, 45, 1, true).
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

func (a *App) populateProjectsViewList(ctx context.Context) {
	var projects []*gitlab.Project
	if a.projectSearchField.GetText() == "" {
		projects = a.projects
	} else {
		projects = a.searchProjects(ctx)
	}
	select {
	case <-ctx.Done():
		return
	default:
		a.projectsViewList.Clear()
		for _, project := range projects {
			a.projectsViewList.AddItem(project.Name, string(project.ID), rune(0), func() {
				a.pages.SwitchToPage("issues" + project.Name)
			})
		}
	}
}

func (a *App) searchProjects(ctx context.Context) []*gitlab.Project {
	if ctx.Err() != nil {
		return nil
	}

	searchString := a.projectSearchField.GetText()
	projectNamesMap := make(map[string]*gitlab.Project, len(a.projects))
	projectNames := make([]string, 0, len(a.projects))

	for _, project := range a.projects {
		select {
		case <-ctx.Done():
			return nil
		default:
			projectNames = append(projectNames, project.Name)
			projectNamesMap[project.Name] = project
		}
	}

	names := fuzzy.RankFindFold(searchString, projectNames)
	sort.Sort(names)
	projects := addToArray([]*gitlab.Project{}, projectNamesMap, names)
	return projects
}

func addToArray(projects []*gitlab.Project, m map[string]*gitlab.Project, values []fuzzy.Rank) []*gitlab.Project {
	for _, value := range values {
		project := m[value.Target]
		if !contains(projects, project) {
			projects = append(projects, project)
		}
	}
	return projects
}

func contains[T comparable](s []T, value T) bool {
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

func (a *App) showProjects() {
	a.projectsViewList = tview.NewList().
		ShowSecondaryText(false)
	a.populateProjectsViewList(context.Background())
	if len(a.projects) > 0 {
		handleProjectSelect(0, "", "", 'a')
	}

	a.projectsViewList = setNavigation(a.projectsViewList)
}

func setNavigation(list *tview.List) *tview.List {
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

	list = setNavigation(list)

	return list
}

func InputFieldChangedFunc(text string) {
	if app.searchCancel != nil {
		app.searchCancel()
		app.searchCancel = nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	app.searchCancel = cancel

	go func() {
		app.tviewApp.QueueUpdateDraw(func() {
			app.populateProjectsViewList(ctx)
		})
	}()
}

func (a *App) CreateSearchField() *tview.InputField {
	a.projectSearchField = tview.NewInputField().
		SetLabel("Search: ")
		// a.projectSearchField.SetInputCapture(SearchInputCapture)
	a.projectSearchField.SetChangedFunc(InputFieldChangedFunc)
	return a.projectSearchField
}

func (a *App) createProjectsView() {
	a.projectsView = projectsGrid(a.projectsViewList)
	a.projectsView.SetInputCapture(projectsViewInputCapture)
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
