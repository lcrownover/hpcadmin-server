package config

import (
	"os"
	"os/user"
)

type CliOptions struct {
	// AzureAuthOptions
	TenantID string
	ClientID string

	// HPCAdmin Config
	ConfigDir string
}

func EnsureConfigDir() (string, error) {
	var dir string
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	// Look config dir in this order:
	// 1. ~/.config/hpcadmin/
	dir = usr.HomeDir + "/.config/hpcadmin"
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return dir, nil
	}
	// 2. ~/.hpcadmin
	dir = usr.HomeDir + "/.hpcadmin"
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return dir, nil
	}
	// Neither are found, so let's create the first one
	dir = usr.HomeDir + "/.config/hpcadmin"
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return "", err
	}
	return dir, nil
}
