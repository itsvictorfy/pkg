package github

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v66/github"
)

type Github struct {
	Client  *github.Client `json:"-"` // Exclude from JSON
	Token   string         `json:"token"`
	OrgName string         `json:"orgName"`
	Repo    string         `json:"repo"`
}

type Workflow struct {
	Name       string                 `json:"name"`
	ID         int64                  `json:"id"`
	ActionFile string                 `json:"actionFile"`
	Inputs     map[string]interface{} `json:"inputs"`
	Ref        string                 `json:"ref"`
}

func (gh *Github) InitGithubClient() error {
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

func (gh *Github) TriggerWorkflowActionByID(wf Workflow) error {
	ctx := context.Background()
	req := github.CreateWorkflowDispatchEventRequest{
		Ref:    wf.Ref,
		Inputs: wf.Inputs,
	}

	_, err := gh.Client.Actions.CreateWorkflowDispatchEventByID(ctx, gh.OrgName, gh.Repo, wf.ID, req)
	if err != nil {
		return fmt.Errorf("github: unable to trigger workflow action: %v", err)
	}
	return nil
}

func (gh *Github) TriggerWorkflowActionByFileName(wf Workflow) error {
	ctx := context.Background()
	req := github.CreateWorkflowDispatchEventRequest{
		Ref:    wf.Ref,
		Inputs: wf.Inputs,
	}
	data, err := gh.Client.Actions.CreateWorkflowDispatchEventByFileName(ctx, gh.OrgName, gh.Repo, wf.ActionFile, req)
	if err != nil {
		return fmt.Errorf("github: unable to trigger workflow action: %v", err)
	}
	fmt.Printf("Workflow Dispatch Event Response: %v\n", data)
	return nil
}

func (gh *Github) ListAllWorkflows() ([]*github.Workflow, error) {
	ctx := context.Background()
	wf, _, err := gh.Client.Actions.ListWorkflows(ctx, gh.OrgName, gh.Repo, nil)
	if err != nil {
		return []*github.Workflow{}, fmt.Errorf("github: unable to list workflows: %v", err)
	}
	for _, workflow := range wf.Workflows {
		fmt.Printf("Workflow Name: %v\n Workflow ID: %v\n", workflow.GetName(), workflow.GetID())
	}
	return wf.Workflows, nil
}
