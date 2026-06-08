package main

import (
	"log/slog"

	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"

	"github.com/UniquityVentures/lamu/plugins/p_dashboard"
	"github.com/UniquityVentures/lamu/plugins/p_google_genai"
	"github.com/UniquityVentures/lamu/plugins/p_livereloading"
	"github.com/UniquityVentures/lamu/plugins/p_pwa"
	"github.com/UniquityVentures/totschool/plugins/p_totschool_appointments"
	"github.com/UniquityVentures/totschool/plugins/p_totschool_clients"
	"github.com/UniquityVentures/totschool/plugins/p_totschool_dashboard"
	"github.com/UniquityVentures/totschool/plugins/p_totschool_export"
	"github.com/UniquityVentures/totschool/plugins/p_totschool_followups"
	"github.com/UniquityVentures/totschool/plugins/p_totschool_proposals"
	"github.com/UniquityVentures/totschool/plugins/p_totschool_tally"
	"github.com/UniquityVentures/totschool/plugins/p_totschool_users"
)

func main() {
	plugins := []registry.Pair[string, lamu.Plugin]{
		p_dashboard.GetPlugin(),
		p_google_genai.GetPlugin(),
		p_totschool_export.ExportPluginForTotschool(),
		p_totschool_clients.GetPlugin(),
		p_totschool_users.UsersPluginForTotschool(),
		p_totschool_users.OtpPluginForTotschool(),
		p_totschool_appointments.GetPlugin(),
		p_totschool_proposals.GetPlugin(),
		p_totschool_followups.GetPlugin(),
		p_totschool_tally.GetPlugin(),
		p_totschool_dashboard.GetPlugin(),
		p_totschool_export.GetPlugin(),
		p_totschool_users.GetPlugin(),
		p_livereloading.GetPlugin(),
		p_pwa.GetPlugin(),
	}
	config, err := lamu.LoadConfigFromFile("totschool.toml", plugins)
	if err != nil {
		panic(err)
	}
	if err := lamu.Start(config, plugins); err != nil {
		slog.Error(err.Error())
	}
}
