package config

type ServerOptions struct {
	Azure AzureOptions    `yaml:"azure"`
	DB    DatabaseOptions `yaml:"database"`
}

type AzureOptions struct {
	TenantID     string `yaml:"tenant_id"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type DatabaseOptions struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

func LoadConfig()
