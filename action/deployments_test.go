package action

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-github/v41/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetSingleDeploymentNoHistory(t *testing.T) {

	var deployment github.Deployment
	deployment.ID = Int64(123)
	deployment.Environment = String("feature")

	var status Status
	status.Id = Int64(234)
	status.State = String("pending")

	mockClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposDeploymentsByOwnerByRepo,
			[]github.Deployment{
				deployment,
			},
		),
		mock.WithRequestMatch(
			mock.GetReposDeploymentsStatusesByOwnerByRepoByDeploymentId,
			[]Status{
				status,
			},
		),
	)
	client := github.NewClient(mockClient)
	context := context.Background()

	history, err := GetDeploymentHistory(context, client, &Args{
		Owner: String("autorama"), Repo: String("nsf"), Ref: String("head"),
	})
	if err != nil {
		assert.Error(t, err)
	}

	assert.Len(t, history, 1)
	assert.Equal(t, history[0].DeploymentId, Int64(123))
	assert.Equal(t, *history[0].Environment, "feature")
	assert.Equal(t, *history[0].Statuses[0].State, "pending")
}

func TestLatestIsPending(t *testing.T) {
	var deployment DeploymentHistory
	deployment.DeploymentId = Int64(123)
	deployment.Environment = String("feature")

	var status Status
	status.Id = Int64(234)
	status.State = String("pending")

	deployment.Statuses = append(deployment.Statuses, &status)

	history := []*DeploymentHistory{&deployment}

	id, err := GetLatestActiveDeploymentId(history)

	assert.Nil(t, id)
	assert.NotNil(t, err)
}

func TestLatestActive(t *testing.T) {
	var deployment DeploymentHistory
	deployment.DeploymentId = Int64(123)
	deployment.Environment = String("feature")

	var status Status
	status.Id = Int64(234)
	status.State = String("active")

	deployment.Statuses = append(deployment.Statuses, &status)

	history := []*DeploymentHistory{&deployment}

	id, err := GetLatestActiveDeploymentId(history)

	assert.Nil(t, err)
	assert.NotNil(t, id)
	assert.Equal(t, *Int64(123), *id, "deployment id is 123")
}

func TestNoDeployments(t *testing.T) {
	history := []*DeploymentHistory{}
	id, err := GetLatestActiveDeploymentId(history)

	assert.Nil(t, id)
	assert.NotNil(t, err)
}
func TestNoStatuses(t *testing.T) {
	var deployment DeploymentHistory
	deployment.DeploymentId = Int64(123)
	deployment.Environment = String("feature")

	history := []*DeploymentHistory{&deployment}

	id, err := GetLatestActiveDeploymentId(history)

	assert.Nil(t, id)
	assert.NotNil(t, err)
}

func TestSortSortedDeployments(t *testing.T) {
	before := github.Timestamp{}
	before.Time = time.Now()

	after := github.Timestamp{}
	after.Time = before.Time.Add(time.Duration(*Int64(100)))

	var d1 DeploymentHistory
	d1.DeploymentId = Int64(123)
	d1.Environment = String("feature")
	d1.CreatedAt = &before

	var d2 DeploymentHistory
	d2.DeploymentId = Int64(234)
	d2.Environment = String("feature")
	d2.CreatedAt = &after

	// d2 is most recent, so this is already sorted correctly
	history := []*DeploymentHistory{&d2, &d1}

	sortDeployments(history)

	assert.Equal(t, d2, *history[0])
	assert.Equal(t, d1, *history[1])
	assert.Len(t, history, 2)
}

func TestSortUnSortedDeployments(t *testing.T) {
	before := github.Timestamp{}
	before.Time = time.Now()

	after := github.Timestamp{}
	after.Time = before.Time.Add(time.Duration(*Int64(100)))

	var d1 DeploymentHistory
	d1.DeploymentId = Int64(123)
	d1.Environment = String("feature")
	d1.CreatedAt = &before

	var d2 DeploymentHistory
	d2.DeploymentId = Int64(234)
	d2.Environment = String("feature")
	d2.CreatedAt = &after

	// d2 is most recent, so this is not sorted correctly
	history := []*DeploymentHistory{&d1, &d2}

	sortDeployments(history)

	assert.Equal(t, d2, *history[0])
	assert.Equal(t, d1, *history[1])
	assert.Len(t, history, 2)
}

func TestSortUnSortedStatuses(t *testing.T) {
	before := github.Timestamp{}
	before.Time = time.Now()

	after := github.Timestamp{}
	after.Time = before.Time.Add(time.Duration(*Int64(100)))

	var d1 Status
	d1.Id = Int64(123)
	d1.CreatedAt = &before

	var d2 Status
	d2.Id = Int64(234)
	d2.CreatedAt = &after

	// d2 is most recent, so this is not sorted correctly
	statuses := []*Status{&d1, &d2}

	sortStatuses(statuses)

	assert.Equal(t, d2, *statuses[0])
	assert.Equal(t, d1, *statuses[1])
	assert.Len(t, statuses, 2)
}

func TestSortSortedStatuses(t *testing.T) {
	before := github.Timestamp{}
	before.Time = time.Now()

	after := github.Timestamp{}
	after.Time = before.Time.Add(time.Duration(*Int64(100)))

	var d1 Status
	d1.Id = Int64(123)
	d1.CreatedAt = &before

	var d2 Status
	d2.Id = Int64(234)
	d2.CreatedAt = &after

	// already sorted
	statuses := []*Status{&d2, &d1}

	sortStatuses(statuses)

	assert.Equal(t, d2, *statuses[0])
	assert.Equal(t, d1, *statuses[1])
	assert.Len(t, statuses, 2)
}
