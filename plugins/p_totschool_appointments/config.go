package p_totschool_appointments

import "github.com/UniquityVentures/lago/lago"

type TotscholAppointmentsConfig struct {
	APIKey string `toml:"apiKey"`
	Model  string `toml:"model"`
}

var totschoolAppointmentConfig = &TotscholAppointmentsConfig{}

func (c *TotscholAppointmentsConfig) PostConfig() {}

func init() {
	lago.RegistryConfig.Register("p_totschool_appointments", totschoolAppointmentConfig)
}
