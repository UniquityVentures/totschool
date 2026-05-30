package p_totschool_dashboard

import (
	"log"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func dashboardPagePatches() lamu.PluginFeatures[components.PageInterface] {
	return lamu.PluginFeatures[components.PageInterface]{
		Patches: []registry.Pair[string, func(components.PageInterface) components.PageInterface]{
			{
				Key:   "dashboard.AppsPage",
				Value: patchDashboardAppsPage,
			},
		},
	}
}

func patchDashboardAppsPage(page components.PageInterface) components.PageInterface {
	scaffold, ok := page.(*components.ShellTopbarScaffold)
	if !ok {
		log.Panic("dashboard.AppsPage was not *components.ShellTopbarScaffold")
	}
	components.ReplaceChild(scaffold, "dashboard.AppsPageLayout", func(layout *components.LayoutSimple) *components.LayoutSimple {
		if len(layout.Children) == 1 && layout.Children[0].GetKey() == appsDashboardKey {
			return layout
		}
		layout.Children = []components.PageInterface{
			&totschoolAppsDashboard{Page: components.Page{Key: appsDashboardKey}},
		}
		return layout
	})
	return scaffold
}

// GetPlugin returns dashboard UI patches for Totschool.
func GetPlugin() registry.Pair[string, lamu.Plugin] {
	return registry.Pair[string, lamu.Plugin]{
		Key: "p_totschool_dashboard",
		Value: lamu.Plugin{
			Type:  lamu.PluginTypeAddon,
			Pages: lamu.PluginStages(dashboardPagePatches),
		},
	}
}
