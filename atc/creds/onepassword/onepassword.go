package onepassword

import (
	"strings"
	"time"

	"github.com/1Password/connect-sdk-go/onepassword"
	"github.com/concourse/concourse/atc/creds"
)

type OnePassword struct {
	Manager   OnePasswordManager
	VaultName string
}

func (op *OnePassword) NewSecretLookupPaths(teamName string, pipelineName string, allowRootPath bool) []creds.SecretLookupPath {
	return []creds.SecretLookupPath{creds.NewSecretLookupWithPrefix("")}
}

func findFieldByName(item *onepassword.Item, fieldName string) *onepassword.ItemField {
	for _, field := range item.Fields {
		if field.Label == fieldName || field.ID == fieldName {
			return field
		}
	}
	return nil
}

func (op OnePassword) Get(secretPath string) (any, *time.Time, bool, error) {
	print("1Password fetch: " + secretPath + "\n")

	secretParts := strings.Split(secretPath, "/")

	fieldName := secretParts[len(secretParts)-1]

	secretName := strings.Join(secretParts[:len(secretParts)-1], "/")

	secret, err := op.Manager.Client.GetItem(
		secretName,
		op.VaultName,
	)

	if err != nil {
		if err.Error() == "idk" {
			return nil, nil, false, nil
		} else {
			return nil, nil, false, err
		}
	}

	field := findFieldByName(secret, fieldName)

	if field == nil {
		return nil, nil, false, nil
	}

	return field.Value, nil, true, nil
}
