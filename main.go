package main

import (
	"fmt"
	"os"

	"github.com/rivo/tview"
	// "gitlab.com/gitlab-org/api/client-go"
)

type Issue struct {
	id   int
	text string
}

func main() {
	git := getGitlab(os.Getenv("GITLAB_TOKEN"), "https://gitlab.utwente.nl")
    projects := listProjects(git)
    fmt.Println(projects)
    // var project string = "s2969912/glabtest"
	app := tview.NewApplication()
	helpList := grid(app)
	if err := app.SetRoot(helpList, true).SetFocus(helpList).Run(); err != nil {
		panic(err)
	}
}

