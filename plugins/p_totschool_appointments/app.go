package p_totschool_appointments

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

const AppUrl = "/appointments/"

// GetPlugin returns registry contributions for [lamu.BuildAllRegistries].
// Appointments is an addon (not a dashboard app); timeline and the admin list
// are reached from the Clients sidebar menu.
func GetPlugin() registry.Pair[string, lamu.Plugin] {
	return registry.Pair[string, lamu.Plugin]{
		Key: "p_totschool_appointments",
		Value: lamu.Plugin{
			Type:        lamu.PluginTypeAddon,
			Views:       lamu.PluginStages(pluginViews),
			Pages:       lamu.PluginStages(pluginPages),
			Routes:      lamu.PluginStages(pluginRoutes),
			Models:      lamu.PluginStages(pluginModels),
			Migrations:  lamu.PluginStages(pluginMigrations),
			Configs:     lamu.PluginStages(pluginConfigs),
			Generators:  lamu.PluginStages(pluginGenerators),
			DBInitHooks: lamu.PluginStages(pluginDBInitHooks),
		},
	}
}
