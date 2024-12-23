package gcp

import (
	"context"
	"fmt"

	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// Retrieve the secret from the secret manager
func (gcp *Gcp) GetSecret(name string) (string, error) {
	secretName := "projects/" + gcp.ProjectID + "/secrets/" + name + "/versions/latest"
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	}
	result, err := gcp.SecretManagerService.AccessSecretVersion(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("gcp: unable to get secret: %v, Error:%v", name, err)
	}
	secretData := string(result.Payload.Data)
	return secretData, nil
}

// Update the secret in the secret manager
func (gcp *Gcp) UpdateSecret(name, secretPayload string) error {
	secretName := "projects/" + gcp.ProjectID + "/secrets/" + name + "/versions/latest"
	req := &secretmanagerpb.AddSecretVersionRequest{
		Parent: secretName,
		Payload: &secretmanagerpb.SecretPayload{
			Data: []byte(secretPayload),
		},
	}
	_, err := gcp.SecretManagerService.AddSecretVersion(context.Background(), req)
	if err != nil {
		return fmt.Errorf("gcp: unable to update secret: %v, Error: %v", name, err)
	}
	return nil

}

// Create a new secret in the secret manager
func (gcp *Gcp) CreateSecret(name, secretPayload string) error {
	newSecretReq := &secretmanagerpb.AddSecretVersionRequest{
		Parent: "projects/" + gcp.ProjectID + "/secrets/website-maintenance-config",
		Payload: &secretmanagerpb.SecretPayload{
			Data: []byte(secretPayload),
		},
	}
	_, err := gcp.SecretManagerService.AddSecretVersion(context.TODO(), newSecretReq)
	if err != nil {
		return fmt.Errorf("gcp: unable to create secret: %v, Error: %v", name, err)
	}
	return nil
}
