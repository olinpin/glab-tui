package main

import (
	"os"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type Issue struct {
	id   int
	text string
}

type TimedCached struct {
	timestamp int64
	value     any
}

type App struct {
	tviewApp         *tview.Application
	projectsView     *tview.Flex
	projectsViewList *tview.List
	projects         []*gitlab.Project
	projectIssues    map[*gitlab.Project][]*gitlab.Issue
	git              *gitlab.Client
	projectsTextView *tview.TextView
	cache            map[string]TimedCached
	pages            *tview.Pages
	currentProject   *gitlab.Project
	safeIssueViews       SafeIssueViews
}

type SafeIssueViews struct {
	issueViews map[*gitlab.Project]*tview.TextView
	mu         sync.Mutex
}

var app App

func createApp() *App {
	a := App{}
	a.git = getGitlab(os.Getenv("GITLAB_TOKEN"), "https://gitlab.utwente.nl")
	a.cache = map[string]TimedCached{}
	a.projectIssues = map[*gitlab.Project][]*gitlab.Issue{}
    a.safeIssueViews = SafeIssueViews{issueViews: map[*gitlab.Project]*tview.TextView{}}
	a.tviewApp = tview.NewApplication()
	a.pages = tview.NewPages()
	a.projects = []*gitlab.Project{}
	a.projectsTextView = createPrimitive("")
	return &a
}

func (a *App) getProjectsAndIssuesRoutine() {
	a.listProjects()
	a.populateProjectsViewList()
	for _, project := range a.projects {
		go a.createIssuePage(project)
	}
}

func (a *App) createIssuePage(project *gitlab.Project) {
	a.projectIssues[project] = listProjectIssues(project)
	issueView, textView := createIssueView(a.projectIssues[project])
    a.safeIssueViews.mu.Lock()
	a.safeIssueViews.issueViews[project] = textView
    a.safeIssueViews.mu.Unlock()
	a.pages.AddPage("issues"+project.Name, issueView, true, false)
}

func main() {
	app = *createApp()
	app.showProjects()
	go app.getProjectsAndIssuesRoutine()

	app.createProjectsView()

	app.tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()
		if k == tcell.KeyEsc {
			app.pages.SwitchToPage("projects")
		} else if event.Rune() == 'q' {
			app.tviewApp.Stop()
		}
		return event
	})

	if err := app.tviewApp.SetRoot(app.pages, true).SetFocus(app.pages).Run(); err != nil {
		panic(err)
		app.tviewApp.Stop()
	}
}
