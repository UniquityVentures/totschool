package p_totschool_appointments

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "appointments.CardTimelineRoute", Value: lamu.Route{Path: AppUrl, Handler: lamu.NewDynamicView("appointments.CardTimelineView")}},
			{Key: "appointments.ListRoute", Value: lamu.Route{Path: AppUrl + "list/", Handler: lamu.NewDynamicView("appointments.ListView")}},
			{Key: "appointments.CreateRoute", Value: lamu.Route{Path: AppUrl + "create/", Handler: lamu.NewDynamicView("appointments.CreateView")}},
			{Key: "appointments.DetailRoute", Value: lamu.Route{Path: AppUrl + "{id}/", Handler: lamu.NewDynamicView("appointments.DetailView")}},
			{Key: "appointments.UpdateRoute", Value: lamu.Route{Path: AppUrl + "{id}/edit/", Handler: lamu.NewDynamicView("appointments.UpdateView")}},
			{Key: "appointments.DeleteRoute", Value: lamu.Route{Path: AppUrl + "{id}/delete/", Handler: lamu.NewDynamicView("appointments.DeleteView")}},
			{Key: "appointments.GenerateRoute", Value: lamu.Route{Path: AppUrl + "{id}/generate/", Handler: lamu.NewDynamicView("appointments.GenerateView")}},
			{Key: "appointments.CancelRoute", Value: lamu.Route{Path: AppUrl + "{id}/cancel/", Handler: lamu.NewDynamicView("appointments.CancelView")}},
			{Key: "appointments.AiEditFormRoute", Value: lamu.Route{Path: AppUrl + "{id}/ai-edit/form/", Handler: lamu.NewDynamicView("appointments.AiEditFormView")}},
			{Key: "appointments.AiEditRoute", Value: lamu.Route{Path: AppUrl + "{id}/ai-edit/", Handler: lamu.NewDynamicView("appointments.AiEditView")}},
			{Key: "appointments.SelectRoute", Value: lamu.Route{Path: AppUrl + "select/", Handler: lamu.NewDynamicView("appointments.SelectView")}},
			{Key: "appointments.UserSelectRoute", Value: lamu.Route{Path: AppUrl + "users/select/", Handler: lamu.NewDynamicView("appointments.UserSelectView")}},
		},
	}
}
