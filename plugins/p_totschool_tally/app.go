package p_totschool_tally

import (
	"log"
	"net/url"

	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

const AppUrl = "/tally/"

// GetPlugin returns registry contributions for [lamu.BuildAllRegistries].
func GetPlugin() registry.Pair[string, lamu.Plugin] {
	u, err := url.Parse(AppUrl)
	if err != nil {
		log.Panic(err)
	}
	return registry.Pair[string, lamu.Plugin]{
		Key: "p_totschool_tally",
		Value: lamu.Plugin{
			Type:        lamu.PluginTypeApp,
			Icon:        "chart-bar",
			URL:         u,
			VerboseName: "Progress Tracker",
			Views:       pluginViews,
			Pages:       pluginPages,
			Routes:      pluginRoutes,
			Models:      pluginModels,
			Migrations:  pluginMigrations,
		},
	}
}
