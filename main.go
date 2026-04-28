package main

import (
	"log/slog"

	"github.com/UniquityVentures/lago/lago"
	_ "github.com/UniquityVentures/lago/plugins/p_dashboard"
	_ "github.com/UniquityVentures/lago/plugins/p_export"
	_ "github.com/UniquityVentures/lago/plugins/p_livereloading"
	_ "github.com/UniquityVentures/lago/plugins/p_otp"
	_ "github.com/UniquityVentures/lago/plugins/p_pwa"
	_ "github.com/UniquityVentures/lago/plugins/p_users"
	_ "github.com/UniquityVentures/totschool_lago/plugins/p_totschool_appointments"
	_ "github.com/UniquityVentures/totschool_lago/plugins/p_totschool_export"
	_ "github.com/UniquityVentures/totschool_lago/plugins/p_totschool_proposals"
	_ "github.com/UniquityVentures/totschool_lago/plugins/p_totschool_tally"
	_ "github.com/UniquityVentures/totschool_lago/plugins/p_totschool_users"
)

func main() {
	config, err := lago.LoadConfigFromFile("totschool.toml")
	if err != nil {
		panic(err)
	}
	if err := lago.Start(config); err != nil {
		slog.Error(err.Error())
	}
}
