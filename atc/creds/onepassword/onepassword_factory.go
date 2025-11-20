package onepassword

import "github.com/concourse/concourse/atc/creds"

type OnePasswordFactory struct {
	VaultName string
	Manager   OnePasswordManager
}

func NewOnePasswordFactory(vaultName string, manager OnePasswordManager) OnePasswordFactory {
	return OnePasswordFactory{
		VaultName: vaultName,
		Manager:   manager,
	}
}

func (factory OnePasswordFactory) NewSecrets() creds.Secrets {
	return &OnePassword{
		Manager:   factory.Manager,
		VaultName: factory.VaultName,
	}
}
