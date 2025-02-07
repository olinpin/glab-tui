package main

import (
	"os"

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

var app *tview.Application
var projectsView *tview.Flex
var projectsViewList *tview.List
var projects []*gitlab.Project
var projectIssues map[*gitlab.Project][]*gitlab.Issue
var git *gitlab.Client
var projectsTextView *tview.TextView
var cache map[string]TimedCached
var pages *tview.Pages
var currentProject *gitlab.Project
var issueViews map[*gitlab.Project]*tview.TextView

func getProjectsAndIssuesRoutine() {
	projects = listProjects()
	populateProjectsViewList(projects)
	for _, project := range projects {
		issues := listProjectIssues(project)
		projectIssues[project] = issues
		issueView, textView := createIssueView(issues)
		issueViews[project] = textView
		pages.AddPage("issues"+project.Name, issueView, true, false)
	}
}
func main() {
	git = getGitlab(os.Getenv("GITLAB_TOKEN"), "https://gitlab.utwente.nl")
	cache = map[string]TimedCached{}
	projectIssues = map[*gitlab.Project][]*gitlab.Issue{}
	issueViews = map[*gitlab.Project]*tview.TextView{}

	// create app and pages
	app = tview.NewApplication()
	pages = tview.NewPages()

	go getProjectsAndIssuesRoutine()

	projectsView = createProjectsView([]*gitlab.Project{})
	pages.AddPage("projects", projectsView, true, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()
		if k == tcell.KeyEsc {
			pages.SwitchToPage("projects")
		} else if event.Rune() == 'q' {
			app.Stop()
		}
		return event
	})

	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
		app.Stop()
	}
}
