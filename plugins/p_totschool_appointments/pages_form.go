package p_totschool_appointments

import (
	"time"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
	"github.com/UniquityVentures/totschool/plugins/p_totschool_clients"
)

func appointmentFormFields() []components.PageInterface {
	return []components.PageInterface{
		components.ContainerError{
			Error: getters.Key[error]("$error.ClientID"),
			Children: []components.PageInterface{
				components.InputForeignKey[p_totschool_clients.Client]{
					Name:        "ClientID",
					Label:       "Client",
					Url:         lamu.RoutePath("clients.SelectRoute", nil),
					Display:     getters.Key[string]("$in.Name"),
					Placeholder: "Select a client...",
					Required:    true,
					Getter:      getters.Association[p_totschool_clients.Client](getters.Key[uint]("$in.ClientID")),
				},
			},
		},
		components.ContainerError{
			Error: getters.Key[error]("$error.Datetime"),
			Children: []components.PageInterface{
				components.InputDatetime{Label: "Date & Time", Name: "Datetime", Required: true, Getter: getters.Key[time.Time]("$in.Datetime")},
			},
		},
		components.ContainerError{
			Error: getters.Key[error]("$error.Status"),
			Children: []components.PageInterface{
				components.InputSelect[AppointmentStatus]{
					Name:     "Status",
					Label:    "Status",
					Required: true,
					Choices:  getters.Static(AppointmentStatusChoices),
					Getter:   appointmentStatusSelectGetter("$in.Status"),
				},
			},
		},
		components.ContainerError{
			Error: getters.Key[error]("$error.Location"),
			Children: []components.PageInterface{
				components.InputText{Label: "Location", Name: "Location", Getter: getters.Key[string]("$in.Location")},
			},
		},
		components.ContainerError{
			Error: getters.Key[error]("$error.Remarks"),
			Children: []components.PageInterface{
				components.InputTextarea{Label: "Remarks", Name: "Remarks", Getter: getters.Key[string]("$in.Remarks"), Rows: 2},
			},
		},
		components.ContainerError{
			Error: getters.Key[error]("$error.ExtraInfo"),
			Children: []components.PageInterface{
				components.InputTextarea{Label: "Extra Info (For AI)", Name: "ExtraInfo", Getter: getters.Key[string]("$in.ExtraInfo"), Rows: 2},
			},
		},
	}
}

func registerForms() []registry.Pair[string, components.PageInterface] {
	createFormName := getters.Static("appointments.AppointmentCreateForm")
	updateFormName := getters.Static("appointments.AppointmentUpdateForm")
	deleteFormName := getters.Static("appointments.AppointmentDeleteForm")
	return []registry.Pair[string, components.PageInterface]{
		{Key: "appointments.AppointmentCreateForm", Value: components.Modal{
			UID: "appointment-create-modal",
			Children: []components.PageInterface{
				components.FormComponent[Appointment]{
					Attr:           getters.FormBubbling(createFormName),
					Title:          "Create Appointment",
					Subtitle:       "Create a new appointment",
					ChildrenInput:  appointmentFormFields(),
					ChildrenAction: []components.PageInterface{components.ButtonSubmit{Label: "Save Appointment"}},
				},
			},
		}},
		{Key: "appointments.AppointmentUpdateForm", Value: components.Modal{
			UID: "appointment-update-modal",
			Children: []components.PageInterface{
				components.FormComponent[Appointment]{
					Getter:        getters.Key[Appointment]("appointment"),
					Attr:          getters.FormBubbling(updateFormName),
					Title:         "Edit Appointment",
					Subtitle:      "Update appointment details",
					ChildrenInput: appointmentFormFields(),
					ChildrenAction: []components.PageInterface{
						components.ContainerRow{
							Classes: "flex flex-wrap justify-between gap-2 mt-2 items-center",
							Children: []components.PageInterface{
								components.ContainerRow{
									Classes: "flex justify-end gap-2",
									Children: []components.PageInterface{
										components.ButtonSubmit{Label: "Save Appointment"},
										components.ButtonModalForm{
											Label:       "Delete",
											Icon:        "trash",
											Name:        deleteFormName,
											Url:         lamu.RoutePath("appointments.DeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("appointment.ID"))}),
											FormPostURL: lamu.RoutePath("appointments.DeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("appointment.ID"))}),
											ModalUID:    "appointment-delete-modal",
											Classes:     "btn-error",
										},
									},
								},
							},
						},
					},
				},
			},
		}},
	}
}
