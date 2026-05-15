package p_totschool_proposals

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

type AIConfig struct {
	APIKey string `toml:"apiKey"`
	Model  string `toml:"model"`
}

var aiConfig = &AIConfig{}

func (c *AIConfig) PostConfig() {}

func pluginConfigs() lamu.PluginFeatures[lamu.Config] {
	return lamu.PluginFeatures[lamu.Config]{
		Entries: []registry.Pair[string, lamu.Config]{
			{Key: "p_totschool_proposals", Value: aiConfig},
		},
	}
}
