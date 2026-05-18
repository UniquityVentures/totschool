package p_totschool_clients

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "clients.ListRoute", Value: lamu.Route{Path: AppUrl, Handler: lamu.NewDynamicView("clients.ListView")}},
			{Key: "clients.CreateRoute", Value: lamu.Route{Path: AppUrl + "create/", Handler: lamu.NewDynamicView("clients.CreateView")}},
			{Key: "clients.DetailRoute", Value: lamu.Route{Path: AppUrl + "{id}/", Handler: lamu.NewDynamicView("clients.DetailView")}},
			{Key: "clients.UpdateRoute", Value: lamu.Route{Path: AppUrl + "{id}/edit/", Handler: lamu.NewDynamicView("clients.UpdateView")}},
			{Key: "clients.DeleteRoute", Value: lamu.Route{Path: AppUrl + "{id}/delete/", Handler: lamu.NewDynamicView("clients.DeleteView")}},
			{Key: "clients.SelectRoute", Value: lamu.Route{Path: AppUrl + "select/", Handler: lamu.NewDynamicView("clients.SelectView")}},
		},
	}
}
