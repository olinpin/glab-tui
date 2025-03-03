package main

import (
	"context"
	"os"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type SafeIssueViews struct {
	issueViews map[ListItem]Page
	mu         sync.Mutex
}

type SafeCache struct {
	cache map[string]TimedCached
	mu    sync.Mutex
}

type App struct {
	tviewApp           *tview.Application
	projectIssues      map[ListItem][]ListItem
	git                *gitlab.Client
	safeCache          SafeCache
	pages              *tview.Pages
	safeIssueViews     SafeIssueViews
	projectsPage       Page
}

type ListItem interface {
	ID() int
	Name() string
	Description() string
}

type ProjectWrapper struct {
	project *gitlab.Project
}

type IssueWrapper struct {
	issue *gitlab.Issue
}

func (p ProjectWrapper) ID() int {
	return p.project.ID
}

func (p ProjectWrapper) Name() string {
	return p.project.Name
}

func (p ProjectWrapper) Description() string {
	return p.project.Description
}

func (i IssueWrapper) ID() int {
	return i.issue.ID
}

func (i IssueWrapper) Name() string {
	return i.issue.Title
}

func (i IssueWrapper) Description() string {
	return i.issue.Description
}

type TimedCached struct {
	timestamp int64
	value     any
}

func createApp() *App {
	a := App{}
	a.git = getGitlab(os.Getenv("GITLAB_TOKEN"), "https://gitlab.utwente.nl")
	a.safeCache = SafeCache{cache: map[string]TimedCached{}}
	a.projectIssues = map[ListItem][]ListItem{}
	a.safeIssueViews = SafeIssueViews{issueViews: map[ListItem]Page{}}
	a.tviewApp = tview.NewApplication()
	a.pages = tview.NewPages()
	a.projectsPage = Page{
		textView: a.createPrimitive("Issues"),
	}
	return &a
}

func (a *App) getProjectsAndIssuesRoutine() {
	a.listProjects()
	a.projectsPage.PopulateListView(context.Background(), a.switchToPageFunc)
	for _, project := range a.projectsPage.listItems {
		a.createIssuePage(project)
	}
}

func viewInputCapture(event *tcell.EventKey) *tcell.EventKey {
	k := event.Key()
	if k == tcell.KeyEsc {
		app.pages.SwitchToPage("projects")
	}
	return event
}

func (a *App) createIssuePage(project ListItem) {
	issues := listProjectIssues(project)
	a.projectIssues[project] = issues
	page := Page{
		textView: a.createPrimitive(""),
	}
	a.safeIssueViews.mu.Lock()
	a.safeIssueViews.issueViews[project] = page
	a.safeIssueViews.mu.Unlock()
	page.listItems = issues
	// TODO: fix where does the switchToPage lead
	page.CreateSearchField(a.SetFocus, a.switchToPageFunc)
	page.CreateListView(a, "Issues", handleIssueSelect)
	page.CreatePageGrid()
	if len(issues) > 0 {
		issueText := getIssueDetails(issues[0])
		page.textView.SetText(issueText)
	}
	a.pages.AddPage("issues"+project.Name(), page.gridView, true, false)
}
