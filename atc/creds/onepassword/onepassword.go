package onepassword

import (
	"context"
	"time"

	"github.com/concourse/concourse/atc/creds"
)

type OnePassword struct {
	Manager   OnePasswordManager
	VaultName string
}

func (op *OnePassword) NewSecretLookupPaths(teamName string, pipelineName string, allowRootPath bool) []creds.SecretLookupPath {
	prefix := "op://" + op.VaultName

	//if !allowRootPath {
	//	return []creds.SecretLookupPath{}
	//}

	return []creds.SecretLookupPath{creds.NewSecretLookupWithPrefix(prefix + "/")}
}

func (op OnePassword) Get(secretPath string) (any, *time.Time, bool, error) {
	print("1Password fetch: " + secretPath + "\n")

	secret, err := op.Manager.Client.Secrets().Resolve(
		context.Background(),
		secretPath,
	)

	if err != nil {
		if err.Error() == "error: error resolving secret reference: no item matched the secret reference query" {
			return nil, nil, false, nil
		} else {
			return nil, nil, false, err
		}
	}

	return secret, nil, true, nil
}
