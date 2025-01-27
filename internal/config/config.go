package config

import (
	"fmt"
	baseConfig "github.com/avalance-rl/otiva-pkg/config"
	"github.com/spf13/viper"
)

type AuthConfig struct {
	*baseConfig.Config
	JWT struct {
		Secret string
		TTL    int
	}
}

func Load(path string) (*AuthConfig, error) {
	baseCFG, err := baseConfig.Load(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load base config: %w", err)
	}

	config := &AuthConfig{
		Config: baseCFG,
	}

	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read auth config: %w", err)
	}

	if err := viper.UnmarshalKey("jwt", &config.JWT); err != nil {
		return nil, fmt.Errorf("failed to unmarshal jwt config: %w", err)
	}
	return config, nil
}
