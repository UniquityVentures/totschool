package p_totschool_appointments

import (
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
		{Key: "appointments.AppointmentMenu", Value: components.SidebarMenu{
			Title: getters.Static("Appointments"),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("Back to All Apps"),
				Url:   lamu.RoutePath("dashboard.AppsPage", nil),
			},
			Children: []components.PageInterface{
				components.SidebarMenuItem{Title: getters.Static("All Appointments"), Url: lamu.RoutePath("appointments.ListRoute", nil)},
				components.SidebarMenuItem{Title: getters.Static("Appointments Timeline"), Url: lamu.RoutePath("appointments.CardTimelineRoute", nil)},
			},
		}},
		{Key: "appointments.AppointmentDetailMenu", Value: components.SidebarMenu{
			Title: getters.Format("Appointment: %s", getters.Any(getters.Key[string]("appointment.Client.Name"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("Back to all Appointments"),
				Url:   lamu.RoutePath("appointments.ListRoute", nil),
			},
			Children: []components.PageInterface{
				components.SidebarMenuItem{Title: getters.Static("Appointment Detail"), Url: lamu.RoutePath("appointments.DetailRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("appointment.ID"))})},
			},
		}},
	}
}
