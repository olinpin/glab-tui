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
var projects []*gitlab.Project
var projectIssues map[*gitlab.Project][]*gitlab.Issue
var git *gitlab.Client
var projectsTextView *tview.TextView
var cache map[string]TimedCached
var pages *tview.Pages
var currentProject *gitlab.Project
var issueViews map[*gitlab.Project]*tview.TextView

// TODO: loading is taking way too long, first open the app and then populate the projects list and then download the issues on select

func main() {
	git = getGitlab(os.Getenv("GITLAB_TOKEN"), "https://gitlab.utwente.nl")
	projects = listProjects()
	cache = map[string]TimedCached{}
	projectIssues = map[*gitlab.Project][]*gitlab.Issue{}
	issueViews = map[*gitlab.Project]*tview.TextView{}

	// var project string = "s2969912/glabtest"
	app = tview.NewApplication()
	pages = tview.NewPages()
	projectsView = createProjectsView(projects)
	pages.AddPage("projects", projectsView, true, true)
	for _, project := range projects {
		issues := listProjectIssues(project)
		projectIssues[project] = issues
		issueView, textView := createIssueView(issues)
		issueViews[project] = textView
		pages.AddPage("issues"+project.Name, issueView, true, false)
	}

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
