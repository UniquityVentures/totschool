package p_totschool_proposals

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "proposals.ListRoute", Value: lamu.Route{Path: AppUrl, Handler: lamu.NewDynamicView("proposals.ListView")}},
			{Key: "proposals.CreateRoute", Value: lamu.Route{Path: AppUrl + "create/", Handler: lamu.NewDynamicView("proposals.CreateView")}},
			{Key: "proposals.DetailRoute", Value: lamu.Route{Path: AppUrl + "{id}/", Handler: lamu.NewDynamicView("proposals.DetailView")}},
			{Key: "proposals.UpdateRoute", Value: lamu.Route{Path: AppUrl + "{id}/edit/", Handler: lamu.NewDynamicView("proposals.UpdateView")}},
			{Key: "proposals.DeleteRoute", Value: lamu.Route{Path: AppUrl + "{id}/delete/", Handler: lamu.NewDynamicView("proposals.DeleteView")}},
			{Key: "proposals.GenerateRoute", Value: lamu.Route{Path: AppUrl + "{id}/generate/", Handler: lamu.NewDynamicView("proposals.GenerateView")}},
			{Key: "proposals.CancelRoute", Value: lamu.Route{Path: AppUrl + "{id}/cancel/", Handler: lamu.NewDynamicView("proposals.CancelView")}},
			{Key: "proposals.AiEditFormRoute", Value: lamu.Route{Path: AppUrl + "{id}/ai-edit/form/", Handler: lamu.NewDynamicView("proposals.AiEditFormView")}},
			{Key: "proposals.AiEditRoute", Value: lamu.Route{Path: AppUrl + "{id}/ai-edit/", Handler: lamu.NewDynamicView("proposals.AiEditView")}},
			{Key: "proposals.ExportPdfRoute", Value: lamu.Route{Path: AppUrl + "{id}/export-pdf/", Handler: lamu.NewDynamicView("proposals.ExportPdfView")}},
			{Key: "proposals.ExportDocxRoute", Value: lamu.Route{Path: AppUrl + "{id}/export-docx/", Handler: lamu.NewDynamicView("proposals.ExportDocxView")}},
		},
	}
}
