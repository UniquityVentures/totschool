package p_totschool_followups

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

const AppUrl = "/followups/"

// GetPlugin returns registry contributions for [lamu.BuildAllRegistries].
// Followups is an addon reached from the Clients sidebar and client detail page.
func GetPlugin() registry.Pair[string, lamu.Plugin] {
	return registry.Pair[string, lamu.Plugin]{
		Key: "p_totschool_followups",
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
