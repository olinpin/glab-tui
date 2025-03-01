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
	issueViews map[*gitlab.Project]*tview.TextView
	mu         sync.Mutex
}

type SafeCache struct {
	cache map[string]TimedCached
	mu    sync.Mutex
}

type App struct {
	tviewApp           *tview.Application
	projectsView       *tview.Flex
	projectsViewList   *tview.List
	projects           []*gitlab.Project
	projectIssues      map[*gitlab.Project][]*gitlab.Issue
	git                *gitlab.Client
	projectsTextView   *tview.TextView
	safeCache          SafeCache
	pages              *tview.Pages
	currentProject     *gitlab.Project
	safeIssueViews     SafeIssueViews
	projectSearchField *tview.InputField
	searchCancel       context.CancelFunc
	projectsPage       Page
}

type ListItem interface {
	ID() int
	Name() string
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

func (i IssueWrapper) ID() int {
	return i.issue.ID
}

func (i IssueWrapper) Name() string {
	return i.issue.Title
}

type TimedCached struct {
	timestamp int64
	value     any
}

func createApp() *App {
	a := App{}
	a.git = getGitlab(os.Getenv("GITLAB_TOKEN"), "https://gitlab.utwente.nl")
	a.safeCache = SafeCache{cache: map[string]TimedCached{}}
	a.projectIssues = map[*gitlab.Project][]*gitlab.Issue{}
	a.safeIssueViews = SafeIssueViews{issueViews: map[*gitlab.Project]*tview.TextView{}}
	a.tviewApp = tview.NewApplication()
	a.pages = tview.NewPages()
	a.projects = []*gitlab.Project{}
	a.projectsTextView = a.createPrimitive("")
	a.projectSearchField = tview.NewInputField()
	a.projectsPage = Page{
		textView: a.createPrimitive("Issues"),
	}
	return &a
}

func (a *App) getProjectsAndIssuesRoutine() {
	a.listProjects()
	a.projectsPage.populateProjectsViewList(context.Background(), a.switchToPageFunc)
	for _, project := range a.projectsPage.listItems {
		a.createIssuePage(project)
	}
}

func issueViewInputCapture(event *tcell.EventKey) *tcell.EventKey {
	k := event.Key()
	if k == tcell.KeyEsc {
		app.pages.SwitchToPage("projects")
	}
	return event
}

func viewInputCapture(event *tcell.EventKey) *tcell.EventKey {
	k := event.Key()
	if k == tcell.KeyEsc {
		app.pages.SwitchToPage("projects")
	}
	return event
}

func (a *App) createIssuePage(project ListItem) {
	// TODO: fix once you start fixing issue pages
	// issues := listProjectIssues(project)
	// a.projectIssues[project] = issues
	// issueView, textView := createIssueView(issues)
	// issueView.SetInputCapture(issueViewInputCapture)
	// a.safeIssueViews.mu.Lock()
	// a.safeIssueViews.issueViews[project] = textView
	// a.safeIssueViews.mu.Unlock()
	// a.pages.AddPage("issues"+project.Name, issueView, true, false)
}
