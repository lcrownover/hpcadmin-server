package config

import (
	"fmt"
	"log/slog"
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

func EnsureCLIConfigDir() (string, error) {
	slog.Debug("ensuring config directory", "method", "EnsureConfigDir")
	var dir string
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	// Look config dir in this order:
	// 1. ~/.config/hpcadmin/
	dir = usr.HomeDir + "/.config/hpcadmin"
	slog.Debug(fmt.Sprintf("checking for config dir: %s", dir), "method", "EnsureConfigDir")
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		slog.Debug("found config dir", "method", "EnsureConfigDir")
		return dir, nil
	}
	// 2. ~/.hpcadmin
	dir = usr.HomeDir + "/.hpcadmin"
	slog.Debug(fmt.Sprintf("checking for config dir: %s", dir), "method", "EnsureConfigDir")
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		slog.Debug("found config dir", "method", "EnsureConfigDir")
		return dir, nil
	}
	// Neither are found, so let's create the first one
	dir = usr.HomeDir + "/.config/hpcadmin"
	slog.Debug(fmt.Sprintf("creating config dir: %s", dir), "method", "EnsureConfigDir")
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return "", err
	}
	return dir, nil
}
