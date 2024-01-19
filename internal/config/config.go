package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Host  string         `yaml:"host"`
	Port  int            `yaml:"port"`
	Oauth OauthConfig    `yaml:"oauth"`
	DB    DatabaseConfig `yaml:"database"`
}

type OauthConfig struct {
	TenantID     string `yaml:"tenant_id"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

// Load loads the configuration from the given path
// If the path is empty, it will load the default configuration
// file from /etc/hpcadmin-server/config.yaml
func LoadFile(configPath string) (*ServerConfig, error) {
	var err error
	cfg := &ServerConfig{}
	// If configPath wasn't provided, and the file doesn't exist, just return an empty config
	if configPath == "" {
		configPath = "/etc/hpcadmin-server/config.yaml"
	}
	slog.Debug("configuration path found", "package", "config", "method", "Load", "path", configPath)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		slog.Debug("configuration file not found", "package", "config", "method", "Load", "path", configPath)
		return cfg, nil
	}

	slog.Debug("reading config file", "package", "config", "method", "Load", "path", configPath)
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %v", err)
	}

	slog.Debug("parsing YAML", "package", "config", "method", "Load", "path", configPath)
	err = yaml.Unmarshal(configData, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	return cfg, nil
}

func LoadEnvironment(cfg *ServerConfig) *ServerConfig {
	// HPCADMIN_SERVER_HOST
	if host, found := os.LookupEnv("HPCADMIN_SERVER_HOST"); found {
		slog.Debug("found host override", "package", "config", "method", "LoadEnvironment", "host", host)
		cfg.Host = host
	}
	// HPCADMIN_SERVER_PORT
	if port, found := os.LookupEnv("HPCADMIN_SERVER_PORT"); found {
		slog.Debug("found port override", "package", "config", "method", "LoadEnvironment", "port", port)
		iport, err := strconv.Atoi(port)
		if err != nil {
			slog.Warn("Invalid port number", "package", "config", "method", "LoadEnvironment", "port", port)
		} else {
			cfg.Port = iport
		}
	}
	// HPCADMIN_SERVER_DATABASE_HOST
	if dbhost, found := os.LookupEnv("HPCADMIN_SERVER_DATABASE_HOST"); found {
		slog.Debug("found database host override", "package", "config", "method", "LoadEnvironment", "host", dbhost)
		cfg.DB.Host = dbhost
	}
	// HPCADMIN_SERVER_DATABASE_PORT
	if dbport, found := os.LookupEnv("HPCADMIN_SERVER_DATABASE_PORT"); found {
		slog.Debug("found database port override", "package", "config", "method", "LoadEnvironment", "port", dbport)
		idbport, err := strconv.Atoi(dbport)
		if err != nil {
			slog.Warn("Invalid database port number", "package", "config", "method", "LoadEnvironment", "port", dbport)
		} else {
			cfg.DB.Port = idbport
		}
	}
	// HPCADMIN_SERVER_DATABASE_USER
	if dbuser, found := os.LookupEnv("HPCADMIN_SERVER_DATABASE_USER"); found {
		slog.Debug("found database user override", "package", "config", "method", "LoadEnvironment", "user", dbuser)
		cfg.DB.User = dbuser
	}
	// HPCADMIN_SERVER_DATABASE_PASSWORD
	if dbpassword, found := os.LookupEnv("HPCADMIN_SERVER_DATABASE_USER"); found {
		slog.Debug("found database user override", "package", "config", "method", "LoadEnvironment", "password", "REDACTED")
		cfg.DB.Password = dbpassword
	}
	// HPCADMIN_SERVER_DATABASE_DBNAME
	if dbname, found := os.LookupEnv("HPCADMIN_SERVER_DATABASE_DBNAME"); found {
		slog.Debug("found database user override", "package", "config", "method", "LoadEnvironment", "dbname", dbname)
		cfg.DB.DBName = dbname
	}
	// HPCADMIN_SERVER_OAUTH_TENANT_ID
	if tenantID, found := os.LookupEnv("HPCADMIN_SERVER_OAUTH_TENANT_ID"); found {
		slog.Debug("found oauth tenantID override", "package", "config", "method", "LoadEnvironment", "tenantID", tenantID)
		cfg.Oauth.TenantID = tenantID
	}
	// HPCADMIN_SERVER_OAUTH_CLIENT_ID
	if clientID, found := os.LookupEnv("HPCADMIN_SERVER_OAUTH_CLIENT_ID"); found {
		slog.Debug("found oauth clientID override", "package", "config", "method", "LoadEnvironment", "clientID", clientID)
		cfg.Oauth.ClientID = clientID
	}
	// HPCADMIN_SERVER_OAUTH_CLIENT_SECRET
	if clientSecret, found := os.LookupEnv("HPCADMIN_SERVER_OAUTH_CLIENT_SECRET"); found {
		slog.Debug("found oauth clientSecret override", "package", "config", "method", "LoadEnvironment", "clientSecret", "REDACTED")
		cfg.Oauth.ClientSecret = clientSecret
	}
	return cfg
}

func Validate(cfg *ServerConfig) error {
	if cfg.Host == "" {
		return fmt.Errorf("missing host")
	}
	if cfg.Port == 0 {
		return fmt.Errorf("missing port")
	}
	if cfg.DB.Host == "" {
		return fmt.Errorf("missing database host")
	}
	if cfg.DB.Port == 0 {
		return fmt.Errorf("missing database port")
	}
	if cfg.DB.User == "" {
		return fmt.Errorf("missing database user")
	}
	if cfg.DB.Password == "" {
		return fmt.Errorf("missing database password")
	}
	if cfg.DB.DBName == "" {
		return fmt.Errorf("missing database name")
	}
	if cfg.Oauth.TenantID == "" {
		return fmt.Errorf("missing oauth tenant ID")
	}
	if cfg.Oauth.ClientID == "" {
		return fmt.Errorf("missing oauth client ID")
	}
	if cfg.Oauth.ClientSecret == "" {
		return fmt.Errorf("missing oauth client secret")
	}
	return nil
}
