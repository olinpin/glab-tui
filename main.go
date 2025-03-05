package main

var app App

// TODO: create a config file where you can put the env variables stuff, instead of getting the environment variable directly
//    - that way user can set the environment variables to whatever they want

func main() {
	app = *createApp()
	app.createProjectsPage()
	go app.getProjectsAndIssuesRoutine()

	if err := app.tviewApp.SetRoot(app.pages, true).SetFocus(app.pages).Run(); err != nil {
		panic(err)
		app.tviewApp.Stop()
	}
}
