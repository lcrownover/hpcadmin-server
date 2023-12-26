package config

import (
	"path/filepath"
	"reflect"
	"testing"
)

// copied from tests/data/testconfig.yaml
//
// host: localhost
// port: 3333
// database:
//   host: localhost
//   port: 5432
//   user: hpcadmin
//   password: "superfancytestpasswordthatnobodyknows&"
//   dbname: hpcadmin_test
// oauth:
//   tenant_id: mock
//   client_id: mock
//   client_secret: mock
//

func TestLoadFile(t *testing.T) {
	// Test case 1: Test with valid config path
	t.Run("ValidConfigPath", func(t *testing.T) {
		configPath, _ := filepath.Abs("../../test/data/testconfig.yaml")
		want := &ServerConfig{
			Host: "localhost",
			Port: 3333,
			DB: DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "hpcadmin",
				Password: "superfancytestpasswordthatnobodyknows&",
				DBName:   "hpcadmin_test",
			},
			Oauth: OauthConfig{
				TenantID:     "mock",
				ClientID:     "mock",
				ClientSecret: "mock",
			},
		}

		got, err := LoadFile(configPath)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Compare the loaded config with the expected config
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Loaded config does not match expected config. got %+v want %+v", got, want)
		}
	})
	// Test case 2: Test with no config path
	t.Run("NoConfigPath", func(t *testing.T) {
		configPath := ""
		want := &ServerConfig{}

		got, err := LoadFile(configPath)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Compare the loaded config with the expected config
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Loaded config does not match expected config. got %+v want %+v", got, want)
		}
	})
}
