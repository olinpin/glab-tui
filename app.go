package main

import (
	"os"
	"sync"

	"github.com/rivo/tview"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type SafeIssueViews struct {
	issueViews map[*gitlab.Project]*tview.TextView
	mu         sync.Mutex
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
	safeIssueViews   SafeIssueViews
}


type TimedCached struct {
	timestamp int64
	value     any
}

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
