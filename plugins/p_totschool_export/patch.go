package p_totschool_export

import (
	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_export"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/registry"
	"github.com/UniquityVentures/lamu/views"
)

// exportAppRoles: dashboard tile + HTTP views ([RoleAuthorizationLayer] still
// allows IsSuperuser). AppsGrid skips role filter when $role is superuser.
var exportAppRoles = []string{"totschool_admin"}

// exportMenuRoles: SidebarMenu uses components.Render ($role string); superuser
// must appear here or sidebar stays empty while export routes still work.
var exportMenuRoles = []string{"totschool_admin", "superuser"}

var exportRoleLayer = p_users.RoleAuthorizationLayer{Roles: exportAppRoles}

// ExportPluginForTotschool returns core p_export with Totschool dashboard roles.
func ExportPluginForTotschool() registry.Pair[string, lamu.Plugin] {
	pair := p_export.GetPlugin()
	pair.Value.Roles = append([]string(nil), exportAppRoles...)
	return pair
}

func exportPagePatches() lamu.PluginFeatures[components.PageInterface] {
	return lamu.PluginFeatures[components.PageInterface]{
		Patches: []registry.Pair[string, func(components.PageInterface) components.PageInterface]{
			{
				Key: "export.Menu",
				Value: func(page components.PageInterface) components.PageInterface {
					menu, ok := page.(components.SidebarMenu)
					if !ok {
						return page
					}
					menu.Roles = append([]string(nil), exportMenuRoles...)
					return menu
				},
			},
			{
				Key: "export.Page",
				Value: func(page components.PageInterface) components.PageInterface {
					shell, ok := page.(*components.ShellScaffold)
					if !ok {
						return page
					}
					shell.Roles = append([]string(nil), exportMenuRoles...)
					return shell
				},
			},
		},
	}
}

func exportViewPatches() lamu.PluginFeatures[*views.View] {
	return lamu.PluginFeatures[*views.View]{
		Patches: []registry.Pair[string, func(*views.View) *views.View]{
			{
				Key: "export.PageView",
				Value: func(v *views.View) *views.View {
					return v.InsertLayerAfter("p_users.auth", "totschool_export.role", exportRoleLayer)
				},
			},
			{
				Key: "export.DownloadView",
				Value: func(v *views.View) *views.View {
					return v.InsertLayerAfter("p_users.auth", "totschool_export.role", exportRoleLayer)
				},
			},
		},
	}
}

// GetPlugin returns export UI patches (roles on p_export come from ExportPluginForTotschool).
func GetPlugin() registry.Pair[string, lamu.Plugin] {
	return registry.Pair[string, lamu.Plugin]{
		Key: "p_totschool_export",
		Value: lamu.Plugin{
			Type:  lamu.PluginTypeAddon,
			Pages: lamu.PluginStages(exportPagePatches),
			Views: lamu.PluginStages(exportViewPatches),
		},
	}
}
