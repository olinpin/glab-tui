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

func listProjects() []*gitlab.Project {
	opt := &gitlab.ListProjectsOptions{}
	projects, _, err := git.Projects.ListProjects(opt)
	if err != nil {
		handleError(err)
	}

	sort.Slice(projects, func(i, j int) bool {
		return strings.ToLower(projects[i].Name) < strings.ToLower(projects[j].Name)
	})
	return projects
}

func listProjectIssues(project *gitlab.Project) []*gitlab.Issue {

	var key string = "project" + string(project.ID)
	var timestamp int64 = time.Now().Unix()

	cacheHit, i := cache[key]
	if i && cacheHit.timestamp > timestamp+60 {
		return cacheHit.value.([]*gitlab.Issue)
	}
	opt := &gitlab.ListProjectIssuesOptions{}
	issues, _, err := git.Issues.ListProjectIssues(project.ID, opt)
	if err != nil {
		handleError(err)
	}
	cache[key] = TimedCached{timestamp, issues}
	return issues
}

func getIssueDetails(issue *gitlab.Issue) string {
	result := ""
	title := issue.Title
	result += "# " + title + "\n\n"
	description := issue.Description
	result += description + "\n\n"
	return result
}
