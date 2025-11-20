package onepassword

import (
	"fmt"

	"github.com/concourse/concourse/atc/creds"
	"github.com/jessevdk/go-flags"
)

type onepasswordManagerFactory struct{}

func init() {
	creds.Register("onepassword", NewManagerFactory())
}

func NewManagerFactory() creds.ManagerFactory {
	return &onepasswordManagerFactory{}
}

func (factory *onepasswordManagerFactory) NewInstance(config any) (creds.Manager, error) {
	configMap, ok := config.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid onepassword credential manager config: %T", config)
	}

	manager := &OnePasswordManager{}

	err := manager.Config(configMap)
	if err != nil {
		return nil, err
	}

	return manager, nil
}

func (factory *onepasswordManagerFactory) AddConfig(group *flags.Group) creds.Manager {
	manager := &OnePasswordManager{}

	subGroup, err := group.AddGroup("1Password Credential Management", "", manager)
	if err != nil {
		panic(err)
	}

	subGroup.Namespace = "onepassword"

	return manager
}
