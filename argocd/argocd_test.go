package argocd

// import (
// 	"context"
// 	"log"
// 	"os"
// 	"testing"

// 	"github.com/argoproj/argo-cd/pkg/apiclient/application"
// 	"github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
// 	"github.com/argoproj/gitops-engine/pkg/health"
// 	"github.com/stretchr/testify/assert"
// )

// var argoClient *ArgoClient

// func TestMain(m *testing.M) {
// 	// Setup ArgoClient for integration tests
// 	address := os.Getenv("ARGOCD_SERVER_ADDRESS")
// 	token := os.Getenv("ARGOCD_TOKEN")

// 	if address == "" || token == "" {
// 		log.Fatal("Environment variables ARGOCD_SERVER_ADDRESS and ARGOCD_TOKEN are required")
// 	}

// 	argoClient = &ArgoClient{}
// 	err := argoClient.NewArgoClient(&ArgoConnection{
// 		Address: address,
// 		Token:   token,
// 	})
// 	if err != nil {
// 		log.Fatalf("Failed to create Argo client: %v", err)
// 	}

// 	// Run tests
// 	code := m.Run()
// 	os.Exit(code)
// }

// func TestGetProjects(t *testing.T) {
// 	projectName := "test-project"
// 	// Check if the project exists
// 	project, err := argoClient.GetProject(projectName)
// 	assert.NoError(t, err)
// 	assert.Equal(t, projectName, project.Name)
// }

// // func TestCreateAndDeleteProject(t *testing.T) {
// // 	// Create a test project
// // 	projectName := "test-project"
// // 	project := &v1alpha1.AppProject{
// // 		ObjectMeta: metav1.ObjectMeta{Name: projectName},
// // 		Spec:       v1alpha1.AppProjectSpec{Description: "Test project for integration testing"},
// // 	}

// // 	_, err := argoClient.projectClient.Create(context.Background(), project)
// // 	assert.NoError(t, err)

// // 	// Verify project creation
// // 	retrievedProject, err := argoClient.GetProject(projectName)
// // 	assert.NoError(t, err)
// // 	assert.Equal(t, projectName, retrievedProject.Name)

// // 	// Delete the project
// // 	err = argoClient.DeleteProject(projectName)
// // 	assert.NoError(t, err)

// // 	// Verify deletion
// // 	_, err = argoClient.GetProject(projectName)
// // 	assert.Error(t, err, "Expected an error since project should be deleted")
// // }

// func TestGetClusters(t *testing.T) {
// 	clusters, err := argoClient.GetClusters()
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, clusters)
// }

// func TestSyncApp(t *testing.T) {
// 	appName := "test-app"

// 	// Ensure application exists in ArgoCD
// 	_, err := argoClient.GetApplications()
// 	assert.NoError(t, err)

// 	// Trigger a sync operation
// 	err = argoClient.SyncApp(appName, 1) // 1-minute timeout
// 	assert.NoError(t, err, "Sync should complete without errors")

// 	// Verify sync status
// 	app, err := argoClient.appClient.Get(context.Background(), &application.ApplicationQuery{Name: &appName})
// 	assert.NoError(t, err)
// 	assert.Equal(t, v1alpha1.SyncStatusCodeSynced, app.Status.Sync.Status)
// 	assert.Equal(t, health.HealthStatusHealthy, app.Status.Health.Status)
// }

// func TestGetApplicationHistory(t *testing.T) {
// 	appName := "test-app"
// 	history, err := argoClient.GetApplicationHistory(appName)
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, history)
// }

// func TestGetLatestRevisionID(t *testing.T) {
// 	appName := "test-app"
// 	latestRevisionID := argoClient.GetLatestRevisionID(appName)
// 	assert.NotZero(t, latestRevisionID)
// }
