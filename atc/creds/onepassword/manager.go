package onepassword

import (
	"context"
	"fmt"

	"code.cloudfoundry.org/lager/v3"
	"github.com/1password/onepassword-sdk-go"
	"github.com/concourse/concourse/atc/creds"
	"github.com/go-viper/mapstructure/v2"
)

type OnePasswordManager struct {
	Token     string `long:"token" description:"1Password Service Account Token" required:"false"`
	VaultName string `long:"vault-name" description:"Vault name to use when looking up secrets." default:"Infrastructure"`
	Client    *onepassword.Client
}

func (manager *OnePasswordManager) Init(log lager.Logger) error {
	var err error

	manager.Client, err = onepassword.NewClient(
		context.TODO(),
		onepassword.WithServiceAccountToken(manager.Token),
		onepassword.WithIntegrationInfo("Concourse CI Credential Provider", "v1.0.0"),
	)

	return err
}

func (manager *OnePasswordManager) Config(config map[string]any) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		ErrorUnused: true,
		Result:      &manager,
	})
	if err != nil {
		return err
	}

	err = decoder.Decode(config)

	return err
}

func (manager *OnePasswordManager) IsConfigured() bool {
	return manager.Token != ""
}

func (manager *OnePasswordManager) Health() (*creds.HealthResponse, error) {
	health := &creds.HealthResponse{
		Method: "Vault List",
	}

	vaults, err := manager.Client.Vaults().List(context.Background())

	if err != nil {
		health.Error = err.Error()
	} else {
		health.Response = fmt.Sprintf("Found %d vaults", len(vaults))
	}

	return health, nil
}

func (manager *OnePasswordManager) Close(logger lager.Logger) {}

func (manager *OnePasswordManager) NewSecretsFactory(logger lager.Logger) (creds.SecretsFactory, error) {
	return &OnePasswordFactory{
		Manager:   *manager,
		VaultName: manager.VaultName,
	}, nil
}

func (manager *OnePasswordManager) Validate() error {
	vaults, err := manager.Client.Vaults().List(context.Background())
	if err != nil {
		return err
	}

	for _, v := range vaults {
		if v.Title == manager.VaultName {
			return nil
		}
	}

	return fmt.Errorf("vault %s not found", manager.VaultName)
}
