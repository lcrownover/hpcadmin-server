package config

import (
	"reflect"
	"testing"
)

func TestLoadFile(t *testing.T) {
	// Test case 1: Test with valid config path
	t.Run("ValidConfigPath", func(t *testing.T) {
		configPath := "tests/data/testconfig.yaml"
		expectedConfig := &ServerConfig{
			Host: "localhost",
			Port: 3333,
			DB: DatabaseConfig{
				Host: "localhost",
				Port: 5432,
			},
		}

		config, err := LoadFile(configPath)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Compare the loaded config with the expected config
		if !reflect.DeepEqual(config, expectedConfig) {
			t.Errorf("Loaded config does not match expected config")
		}
	})
	// Test case 2: Test with no config path
	t.Run("NoConfigPath", func(t *testing.T) {
		configPath := ""
		expectedConfig := &ServerConfig{}

		config, err := LoadFile(configPath)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Compare the loaded config with the expected config
		if !reflect.DeepEqual(config, expectedConfig) {
			t.Errorf("Loaded config does not match expected config")
		}
	})
}
