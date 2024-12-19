package gcp

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	storagetransfer "cloud.google.com/go/storagetransfer/apiv1"
	"cloud.google.com/go/storagetransfer/apiv1/storagetransferpb"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

type Gcp struct {
	IamService            *iam.Service            `json:"-"` // Exclude from JSON
	SecretManagerService  *secretmanager.Client   `json:"-"` // Exclude from JSON
	SqlService            *sqladmin.Service       `json:"-"` // Exclude from JSON
	storageTransferClient *storagetransfer.Client `json:"-"` // Exclude from JSON
	Env                   string                  `json:"env"`
	ProjectID             string                  `json:"projectId"`
}

// InitGcp initializes the GCP client
func (gcp *Gcp) InitGcp() error { //V
	slog.Info("GCP: init client", slog.String("Env", gcp.Env))
	credentials := "config/sa.json"
	gcpContext := context.Background()

	var err error
	gcp.IamService, err = iam.NewService(gcpContext, option.WithCredentialsFile(credentials))
	if err != nil {
		return fmt.Errorf("gcp: unable to create iam service client: %v", err)
	}
	gcp.SecretManagerService, err = secretmanager.NewClient(gcpContext, option.WithCredentialsFile(credentials))
	if err != nil {
		return fmt.Errorf("gcp: unable to create secretManager service client: %v", err)
	}
	gcp.SqlService, err = sqladmin.NewService(gcpContext, option.WithCredentialsFile(credentials))
	if err != nil {
		return fmt.Errorf("gcp: unable to create sql service client: %v", err)
	}
	gcp.storageTransferClient, err = storagetransfer.NewClient(context.Background(), option.WithCredentialsFile(credentials))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	slog.Info("GCP: client initialized", slog.String("Project", gcp.ProjectID), slog.String("Env", gcp.Env))
	return nil
}

// Triggers GCP Transfer Job
func (gcp *Gcp) TriggerTransferJob(projectID, transferJobName string) error {
	ctx := context.Background()
	req := &storagetransferpb.RunTransferJobRequest{
		JobName:   transferJobName,
		ProjectId: projectID,
	}

	_, err := gcp.storageTransferClient.RunTransferJob(ctx, req)
	if err != nil {
		return fmt.Errorf("client.RunTransferJob: %v", err)
	}

	log.Printf("Transfer job triggered successfully")
	return nil
}
