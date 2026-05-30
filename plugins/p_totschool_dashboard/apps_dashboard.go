package p_totschool_dashboard

import (
	"context"
	"fmt"
	"slices"
	"sort"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

const appsDashboardKey = "totschool_dashboard.AppsDashboard"

// totschoolAppsDashboard replaces dashboard.AppsGrid with core apps plus an
// "Admin only" section for PluginTypeApp entries that declare Roles.
type totschoolAppsDashboard struct {
	components.Page
}

func (e totschoolAppsDashboard) GetKey() string     { return e.Key }
func (e totschoolAppsDashboard) GetRoles() []string { return e.Roles }

func (e totschoolAppsDashboard) Build(ctx context.Context) Node {
	coreApps, adminApps := totschoolDashboardApps(ctx)

	sections := Group{appTileGrid(coreApps, ctx)}
	if len(adminApps) > 0 {
		sections = append(sections,
			Div(Class("mt-8 mb-4"), H2(Class("text-lg font-semibold"), Text("Admin only"))),
			appTileGrid(adminApps, ctx),
		)
	}

	return Div(Class("container max-w-5xl mx-auto mt-4 @container"), Attr("x-data", "{ search: ''}"),
		Div(Class("mb-4"),
			Input(Type("text"), Attr("x-model", "search"), Placeholder("Search apps..."), Class("input input-bordered w-full")),
		),
		sections,
	)
}

func totschoolDashboardApps(ctx context.Context) (core, admin []lamu.Plugin) {
	pluginsMap := lamu.RegistryPlugin.AllStable()
	roleName := p_users.RoleFromContext(ctx, appsDashboardKey)

	for _, pluginItem := range *pluginsMap {
		plugin := pluginItem.Value
		if plugin.Type != lamu.PluginTypeApp {
			continue
		}
		if !totschoolDashboardAppVisible(roleName, plugin) {
			continue
		}
		if len(plugin.Roles) > 0 {
			admin = append(admin, plugin)
		} else {
			core = append(core, plugin)
		}
	}

	sortTotschoolDashboardApps(core)
	sortTotschoolDashboardApps(admin)
	return core, admin
}

func totschoolDashboardAppVisible(roleName string, plugin lamu.Plugin) bool {
	if roleName == "superuser" || len(plugin.Roles) == 0 {
		return true
	}
	return slices.Contains(plugin.Roles, roleName)
}

func sortTotschoolDashboardApps(apps []lamu.Plugin) {
	sort.Slice(apps, func(i, j int) bool {
		return apps[i].VerboseName < apps[j].VerboseName
	})
}

func appTileGrid(apps []lamu.Plugin, ctx context.Context) Node {
	group := Group{}
	for _, app := range apps {
		group = append(group, A(
			Href(app.URL.String()),
			Class("btn btn-md h-auto flex-col space-y-1 py-4"),
			Attr("x-show", fmt.Sprintf("'%s'.toLowerCase().includes(search.toLowerCase())", app.VerboseName)),
			Attr("x-cloak"), components.Render(components.Icon{Name: app.Icon, Classes: "w-8 h-8"}, ctx), Div(
				Class("text-sm truncate min-w-0 w-full"),
				Text(app.VerboseName),
			),
		))
	}
	return Div(Class("grid grid-cols-2 @md:grid-cols-4 @2xl:grid-cols-6 gap-2"), group)
}
