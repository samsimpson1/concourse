package onepassword

import (
	"fmt"

	"code.cloudfoundry.org/lager/v3"
	"github.com/1Password/connect-sdk-go/connect"
	"github.com/concourse/concourse/atc/creds"
	"github.com/go-viper/mapstructure/v2"
)

type OnePasswordManager struct {
	ConnectHost  string `long:"connect-host" description:"1Password Connect Host" required:"false"`
	ConnectToken string `long:"connect-token" description:"1Password Connect Token" required:"false"`
	VaultName    string `long:"vault-name" description:"Vault name to use when looking up secrets." default:"Infrastructure"`
	Client       connect.Client
}

func (manager *OnePasswordManager) Init(log lager.Logger) error {
	var err error

	manager.Client = connect.NewClient(manager.ConnectHost, manager.ConnectToken)

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
	return manager.ConnectHost != "" && manager.ConnectToken != ""
}

func (manager *OnePasswordManager) Health() (*creds.HealthResponse, error) {
	health := &creds.HealthResponse{
		Method: "Vault List",
	}

	vaults, err := manager.Client.GetVaults()

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
	_, err := manager.Client.GetVaultByTitle(manager.VaultName)

	if err != nil {
		return err
	}

	return nil
}
