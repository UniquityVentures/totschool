package p_totschool_followups

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

type FollowupAIConfig struct {
	APIKey string `toml:"apiKey"`
	Model  string `toml:"model"`
}

var followupAIConfig = &FollowupAIConfig{}

func (c *FollowupAIConfig) PostConfig() {}

func pluginConfigs() lamu.PluginFeatures[lamu.Config] {
	return lamu.PluginFeatures[lamu.Config]{
		Entries: []registry.Pair[string, lamu.Config]{
			{Key: "p_totschool_followups", Value: followupAIConfig},
		},
	}
}
