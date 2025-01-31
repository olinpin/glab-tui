package main

import gitlab "gitlab.com/gitlab-org/api/client-go"

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
	return projects
}

func listProjectIssues(project *gitlab.Project) []*gitlab.Issue {
	opt := &gitlab.ListProjectIssuesOptions{}
	issues, _, err := git.Issues.ListProjectIssues(project.ID, opt)
	if err != nil {
		handleError(err)
	}
	return issues
}
