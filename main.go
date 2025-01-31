package main

import (
	"fmt"
	"os"

	"github.com/rivo/tview"
	gitlab "gitlab.com/gitlab-org/api/client-go"
	// "gitlab.com/gitlab-org/api/client-go"
)

type Issue struct {
	id   int
	text string
}

var app *tview.Application
var mainUI *tview.Flex
var projects []*gitlab.Project
var git *gitlab.Client
var textView *tview.TextView

func main() {
	git = getGitlab(os.Getenv("GITLAB_TOKEN"), "https://gitlab.utwente.nl")
	projects = listProjects()

	// var project string = "s2969912/glabtest"
	app = tview.NewApplication()
	projectsUI := showProjects(projects)
    textView = createPrimitive("")
	mainUI = grid(projectsUI, textView)
	if err := app.SetRoot(mainUI, true).SetFocus(mainUI).Run(); err != nil {
		fmt.Println(err)
        app.Stop()
	}
}
