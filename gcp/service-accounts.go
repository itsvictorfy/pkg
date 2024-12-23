package gcp

import (
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/api/iam/v1"
)

// Retrive all the keys for a service account
func (gcp *Gcp) GetServiceAccountKeys(saEmail string) ([]*iam.ServiceAccountKey, error) { //V
	slog.Info("GCP: retrieving Keys", slog.String("SaEmail", saEmail))
	resource := "projects/-/serviceAccounts/" + saEmail
	response, err := gcp.IamService.Projects.ServiceAccounts.Keys.List(resource).Do()
	if err != nil {
		return nil, fmt.Errorf("gcp: unable to retrieve keys for %s: %v", saEmail, err)
	}
	slog.Info("GCP: retrieved keys successfully", slog.String("SaEmail", saEmail), slog.String("Env", gcp.Env))
	return response.Keys, nil
}

// Retrieve all the service accounts
func (gcp *Gcp) GetAllServiceAccounts() ([]*iam.ServiceAccount, error) { //V
	response, err := gcp.IamService.Projects.ServiceAccounts.List("projects/" + gcp.ProjectID).Do()
	if err != nil {
		return nil, fmt.Errorf("gcp: unable to retrieve service accounts: %v", err)
	}
	return response.Accounts, nil
}

// Create a service account
func (gcp *Gcp) CreateServiceAccountKey(saEmail string) (*iam.ServiceAccountKey, error) { //V
	resource := "projects/-/serviceAccounts/" + saEmail
	request := &iam.CreateServiceAccountKeyRequest{
		KeyAlgorithm:   "KEY_ALG_RSA_2048",
		PrivateKeyType: "TYPE_GOOGLE_CREDENTIALS_FILE",
	}
	key, err := gcp.IamService.Projects.ServiceAccounts.Keys.Create(resource, request).Do()
	if err != nil {
		return nil, fmt.Errorf("gcp: unable to create key for: %w", err)
	}
	return key, nil
}

// Delete a service account key
func (gcp *Gcp) DeleteServiceAccountKey(saEmail string, key *iam.ServiceAccountKey) error { //V
	_, err := gcp.IamService.Projects.ServiceAccounts.Keys.Delete(key.Name).Do()
	if err != nil {
		return fmt.Errorf("gcp: unable to delete %s key: %v", saEmail, err)
	}
	return nil
}

// Check if a key is outdated
func IsKeyOutdated(key *iam.ServiceAccountKey) bool { //V
	slog.Info("GCP: Checking if Key is Outdated", slog.String("keyName", key.Name))
	parsedTime, err := time.Parse(time.RFC3339, key.ValidAfterTime)
	if err != nil {
		slog.Error("GCP: Error parsing ValidAfterTime", slog.String("keyName", key.Name), slog.String("Error:", err.Error()))
		return false
	}
	keyAge := time.Since(parsedTime)
	slog.Info("GCP:", slog.String("keyName", key.Name), slog.Bool("Key Outdated", keyAge > (90*24*time.Hour)))
	return keyAge > (90 * 24 * time.Hour)
}
