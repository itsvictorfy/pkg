package github

import (
	"testing"

	"github.com/google/go-github/v66/github"
	"github.com/stretchr/testify/assert"
)

// Test function for CreateGithubClient and TriggerWorkflowActionByID
func TestTriggerWorkflowActionByID_Integration(t *testing.T) {
	gh := &Github{
		OrgName:    "",
		Repo:       "",
		ActionId:   1, // Replace with a valid workflow ID
		ActionFile: "CD-Deploy-Env.yaml",
	}
	gh.CreateGithubClient()

	// Define request
	req := github.CreateWorkflowDispatchEventRequest{
		Ref: "main",
		Inputs: map[string]interface{}{
			"environment": "staging",
			"image_tag":   "1",
		}}

	// Attempt to trigger workflow
	err := gh.TriggerWorkflowActionByID(req)
	assert.NoError(t, err, "Expected no error when triggering workflow action")
}
func TestTriggerWorkflowActionByFileName_Integration(t *testing.T) {
	gh := &Github{
		OrgName:    "",
		Repo:       "",
		ActionId:   123456, // Replace with a valid workflow ID
		ActionFile: "CD-Deploy-Env.yaml",
		Token:      "",
	}
	gh.CreateGithubClient()

	// Define request
	req := github.CreateWorkflowDispatchEventRequest{
		Ref: "main",
		Inputs: map[string]interface{}{
			"environment": "staging",
			"image_tag":   "1",
		}}

	// Attempt to trigger workflow
	err := gh.TriggerWorkflowActionByFileName(req)
	assert.NoError(t, err, "Expected no error when triggering workflow action by filename")
}
func TestGithubAuth(t *testing.T) {
	gh := &Github{
		OrgName:    "",
		Repo:       "",
		ActionId:   123456, // Replace with a valid workflow ID
		ActionFile: "CD-Deploy-Env.yaml",
		Token:      "",
	}
	gh.CreateGithubClient()
	assert.NotNil(t, gh.Client, "Expected a non-nil Github client")
}

func TestGithubListWorkflows(t *testing.T) {
	gh := &Github{
		OrgName:    "50k-trade",
		Repo:       "server",
		ActionId:   123456, // Replace with a valid workflow ID
		ActionFile: "CD-Deploy-Env.yaml",
		Token:      "",
	}
	gh.CreateGithubClient()
	wf, err := gh.ListAllWorkflows()
	assert.NotNil(t, wf, "Expected a non-nil list of workflows")
	assert.NoError(t, err, "Expected no error when listing workflows")
}
