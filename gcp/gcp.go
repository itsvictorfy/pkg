package gcp

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	storage "cloud.google.com/go/storage"
	storagetransfer "cloud.google.com/go/storagetransfer/apiv1"
	"cloud.google.com/go/storagetransfer/apiv1/storagetransferpb"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	cloudrun "google.golang.org/api/run/v1"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

type Gcp struct {
	IamService            *iam.Service            `json:"-"` // Exclude from JSON
	SecretManagerService  *secretmanager.Client   `json:"-"` // Exclude from JSON
	SqlService            *sqladmin.Service       `json:"-"` // Exclude from JSON
	storageTransferClient *storagetransfer.Client `json:"-"` // Exclude from JSON
	CloudRunService       *cloudrun.APIService    `json:"-"` // Exclude from JSON
	StorageService        *storage.Client         `json:"-"` // Exclude from JSON
	Env                   string                  `json:"env"`
	ProjectID             string                  `json:"projectId"`
}

// Initc initializes the c client
func (c *Gcp) InitGcp() error { //V
	slog.Info("GCP: init client", slog.String("Env", c.Env))
	credentials := "config/sa.json"
	gcpContext := context.Background()

	var err error
	c.IamService, err = iam.NewService(gcpContext, option.WithCredentialsFile(credentials))
	if err != nil {
		return fmt.Errorf("gcp: unable to create iam service client: %v", err)
	}
	c.SecretManagerService, err = secretmanager.NewClient(gcpContext, option.WithCredentialsFile(credentials))
	if err != nil {
		return fmt.Errorf("gcp: unable to create secretManager service client: %v", err)
	}
	c.SqlService, err = sqladmin.NewService(gcpContext, option.WithCredentialsFile(credentials))
	if err != nil {
		return fmt.Errorf("gcp: unable to create sql service client: %v", err)
	}
	c.storageTransferClient, err = storagetransfer.NewClient(context.Background(), option.WithCredentialsFile(credentials))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	c.CloudRunService, err = cloudrun.NewService(gcpContext, option.WithCredentialsFile(credentials))
	if err != nil {
		return fmt.Errorf("gcp: unable to create cloudrun service client: %v", err)
	}
	c.StorageService, err = storage.NewClient(gcpContext, option.WithCredentialsFile(credentials))
	if err != nil {
		return fmt.Errorf("gcp: unable to create storage service client: %v", err)
	}
	slog.Info("GCP: client initialized", slog.String("Project", c.ProjectID), slog.String("Env", c.Env))
	return nil
}

// Triggers c Transfer Job
func (c *Gcp) TriggerTransferJob(projectID, transferJobName string) error {
	ctx := context.Background()
	req := &storagetransferpb.RunTransferJobRequest{
		JobName:   transferJobName,
		ProjectId: projectID,
	}

	_, err := c.storageTransferClient.RunTransferJob(ctx, req)
	if err != nil {
		return fmt.Errorf("client.RunTransferJob: %v", err)
	}

	log.Printf("Transfer job triggered successfully")
	return nil
}

func (c *Gcp) DeployCloudRunApp(serviceName, imageURL, region string) (string, string, error) {
	// Define the parent location in Cloud Run
	parent := fmt.Sprintf("projects/%s/locations/%s", c.ProjectID, region)

	// Define the Cloud Run service specification
	service := &cloudrun.Service{
		ApiVersion: "serving.knative.dev/v1",
		Kind:       "Service",
		Metadata: &cloudrun.ObjectMeta{
			Name: serviceName,
		},
		Spec: &cloudrun.ServiceSpec{
			Template: &cloudrun.RevisionTemplate{
				Spec: &cloudrun.RevisionSpec{
					Containers: []*cloudrun.Container{
						{
							Image: imageURL,
							Ports: []*cloudrun.ContainerPort{
								{
									ContainerPort: 8080, // Cloud Run requires port 8080
								},
							},
						},
					},
				},
			},
		},
	}

	// Deploy the service to Cloud Run
	app, err := c.CloudRunService.Projects.Locations.Services.Create(parent, service).Do()
	if err != nil {
		return "", "", fmt.Errorf("failed to deploy service to Cloud Run: %v", err)
	}

	fmt.Printf("Service %s deployed successfully to region %s!\n", serviceName, region)
	return app.Metadata.Uid, app.Status.Url, nil
}

// UploadFileToGCS uploads a local file to a specified c bucket
func (c *Gcp) UploadFileToGCS(bucketName, objectName, filePath string) error {

	// Open the local file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	// Create a writer for the bucket
	wc := c.StorageService.Bucket(bucketName).Object(objectName).NewWriter(context.Background())
	defer func() {
		if closeErr := wc.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close writer for bucket %s: %v", bucketName, closeErr)
		}
	}()

	// Copy the file content to the bucket
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("failed to write file to bucket %s: %v", bucketName, err)
	}

	fmt.Printf("File %s uploaded to bucket %s as object %s\n", filePath, bucketName, objectName)
	return nil
}
