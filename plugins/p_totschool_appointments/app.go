package p_totschool_appointments

import (
	"log"
	"net/url"

	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

const AppUrl = "/appointments/"

// GetPlugin returns registry contributions for [lamu.BuildAllRegistries].
func GetPlugin() registry.Pair[string, lamu.Plugin] {
	u, err := url.Parse(AppUrl)
	if err != nil {
		log.Panic(err)
	}
	return registry.Pair[string, lamu.Plugin]{
		Key: "p_totschool_appointments",
		Value: lamu.Plugin{
			Type:        lamu.PluginTypeApp,
			Icon:        "calendar-days",
			URL:         u,
			VerboseName: "Appointments",
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
