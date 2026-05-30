package p_totschool_proposals

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

const AppUrl = "/proposals/"

// GetPlugin returns registry contributions for [lamu.BuildAllRegistries].
// Proposals is an addon (not a dashboard app); unassigned proposals are reached
// from the Clients sidebar menu.
func GetPlugin() registry.Pair[string, lamu.Plugin] {
	return registry.Pair[string, lamu.Plugin]{
		Key: "p_totschool_proposals",
		Value: lamu.Plugin{
			Type:        lamu.PluginTypeAddon,
			Views:       lamu.PluginStages(pluginViews),
			Pages:       lamu.PluginStages(pluginPages),
			Routes:      lamu.PluginStages(pluginRoutes),
			Models:      lamu.PluginStages(pluginModels),
			Migrations:  lamu.PluginStages(pluginMigrations),
			DBInitHooks: lamu.PluginStages(pluginDBInitHooks),
			Configs:     lamu.PluginStages(pluginConfigs),
		},
	}
}
