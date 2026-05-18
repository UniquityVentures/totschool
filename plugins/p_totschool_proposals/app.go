package p_totschool_proposals

import (
	"log"
	"net/url"

	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

const AppUrl = "/proposals/"

// GetPlugin returns registry contributions for [lamu.BuildAllRegistries].
func GetPlugin() registry.Pair[string, lamu.Plugin] {
	u, err := url.Parse(AppUrl)
	if err != nil {
		log.Panic(err)
	}
	return registry.Pair[string, lamu.Plugin]{
		Key: "p_totschool_proposals",
		Value: lamu.Plugin{
			Type:        lamu.PluginTypeApp,
			Icon:        "document-text",
			URL:         u,
			VerboseName: "Proposals",
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
