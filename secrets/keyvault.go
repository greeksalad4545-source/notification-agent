package secrets

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

type KeyVault struct {
	client *azsecrets.Client
}

func NewKeyVault(vaultName string) (*KeyVault, error) {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create Azure credential: %w",
			err,
		)
	}

	vaultURL := fmt.Sprintf(
		"https://%s.vault.azure.net/",
		vaultName,
	)

	client, err := azsecrets.NewClient(
		vaultURL,
		credential,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create Key Vault client: %w",
			err,
		)
	}

	return &KeyVault{
		client: client,
	}, nil
}

func (k *KeyVault) GetSecret(
	ctx context.Context,
	secretName string,
) (string, error) {
	response, err := k.client.GetSecret(
		ctx,
		secretName,
		"",
		nil,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get secret %s: %w",
			secretName,
			err,
		)
	}

	if response.Value == nil {
		return "", fmt.Errorf(
			"secret %s has no value",
			secretName,
		)
	}

	return *response.Value, nil
}
