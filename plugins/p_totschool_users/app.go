package p_totschool_users

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
	"github.com/UniquityVentures/lamu/views"
)

// GetPlugin patches p_users views to allow totschool_admin for role-gated routes.
func GetPlugin() registry.Pair[string, lamu.Plugin] {
	var patches []registry.Pair[string, func(*views.View) *views.View]
	for _, key := range userViewsWithRoleLayer {
		k := key
		patches = append(patches, registry.Pair[string, func(*views.View) *views.View]{
			Key: k,
			Value: func(v *views.View) *views.View {
				return v.PatchLayer("p_users.role", userRolePatcher)
			},
		})
	}
	return registry.Pair[string, lamu.Plugin]{
		Key: "p_totschool_users",
		Value: lamu.Plugin{
			Type: lamu.PluginTypeAddon,
			DBInitHooks: lamu.PluginStages(pluginDBInitHooks),
			Migrations:  lamu.PluginStages(pluginMigrations),
			Views: lamu.PluginStages(func() lamu.PluginFeatures[*views.View] {
				return lamu.PluginFeatures[*views.View]{Patches: patches}
			}),
		},
	}
}
