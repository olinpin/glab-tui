package main

var app App

// TODO: rewrite have a type/struct for the entire page and then generalize the search function and the grid function and all those functions to work with the entire page

func main() {
	app = *createApp()
	app.createProjectsPage()
	go app.getProjectsAndIssuesRoutine()

	if err := app.tviewApp.SetRoot(app.pages, true).SetFocus(app.pages).Run(); err != nil {
		panic(err)
		app.tviewApp.Stop()
	}
}
