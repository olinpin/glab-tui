package main

import (
	"sort"
	"strings"
	"time"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func getGitlab(token string, url string) *gitlab.Client {
	git, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		handleError(err)
	}
	return git
}

func (a *App) listProjects() {
	opt := &gitlab.ListProjectsOptions{}
	projects, _, err := a.git.Projects.ListProjects(opt)
	if err != nil {
		handleError(err)
	}

	sort.Slice(projects, func(i, j int) bool {
		return strings.ToLower(projects[i].Name) < strings.ToLower(projects[j].Name)
	})
	a.projectsPage.listItems = []ListItem{}
	for _, project := range projects {
		a.projectsPage.listItems = append(a.projectsPage.listItems, ProjectWrapper{project})
	}
}

func listProjectIssues(project ListItem) []ListItem {
	var key string = "project" + string(project.ID())
	var timestamp int64 = time.Now().Unix()

	app.safeCache.mu.Lock()
	cacheHit, i := app.safeCache.cache[key]
	app.safeCache.mu.Unlock()
	if i && cacheHit.timestamp > timestamp+60 {
		return cacheHit.value.([]ListItem)
	}
	opt := &gitlab.ListProjectIssuesOptions{}
	issues, _, err := app.git.Issues.ListProjectIssues(project.ID(), opt)
	if err != nil {
		handleError(err)
	}

	issueItems := []ListItem{}
	for _, issue := range issues {
		issueItems = append(issueItems, IssueWrapper{issue})
	}

	app.safeCache.mu.Lock()
	app.safeCache.cache[key] = TimedCached{timestamp, issueItems}
	app.safeCache.mu.Unlock()

	return issueItems
}

func getIssueDetails(issue ListItem) string {
	result := ""
	title := issue.Name()
	result += "# " + title + "\n\n"
	description := issue.Description()
	result += description + "\n\n"
	return result
}
