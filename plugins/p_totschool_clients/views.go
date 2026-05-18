package p_totschool_clients

import (
	"net/http"

	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/registry"
	"github.com/UniquityVentures/lamu/views"
	"gorm.io/gorm"
)

type clientQueryPatcher struct{}

func (clientQueryPatcher) Patch(_ views.View, r *http.Request, query gorm.ChainInterface[Client]) gorm.ChainInterface[Client] {
	user, role := p_users.UserAndRoleFromContext(r.Context(), "clientQueryPatcher")
	if user.IsSuperuser || role == "totschool_admin" {
		return query
	}
	return query.Where("created_by_id = ?", user.ID)
}

type clientFormPatcher struct{}

func (clientFormPatcher) Patch(_ views.View, r *http.Request, formData map[string]any, formErrors map[string]error) (map[string]any, map[string]error) {
	user := p_users.UserFromContext(r.Context(), "clientFormPatcher")
	formData["CreatedByID"] = user.ID
	return formData, formErrors
}

func pluginViews() lamu.PluginFeatures[*views.View] {
	return lamu.PluginFeatures[*views.View]{
		Entries: []registry.Pair[string, *views.View]{
			{Key: "clients.ListView", Value: lamu.GetPageView("clients.ClientTable").
				WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
				WithLayer("clients.list", views.LayerList[Client]{
					Key: getters.Static("clients"),
					QueryPatchers: views.QueryPatchers[Client]{
						{Key: "clients.query", Value: clientQueryPatcher{}},
					},
				})},
			{Key: "clients.DetailView", Value: lamu.GetPageView("clients.ClientDetail").
				WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
				WithLayer("clients.detail", views.LayerDetail[Client]{
					Key:          getters.Static("client"),
					PathParamKey: getters.Static("id"),
				})},
			{Key: "clients.CreateView", Value: lamu.GetPageView("clients.ClientCreateForm").
				WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
				WithLayer("clients.create", views.LayerCreate[Client]{
					SuccessURL: lamu.RoutePath("clients.DetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$id")),
					}),
					FormPatchers: views.FormPatchers{
						{Key: "clients.form", Value: clientFormPatcher{}},
					},
				})},
			{Key: "clients.UpdateView", Value: lamu.GetPageView("clients.ClientUpdateForm").
				WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
				WithLayer("clients.detail", views.LayerDetail[Client]{
					Key:          getters.Static("client"),
					PathParamKey: getters.Static("id"),
				}).
				WithLayer("clients.update", views.LayerUpdate[Client]{
					Key: getters.Static("client"),
					SuccessURL: lamu.RoutePath("clients.DetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("client.ID")),
					}),
					FormPatchers: views.FormPatchers{
						{Key: "clients.form", Value: clientFormPatcher{}},
					},
				})},
			{Key: "clients.DeleteView", Value: lamu.GetPageView("clients.ClientDeleteForm").
				WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
				WithLayer("clients.detail", views.LayerDetail[Client]{
					Key:          getters.Static("client"),
					PathParamKey: getters.Static("id"),
				}).
				WithLayer("clients.delete", views.LayerDelete[Client]{
					Key:        getters.Static("client"),
					SuccessURL: lamu.RoutePath("clients.ListRoute", nil),
				})},
			{Key: "clients.SelectView", Value: lamu.GetPageView("clients.ClientSelectionTable").
				WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
				WithLayer("clients.select_list", views.LayerList[Client]{
					Key: getters.Static("clients"),
					QueryPatchers: views.QueryPatchers[Client]{
						{Key: "clients.query", Value: clientQueryPatcher{}},
					},
				})},
		},
	}
}
