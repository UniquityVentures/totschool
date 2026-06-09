package p_totschool_appointments

import (
	"context"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginPages() lamu.PluginFeatures[components.PageInterface] {
	var entries []registry.Pair[string, components.PageInterface]
	entries = append(entries, registerMenus()...)
	entries = append(entries, registerFilter()...)
	entries = append(entries, registerForms()...)
	entries = append(entries, registerTable()...)
	entries = append(entries, registerDetail()...)
	entries = append(entries, registerModal()...)
	entries = append(entries, registerDelete()...)
	entries = append(entries, registerSelectionPages()...)
	return pluginPagesWithPatches(entries)
}

func registerMenus() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "appointments.AppointmentDetailMenu", Value: components.SidebarMenu{
			Title: getters.Format("Appointment: %s", getters.Any(getters.Key[string]("appointment.Client.Name"))),
			Back:  appointmentDetailBackItem(),
			Children: []components.PageInterface{
				components.SidebarMenuItem{Title: getters.Static("Appointment Detail"), Url: lamu.RoutePath("appointments.DetailRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("appointment.ID"))})},
			},
		}},
	}
}

func appointmentDetailBackItem() *components.SidebarMenuItem {
	return &components.SidebarMenuItem{
		Title: appointmentDetailBackTitle(),
		Url:   appointmentDetailBackURL(),
	}
}

func appointmentDetailBackTitle() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		clientID, err := getters.Key[uint]("appointment.ClientID")(ctx)
		if err == nil && clientID != 0 {
			return "Back to Client", nil
		}
		return "Back to Appointments Timeline", nil
	}
}

func appointmentDetailBackURL() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		clientID, err := getters.Key[uint]("appointment.ClientID")(ctx)
		if err == nil && clientID != 0 {
			return lamu.RoutePath("clients.DetailRoute", map[string]getters.Getter[any]{
				"id": getters.Any(getters.Static(clientID)),
			})(ctx)
		}
		return lamu.RoutePath("appointments.CardTimelineRoute", nil)(ctx)
	}
}
