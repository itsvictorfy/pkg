package github

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v66/github"
)

type Github struct {
	Client     *github.Client
	Token      string
	OrgName    string
	ActionId   int64
	Repo       string
	ActionFile string
}

func (gh *Github) CreateGithubClient() error {
	if gh.Token == "" {
		log.Fatal("github: token is required")
	}
	repoaut := github.RepositoryListByOrgOptions{}
	gh.Client = github.NewClient(nil).WithAuthToken(gh.Token)
	repo, resp, err := gh.Client.Repositories.ListByOrg(context.Background(), gh.OrgName, &repoaut)
	if repo == nil && resp.StatusCode == 200 {
		return fmt.Errorf("github: unable to create github client: %v", err)
	}
	return nil

}

func (gh *Github) TriggerWorkflowActionByID(req github.CreateWorkflowDispatchEventRequest) error {
	ctx := context.Background()
	_, err := gh.Client.Actions.CreateWorkflowDispatchEventByID(ctx, gh.OrgName, gh.Repo, gh.ActionId, req)
	if err != nil {
		return fmt.Errorf("github: unable to trigger workflow action: %v", err)
	}
	return nil
}

func (gh *Github) TriggerWorkflowActionByFileName(req github.CreateWorkflowDispatchEventRequest) error {
	ctx := context.Background()
	data, err := gh.Client.Actions.CreateWorkflowDispatchEventByFileName(ctx, gh.OrgName, gh.Repo, gh.ActionFile, req)
	if err != nil {
		return fmt.Errorf("github: unable to trigger workflow action: %v", err)
	}
	fmt.Printf("Workflow Dispatch Event Response: %v\n", data)
	return nil
}

func (gh *Github) ListWorkflows() error {
	ctx := context.Background()
	wf, _, err := gh.Client.Actions.ListWorkflows(ctx, gh.OrgName, gh.Repo, nil)
	if err != nil {
		return fmt.Errorf("github: unable to list workflows: %v", err)
	}
	for _, workflow := range wf.Workflows {
		fmt.Printf("Workflow Name: %v\n Workflow ID: %v, WorkFlowFile: %v", workflow.GetName(), workflow.GetID(), workflow.Path)
	}
	return nil
}
