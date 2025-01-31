package main

import gitlab "gitlab.com/gitlab-org/api/client-go"

func getGitlab(token string, url string) *gitlab.Client {
	git, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		handleError(err)
	}
    return git
}

func listProjects(git *gitlab.Client) []*gitlab.Project{
    opt:= &gitlab.ListProjectsOptions{Search: gitlab.Ptr("glab")}
    projects, _, err := git.Projects.ListProjects(opt)
	if err != nil {
		handleError(err)
	}
    return projects
}
