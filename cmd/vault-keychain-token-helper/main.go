package main

import (
	"io/ioutil"
	"net/url"
	"os"

	vaultkeychain "github.com/janslow/vault-keychain"
	"github.com/sirupsen/logrus"
)

func init() {
	l := logrus.WarnLevel
	if os.Getenv("VAULT_KEYCHAIN_DEBUG") == "true" {
		l = logrus.DebugLevel
	}
	logrus.StandardLogger().SetLevel(l)
}
func main() {
	s := vaultkeychain.Server{Address: vaultAddress()}
	l := logrus.WithField("address", s.Address.String())
	l.Debug("Parsed VAULT_ADDR")

	c := command()
	var err error
	switch c {
	case "get":
		l.Debug("Reading Token")
		var t string
		t, err = s.Token()
		if err == vaultkeychain.ErrTokenNotFound {
			l.WithError(err).Debug("No token found in Keychain")
			err = nil
		} else if err == nil {
			l.Debug("Token found in Keychain")
			_, err = os.Stdout.WriteString(t)
		}
	case "store":
		t := readLine()
		l.Debug("Setting Token")
		err = s.SetToken(t)
	case "erase":
		l.Debug("Clearing Token")
		err = s.ClearToken()
	default:
		logrus.WithField("command", c).Fatal("Unknown command")
	}
	if err != nil {
		l.WithError(err).Fatal("Keychain access failed")
	}
}

func command() string {
	if len(os.Args) != 2 {
		logrus.Fatalf("Usage: %s get|store|erase", os.Args[1])
	}
	return os.Args[1]
}

func readLine() string {
	logrus.Debug("Reading token from stdin")
	l, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to read token from stdin")
	}
	return string(l)
}

func vaultAddress() *url.URL {
	raw := os.Getenv("VAULT_ADDR")
	if raw == "" {
		logrus.Fatal("VAULT_ADDR not set")
	}
	u, err := url.Parse(raw)
	if err != nil {
		logrus.WithError(err).Fatal("VAULT_ADDR is not a valid URL")
	}
	return u
}
