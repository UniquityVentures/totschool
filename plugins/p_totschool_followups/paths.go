package p_totschool_followups

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "followups.ListRoute", Value: lamu.Route{Path: AppUrl, Handler: lamu.NewDynamicView("followups.ListView")}},
			{Key: "followups.CreateRoute", Value: lamu.Route{Path: AppUrl + "create/", Handler: lamu.NewDynamicView("followups.CreateView")}},
			{Key: "followups.DetailRoute", Value: lamu.Route{Path: AppUrl + "{id}/", Handler: lamu.NewDynamicView("followups.DetailView")}},
			{Key: "followups.UpdateRoute", Value: lamu.Route{Path: AppUrl + "{id}/edit/", Handler: lamu.NewDynamicView("followups.UpdateView")}},
			{Key: "followups.DeleteRoute", Value: lamu.Route{Path: AppUrl + "{id}/delete/", Handler: lamu.NewDynamicView("followups.DeleteView")}},
			{Key: "followups.GenerateRoute", Value: lamu.Route{Path: AppUrl + "{id}/generate/", Handler: lamu.NewDynamicView("followups.GenerateView")}},
			{Key: "followups.CancelRoute", Value: lamu.Route{Path: AppUrl + "{id}/cancel/", Handler: lamu.NewDynamicView("followups.CancelView")}},
			{Key: "followups.AiEditFormRoute", Value: lamu.Route{Path: AppUrl + "{id}/ai-edit/form/", Handler: lamu.NewDynamicView("followups.AiEditFormView")}},
			{Key: "followups.AiEditRoute", Value: lamu.Route{Path: AppUrl + "{id}/ai-edit/", Handler: lamu.NewDynamicView("followups.AiEditView")}},
			{Key: "followups.ExportPdfRoute", Value: lamu.Route{Path: AppUrl + "{id}/export-pdf/", Handler: lamu.NewDynamicView("followups.ExportPdfView")}},
			{Key: "followups.ExportDocxRoute", Value: lamu.Route{Path: AppUrl + "{id}/export-docx/", Handler: lamu.NewDynamicView("followups.ExportDocxView")}},
			{Key: "followups.UserSelectRoute", Value: lamu.Route{Path: AppUrl + "users/select/", Handler: lamu.NewDynamicView("followups.UserSelectView")}},
		},
	}
}
