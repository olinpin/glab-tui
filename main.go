package main

import (
	"fmt"
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
var git *gitlab.Client
var textView *tview.TextView
var cache map[string]TimedCached
var pages *tview.Pages

func main() {
	git = getGitlab(os.Getenv("GITLAB_TOKEN"), "https://gitlab.utwente.nl")
	projects = listProjects()
	cache = map[string]TimedCached{}

	// var project string = "s2969912/glabtest"
	app = tview.NewApplication()
	pages = tview.NewPages()
	projectsView = createProjectsView(projects)
	pages.AddPage("projects", projectsView, true, true)
	for _, project := range projects {
		issueView := createIssueView(listProjectIssues(project))
		pages.AddPage("issues"+project.Name, issueView, true, false)
	}

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()
		if k == tcell.KeyEsc {
			pages.SwitchToPage("projects")
		}
		return event
	})

	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		fmt.Println(err)
		app.Stop()
	}
}
