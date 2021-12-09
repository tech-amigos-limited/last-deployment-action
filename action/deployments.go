package action

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"github.com/google/go-github/v41/github"
	"github.com/sethvargo/go-githubactions"
	"golang.org/x/oauth2"
)

type DeploymentHistory struct {
	DeploymentId *int64
	Ref          *string
	Environment  *string
	CreatedAt    *github.Timestamp
	Statuses     []*Status
}

type Status struct {
	Id        *int64
	State     *string
	CreatedAt *github.Timestamp
}

type Args struct {
	Owner *string
	Repo  *string
	Ref   *string
}

func (d *DeploymentHistory) LastStatus() *Status {
	if len(d.Statuses) > 0 {
		return d.Statuses[0]
	}
	return nil
}

func ActionImpl(token *string, repo *string, ref *string) (id int64, status string) {
	context := context.Background()
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *token},
	)

	client := github.NewClient(oauth2.NewClient(context, tokenSource))

	s := strings.Split(*repo, "/")

	history, err := GetDeploymentHistory(context, client, &Args{
		Owner: String(strings.TrimSpace(s[0])),
		Repo:  String(strings.TrimSpace(s[1])),
		Ref:   String(*ref),
	})

	if err != nil {
		githubactions.Infof("Unable to get deployment history: %s\n", err.Error())
		return
	}

	// log some useful deployment history
	if len(history) > 0 {
		jsonHistory, err := json.MarshalIndent(history, "", " ")
		if err != nil {
			githubactions.Fatalf(err.Error())
		}
		var sb strings.Builder
		sb.WriteString(string(jsonHistory))
		sb.WriteString("\n")

		githubactions.Group("Ordered deployment history")
		githubactions.Infof(sb.String())
		githubactions.EndGroup()
	} else {
		githubactions.Infof("No deployment history found\n")
		return
	}

	id, status = GetLatestDeploymentInfo(history)

	if err != nil {
		githubactions.Infof("Unable to get latest active deployment id: %s\n", err.Error())
		return
	}

	return
}

func GetLatestDeploymentInfo(history []*DeploymentHistory) (id int64, status string) {
	if len(history) == 0 {
		githubactions.Infof("no deployments found in history")
		return
	}

	d := history[0]
	id = *d.DeploymentId

	if len(d.Statuses) == 0 {
		githubactions.Infof("no statuses found for the most recent deployment id [%d]", *d.DeploymentId)
		return
	}

	status = *d.Statuses[0].State
	return
}

// return an ordered array of deployments (most recent first), each with an ordered array
// of statuses (more recent first)
func GetDeploymentHistory(context context.Context, client *github.Client, args *Args) (history []*DeploymentHistory, err error) {

	deployments, _, err := client.Repositories.ListDeployments(context, *args.Owner, *args.Repo,
		&github.DeploymentsListOptions{Ref: *args.Ref})

	if err != nil {
		return
	}

	var deploymentHistories []*DeploymentHistory
	for _, d := range deployments {

		var ghStatuses []*github.DeploymentStatus
		ghStatuses, _, err = client.Repositories.ListDeploymentStatuses(context, *args.Owner, *args.Repo, *d.ID, nil)

		if err != nil {
			return
		}

		var deployment DeploymentHistory
		deployment.DeploymentId = d.ID
		deployment.Ref = d.Ref
		deployment.Environment = d.Environment
		deployment.CreatedAt = d.CreatedAt

		var statuses []*Status
		for _, s := range ghStatuses {
			var status Status
			status.Id = s.ID
			status.State = s.State
			status.CreatedAt = s.CreatedAt

			statuses = append(statuses, &status)
		}

		// Sort with latest status first within each deployment
		deployment.Statuses = sortStatuses(statuses)
		deploymentHistories = append(deploymentHistories, &deployment)
	}

	// Sort deployments, latest first
	return sortDeployments(deploymentHistories), nil
}

func sortDeployments(deployments []*DeploymentHistory) []*DeploymentHistory {
	sort.SliceStable(deployments, func(i, j int) bool {
		return deployments[i].CreatedAt.Time.After(deployments[j].CreatedAt.Time)
	})
	return deployments
}

func sortStatuses(statuses []*Status) []*Status {
	sort.SliceStable(statuses, func(i, j int) bool {
		return statuses[i].CreatedAt.Time.After(statuses[j].CreatedAt.Time)
	})
	return statuses
}

func String(v string) *string { return &v }
func Int64(v int64) *int64    { return &v }
