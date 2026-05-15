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
			Views:       pluginViews,
			Pages:       pluginPages,
			Routes:      pluginRoutes,
			Models:      pluginModels,
			Migrations:  pluginMigrations,
			DBInitHooks: pluginDBInitHooks,
			Configs:     pluginConfigs,
		},
	}
}
