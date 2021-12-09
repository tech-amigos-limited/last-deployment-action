package main

import (
	"get-deployment-action/action"
	"strconv"

	"github.com/sethvargo/go-githubactions"
)

func main() {
	token := githubactions.GetInput("github-token")
	repo := githubactions.GetInput("repo")
	ref := githubactions.GetInput("ref")

	githubactions.Debugf("received github-token %s", token)
	githubactions.Debugf("received repo %s", repo)
	githubactions.Debugf("received ref %s", ref)

	id, status := action.ActionImpl(&token, &repo, &ref)

	if id == 0 {
		githubactions.SetOutput("deployment_id", "")
	} else {
		githubactions.SetOutput("deployment_id", strconv.FormatInt(id, 10))
	}

	githubactions.SetOutput("status", status)
}
