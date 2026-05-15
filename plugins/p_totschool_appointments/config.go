package p_totschool_appointments

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

type TotscholAppointmentsConfig struct {
	APIKey string `toml:"apiKey"`
	Model  string `toml:"model"`
}

var totschoolAppointmentConfig = &TotscholAppointmentsConfig{}

func (c *TotscholAppointmentsConfig) PostConfig() {}

func pluginConfigs() lamu.PluginFeatures[lamu.Config] {
	return lamu.PluginFeatures[lamu.Config]{
		Entries: []registry.Pair[string, lamu.Config]{
			{Key: "p_totschool_appointments", Value: totschoolAppointmentConfig},
		},
	}
}
