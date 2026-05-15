package p_totschool_users

import (
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/views"
)

// userViewsWithRoleLayer are p_users views that enforce "p_users.role".
var userViewsWithRoleLayer = []string{
	"p_users.ListView", "p_users.DetailView", "p_users.CreateView", "p_users.UpdateView",
	"p_users.DeleteView", "p_users.ChangePasswordView", "p_users.SelectView",
	"p_users.RoleSelectView", "p_users.RoleListView", "p_users.RoleDetailView",
	"p_users.RoleCreateView", "p_users.RoleUpdateView", "p_users.RoleDeleteView",
}

func userRolePatcher(current views.Layer) views.Layer {
	return p_users.RoleAuthorizationLayer{Roles: []string{"", "totschool_admin"}}
}
