package p_totschool_users

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/registry"
)

// UsersPluginForTotschool returns p_users with Totschool dashboard visibility roles.
func UsersPluginForTotschool() registry.Pair[string, lamu.Plugin] {
	pair := p_users.GetPlugin()
	pair.Value.Roles = []string{"superuser", "totschool_admin"}
	return pair
}
