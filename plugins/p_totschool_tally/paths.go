package p_totschool_tally

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "tally.TallyListRoute", Value: lamu.Route{
				Path:    "/tally/list/",
				Handler: lamu.NewDynamicView("tally.TallyListView"),
			}},
			{Key: "tally.TallyDashboardRoute", Value: lamu.Route{
				Path:    "/tally/",
				Handler: lamu.NewDynamicView("tally.TallyDashboardView"),
			}},
			{Key: "tally.TallyLeaderboardRoute", Value: lamu.Route{
				Path:    "/tally/leaderboard/",
				Handler: lamu.NewDynamicView("tally.TallyLeaderboardView"),
			}},
			{Key: "tally.TallyDailyFormRoute", Value: lamu.Route{
				Path:    "/tally/daily/",
				Handler: lamu.NewDynamicView("tally.TallyDailyFormView"),
			}},
			{Key: "tally.TallyCreateRoute", Value: lamu.Route{
				Path:    "/tally/create/",
				Handler: lamu.NewDynamicView("tally.TallyCreateView"),
			}},
			{Key: "tally.TallyUpdateRoute", Value: lamu.Route{
				Path:    "/tally/{id}/update/",
				Handler: lamu.NewDynamicView("tally.TallyUpdateView"),
			}},
			{Key: "tally.TallyDeleteRoute", Value: lamu.Route{
				Path:    "/tally/{id}/delete/",
				Handler: lamu.NewDynamicView("tally.TallyDeleteView"),
			}},
			{Key: "tally.TallyDetailRoute", Value: lamu.Route{
				Path:    "/tally/{id}/",
				Handler: lamu.NewDynamicView("tally.TallyDetailView"),
			}},
		},
	}
}
