package github

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v66/github"
)

type Github struct {
	Client     *github.Client `json:"-"` // Exclude from JSON
	Token      string         `json:"token"`
	OrgName    string         `json:"orgName"`
	ActionId   int64          `json:"actionId"`
	Repo       string         `json:"repo"`
	ActionFile string         `json:"actionFile"`
}

type Workflow struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

var (
	WorkflowCIBuildDeploy          = Workflow{Name: "CI-Build-Deploy", ID: 119110820}
	WorkflowCIManualSingleService  = Workflow{Name: "CI-Manual-SingleService", ID: 119110821}
	WorkflowCIOnCommitTests        = Workflow{Name: "CI-OnCommit-Tests", ID: 119110822}
	WorkflowCIBuildDeployHtmlToPdf = Workflow{Name: "CI-Build-Deploy-HtmlToPdf", ID: 119110823}
	WorkflowMaintenance            = Workflow{Name: "maintenance", ID: 119116978}
	WorkflowCIDeploy               = Workflow{Name: "CI-Deploy", ID: 125946370}
	WorkflowTestWorkflow           = Workflow{Name: "test-workflow", ID: 126133680}
	WorkflowCDDeployEnvironment    = Workflow{Name: "CD-Deploy-Environment", ID: 126723957}
)

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
