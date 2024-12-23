package github

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test function for CreateGithubClient and TriggerWorkflowActionByID
func TestTriggerWorkflowActionByID_Integration(t *testing.T) {
	gh := &Github{
		OrgName: os.Getenv("GITHUB_ORG_NAME"),
		Repo:    os.Getenv("GITHUB_REPO"),
		Token:   os.Getenv("GITHUB_TOKEN"),
	}
	gh.InitGithubClient()
	wf := Workflow{
		ID:  123456789,
		Ref: "master",
		Inputs: map[string]interface{}{
			"environment": "Staging",
			"image_tag":   "1",
		},
		ActionFile: "workflow.yml",
	}
	// Attempt to trigger workflow
	err := gh.TriggerWorkflowActionByID(wf)
	assert.NoError(t, err, "Expected no error when triggering workflow action")
}
func TestTriggerWorkflowActionByFileName_Integration(t *testing.T) {
	gh := &Github{
		OrgName: os.Getenv("GITHUB_ORG_NAME"),
		Repo:    os.Getenv("GITHUB_REPO"),
		Token:   os.Getenv("GITHUB_TOKEN"),
	}
	gh.InitGithubClient()

	// Define workflow
	wf := Workflow{
		ID:  123456789,
		Ref: "master",
		Inputs: map[string]interface{}{
			"environment": "Staging",
			"image_tag":   "1",
		},
		ActionFile: "workflow.yml",
	}

	// Attempt to trigger workflow
	err := gh.TriggerWorkflowActionByFileName(wf)
	assert.NoError(t, err, "Expected no error when triggering workflow action by filename")
}
func TestGithubAuth_Integration(t *testing.T) {
	gh := &Github{
		OrgName: os.Getenv("GITHUB_ORG_NAME"),
		Repo:    os.Getenv("GITHUB_REPO"),
		Token:   os.Getenv("GITHUB_TOKEN"),
	}
	gh.InitGithubClient()
	assert.NotNil(t, gh.Client, "Expected a non-nil Github client")
}

func TestGithubListWorkflows_Integration(t *testing.T) {
	gh := &Github{
		OrgName: os.Getenv("GITHUB_ORG_NAME"),
		Repo:    os.Getenv("GITHUB_REPO"),
		Token:   os.Getenv("GITHUB_TOKEN"),
	}
	gh.InitGithubClient()
	wf, err := gh.ListAllWorkflows()
	assert.NotNil(t, wf, "Expected a non-nil list of workflows")
	assert.NoError(t, err, "Expected no error when listing workflows")
}
