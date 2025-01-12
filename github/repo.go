package github

import (
	"context"
	"fmt"
	"log"
)

func (gh *Github) ListBranchsinRepo() ([]string, error) {
	ctx := context.Background()
	var branchList []string
	branches, _, err := gh.Client.Repositories.ListBranches(ctx, gh.OrgName, gh.Repo, nil)
	if err != nil {
		log.Fatalf("github: unable to list branches in repo: %v", err)
	}
	for _, branch := range branches {
		fmt.Println(*branch.Name)
		branchList = append(branchList, *branch.Name)
	}
	return branchList, nil
}
