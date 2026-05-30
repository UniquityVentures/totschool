package p_totschool_appointments

import (
	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/registry"
	"github.com/UniquityVentures/lamu/views"
)

var appointmentListAdminRoles = []string{"totschool_admin", "superuser"}

var (
	pluginPagePatches []registry.Pair[string, func(components.PageInterface) components.PageInterface]
	pluginViewPatches []registry.Pair[string, func(*views.View) *views.View]
)

func init() {
	registerAppointmentListAdminPatch()
}

func registerAppointmentListAdminPatch() {
	patchPluginView("appointments.ListView", func(v *views.View) *views.View {
		return v.InsertLayerAfter("p_users.auth", "totschool_appointments.list_admin", p_users.RoleAuthorizationLayer{Roles: appointmentListAdminRoles})
	})
}

func patchPluginPage(key string, patch func(components.PageInterface) components.PageInterface) {
	pluginPagePatches = append(pluginPagePatches, registry.Pair[string, func(components.PageInterface) components.PageInterface]{
		Key: key, Value: patch,
	})
}

func patchPluginView(key string, patch func(*views.View) *views.View) {
	pluginViewPatches = append(pluginViewPatches, registry.Pair[string, func(*views.View) *views.View]{
		Key: key, Value: patch,
	})
}

func pluginPagesWithPatches(entries []registry.Pair[string, components.PageInterface]) lamu.PluginFeatures[components.PageInterface] {
	return lamu.PluginFeatures[components.PageInterface]{
		Entries: entries,
		Patches: pluginPagePatches,
	}
}

func pluginViewsWithPatches(entries []registry.Pair[string, *views.View]) lamu.PluginFeatures[*views.View] {
	return lamu.PluginFeatures[*views.View]{
		Entries: entries,
		Patches: pluginViewPatches,
	}
}
