package argocd

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type ArgoConnection struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
// 	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
// 	"github.com/argoproj/argo-cd/v2/pkg/apiclient/cluster"
// 	"github.com/argoproj/argo-cd/v2/pkg/apiclient/project"
// 	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
// 	"github.com/argoproj/gitops-engine/pkg/health"
// )

// // Argocd Details for connection
// type ArgoConnection struct {
// 	Address string
// 	Token   string
// }

// // Argocd Clients
// type ArgoClient struct {
// 	projectClient project.ProjectServiceClient
// 	clusterClient cluster.ClusterServiceClient
// 	appClient     application.ApplicationServiceClient
// }

// // NewArgoClient creates a new ArgoCD clients
// func (argocd *ArgoClient) NewArgoClient(c *ArgoConnection) error {
// 	apiClient, err := apiclient.NewClient(&apiclient.ClientOptions{
// 		ServerAddr: fmt.Sprintf(c.Address),
// 		Insecure:   true,
// 		AuthToken:  c.Token,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	_, argocd.projectClient, err = apiClient.NewProjectClient()
// 	if err != nil {
// 		return err
// 	}

// 	_, argocd.clusterClient, err = apiClient.NewClusterClient()
// 	if err != nil {
// 		return err
// 	}

// 	_, argocd.appClient, err = apiClient.NewApplicationClient()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Delete ArgoCD Project
// func (c *ArgoClient) DeleteProject(name string) error {
// 	_, err := c.projectClient.Delete(context.Background(), &project.ProjectQuery{
// 		Name: name,
// 	})
// 	log.Printf("Deleted Project %s", name)
// 	return err
// }

// // Get ArgoCD Project by Name
// func (c *ArgoClient) GetProject(name string) (*v1alpha1.AppProject, error) {
// 	log.Printf("Getting Project %s", name)
// 	return c.projectClient.Get(context.Background(), &project.ProjectQuery{
// 		Name: name,
// 	})
// }

// // // Get ArgoCD Project by Name
// // func (c *Client) GetAllProjects(name string) (*v1alpha1.AppProject, error) {
// // 	log.Printf("Getting Project %s", name)
// // 	return c.projectClient.list(context.Background(), c.projectClient.List(context.Background(), &project.ProjectQuery{}))
// // }

// // Get All ArgoCD Clusters
// func (c *ArgoClient) GetClusters() ([]v1alpha1.Cluster, error) {
// 	cl, err := c.clusterClient.List(context.Background(), &cluster.ClusterQuery{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	log.Printf("Getting Clusters")
// 	return cl.Items, nil
// }

// // Get All ArgoCD Applications
// func (c *ArgoClient) GetApplications() ([]v1alpha1.Application, error) {
// 	apps, err := c.appClient.List(context.Background(), &application.ApplicationQuery{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	log.Printf("Getting Applications")
// 	return apps.Items, nil
// }

// // Get ArgoCD Application History
// func (c *ArgoClient) GetApplicationHistory(appName string) (v1alpha1.RevisionHistories, error) {
// 	app, err := c.appClient.Get(context.Background(), &application.ApplicationQuery{
// 		Name: &appName,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get application: %w", err)
// 	}

// 	return app.Status.History, nil
// }

// // Sync ArgoCD Application
// func (c *ArgoClient) SyncApp(name string, syncTime int) error { //need to figure out edge case where Kustomize updated to a new version while current sync is in progress
// 	log.Printf("Syncing Application: %s", name)

// 	_, err := c.appClient.Sync(context.Background(), &application.ApplicationSyncRequest{
// 		Name: &name,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to initiate sync: %w", err)
// 	}
// 	time.Sleep(10 * time.Second)
// 	// Define the timeout duration
// 	timeout := time.Duration(syncTime) * time.Minute
// 	timeoutChan := time.After(timeout)
// 	refreshMode := "normal"
// 	for {
// 		select {
// 		case <-timeoutChan:
// 			log.Printf("Timeout reached for application: %s", name)
// 			return fmt.Errorf("timeout reached for application: %s", name)
// 		default:
// 			app, err := c.appClient.Get(context.Background(), &application.ApplicationQuery{
// 				Name:    &name,
// 				Refresh: &refreshMode,
// 			})
// 			if err != nil {
// 				return fmt.Errorf("failed to get application status: %w", err)
// 			}

// 			if app.Status.Sync.Status == v1alpha1.SyncStatusCodeSynced &&
// 				app.Status.Health.Status == health.HealthStatusHealthy {
// 				log.Printf("Application %s is synced and healthy", name)
// 				return nil
// 			}
// 			if app.Status.Health.Status == health.HealthStatusDegraded {
// 				log.Printf("Application %s is degraded", name)
// 				return fmt.Errorf("application %s is degraded", name)
// 			}
// 			time.Sleep(5 * time.Second)
// 		}
// 	}
// }

// // GetLatestRevisionID returns the latest revision ID for an application
// func (c *ArgoClient) GetLatestRevisionID(appName string) int64 {
// 	history, err := c.GetApplicationHistory(appName)
// 	if err != nil {
// 		log.Printf("Error getting application history: %v", err)
// 		return 0
// 	}

// 	if len(history) == 0 {
// 		log.Printf("No history found for application: %s", appName)
// 		return 0
// 	}

// 	// The latest revision should be the first entry in the history list
// 	latestRevision := history[len(history)-2].ID

//		return latestRevision
//	}
func (argo *ArgoConnection) APISyncApp(appName string) error {
	slog.Info("Argocd:Syncing ArgoCd Application", slog.String("App Name", appName))

	url := fmt.Sprintf("%s/api/v1/applications/%s/sync", argo.Address, appName)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("argocd: Error Marshaling syncRequest: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+argo.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("argocd: Error send HTTP request: %s", err)
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("argocd: Error reading HTTP response body: %s", err)
	}

	slog.Info("Argocd: App Synced Successfuly", slog.String("appName", appName), slog.String("ResponseStatus", resp.Status))

	return nil
}
func (argo *ArgoConnection) APIGetAppStatus(appName string) (string, error) {
	url := fmt.Sprintf("%s/api/v1/applications/%s", argo.Address, appName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("argocd: Error creating HTTP request: %s", err)
	}

	req.Header.Set("Authorization", "Bearer "+argo.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("argocd: Error sending HTTP request: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("argocd: Error reading HTTP response body: %s", err)
	}
	fmt.Println(string(body))

	// Unmarshal the JSON response into a struct
	var appStatus struct {
		Status struct {
			OperationState struct {
				Phase string `json:"phase"`
			} `json:"operationState"`
		} `json:"status"`
	}

	err = json.Unmarshal(body, &appStatus)
	if err != nil {
		return "", fmt.Errorf("argocd: Error unmarshaling response: %s", err)
	}

	return appStatus.Status.OperationState.Phase, nil
}
