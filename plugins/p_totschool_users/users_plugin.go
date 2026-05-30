package p_totschool_users

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_otp"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/registry"
)

var adminDashboardRoles = []string{"superuser", "totschool_admin"}

// UsersPluginForTotschool returns p_users with Totschool dashboard visibility roles.
func UsersPluginForTotschool() registry.Pair[string, lamu.Plugin] {
	pair := p_users.GetPlugin()
	pair.Value.Roles = append([]string(nil), adminDashboardRoles...)
	return pair
}

// OtpPluginForTotschool returns p_otp with Totschool admin-only dashboard visibility.
func OtpPluginForTotschool() registry.Pair[string, lamu.Plugin] {
	pair := p_otp.GetPlugin()
	pair.Value.Roles = append([]string(nil), adminDashboardRoles...)
	return pair
}
