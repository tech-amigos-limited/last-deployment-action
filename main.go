package main

import (
	"last-deployment-action/action"

	"github.com/sethvargo/go-githubactions"
)

func main() {
	token := githubactions.GetInput("github-token")
	repo := githubactions.GetInput("repo")
	ref := githubactions.GetInput("ref")

	githubactions.Debugf("received github-token %s", token)
	githubactions.Debugf("received repo %s", repo)
	githubactions.Debugf("received ref %s", ref)

	id := action.ActionImpl(&token, &repo, &ref)

	githubactions.SetOutput("deployment_id", id)
}
