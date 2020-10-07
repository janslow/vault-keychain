package vaultkeychain

import (
	"fmt"
	"net/url"

	"github.com/keybase/go-keychain"
	"github.com/sirupsen/logrus"
)

const (
	securityClass = "1.com.jayanslow.vault-keychain"
)

func write(address *url.URL, token string) error {
	i := item(address)

	logrus.WithField("VAULT_ADDR", address.String()).Debug("Attempting to clear existing Keychain entry")
	err := clear(address)
	if err != nil {
		return err
	}

	i.SetData([]byte(token))

	logrus.Info("Adding Vault token to Keychain")
	return keychain.AddItem(i)
}

func read(address *url.URL) (token string, err error) {
	i, err := query(address)
	if err != nil {
		return
	}
	token = string(i.Data)
	return
}

func clear(address *url.URL) error {
	err := keychain.DeleteItem(item(address))
	if err == nil {
		logrus.Info("Deleted existing Keychain entry")
		return nil
	}
	if err == keychain.ErrorItemNotFound {
		logrus.WithError(err).Debug("No existing Keychain entry")
		return nil
	}
	logrus.WithError(err).Error("Failed to clear Keychain entry")
	return fmt.Errorf("failed to clear Keychain entry")
}

func query(address *url.URL) (_ keychain.QueryResult, err error) {
	q := item(address)
	q.SetMatchLimit(keychain.MatchLimitOne)
	q.SetReturnAttributes(true)
	q.SetReturnData(true)

	logrus.Debug("Reading Vault token from Keychain")
	results, err := keychain.QueryItem(q)
	if err == keychain.Error(-128) {
		err = fmt.Errorf("keychain access denied: %w", err)
		return
	} else if err != nil {
		err = fmt.Errorf("failed to read from Keychain: %w", err)
		return
	}
	if len(results) == 0 {
		err = ErrTokenNotFound
		return
	}

	return results[0], nil
}

func item(address *url.URL) keychain.Item {
	q := keychain.NewItem()
	q.SetSecClass(keychain.SecClassGenericPassword)
	q.SetAccessGroup(securityClass)
	q.SetLabel(fmt.Sprintf("Hashicorp Vault (%s)", address))
	q.SetService(address.String())
	q.SetAccount("token")
	return q
}
