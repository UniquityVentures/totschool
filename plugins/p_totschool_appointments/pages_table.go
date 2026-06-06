package p_totschool_appointments

import (
	"time"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/registry"
	"github.com/UniquityVentures/totschool/plugins/p_totschool_clients"
)

func registerFilter() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "appointments.AppointmentFilter", Value: components.FormComponent[Appointment]{
			Attr: getters.FormBoostedGet(lamu.RoutePath("appointments.ListRoute", nil)),

			ChildrenInput: []components.PageInterface{
				components.ContainerError{
					Error: getters.Key[error]("$error.ClientID"),
					Children: []components.PageInterface{
						components.InputForeignKey[p_totschool_clients.Client]{
							Label:       "Client",
							Name:        "ClientID",
							Url:         lamu.RoutePath("clients.SelectRoute", nil),
							Placeholder: "Filter by client…",
							Display:     getters.Key[string]("$in.Name"),
							Getter:      getters.Association[p_totschool_clients.Client](getters.Key[uint]("$get.ClientID")),
						},
					},
				},
				components.ContainerError{
					Error: getters.Key[error]("$error.CreatedByID"),
					Children: []components.PageInterface{
						components.InputForeignKey[p_users.User]{
							Label:       "Created By",
							Name:        "CreatedByID",
							Url:         lamu.RoutePath("appointments.UserSelectRoute", nil),
							Placeholder: "Filter by user…",
							Display:     getters.Key[string]("$in.Name"),
							Getter:      getters.Association[p_users.User](getters.Key[uint]("$get.CreatedByID")),
						},
					},
				},
				components.ContainerError{
					Error: getters.Key[error]("$error.Status"),
					Children: []components.PageInterface{
						components.InputSelect[AppointmentStatus]{
							Label:   "Status",
							Name:    "Status",
							Choices: getters.Static(AppointmentStatusChoices),
							Getter:  appointmentStatusSelectGetter("$get.Status"),
						},
					},
				},
				components.ContainerError{
					Error: getters.Key[error]("$error.Date"),
					Children: []components.PageInterface{
						components.InputDate{Label: "Date", Name: "Date", Getter: getters.Key[time.Time]("$get.Date")},
					},
				},
				components.InputCheckbox{Label: "Overlaps Only", Name: "Overlapping", Getter: getters.Key[bool]("$get.Overlapping")},
			},
			ChildrenAction: []components.PageInterface{
				components.ContainerRow{Classes: "flex gap-2", Children: []components.PageInterface{
					components.ButtonSubmit{Label: "Apply Filters"},
					components.ButtonClear{Label: "Clear"},
				}},
			},
		}},
	}
}

func registerTable() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "appointments.AppointmentTable", Value: components.ShellScaffold{
			Page:    components.Page{Roles: []string{"totschool_admin", "superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "clients.ClientMenu"}},
			Children: []components.PageInterface{
				components.DataTable[Appointment]{
					UID:         "appointment-table",
					Data:        getters.Key[components.ObjectList[Appointment]]("appointments"),
					Title:       "Appointments",
					Subtitle:    "List of appointments",
					DefaultView: "Grid",
					Actions: []components.PageInterface{
						&components.TableButtonFilter{Child: lamu.DynamicPage{Name: "appointments.AppointmentFilter"}},
					},
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("appointments.DetailRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("$row.ID"))})),
					Columns: []components.TableColumn{
						{Label: "Client", Name: "Client", Children: []components.PageInterface{components.FieldText{Getter: getters.ForeignKey[p_totschool_clients.Client, uint, string](getters.Key[uint]("$row.ClientID"), "Name")}}},
						{Label: "Phone", Name: "Phone", Children: []components.PageInterface{components.FieldPhone{Getter: clientPhoneFromRow()}}},
						{Label: "Address", Name: "Address", Children: []components.PageInterface{components.FieldText{Getter: clientAddressFromRow()}}},
						{Label: "Status", Name: "Status", Children: []components.PageInterface{components.FieldText{Getter: appointmentStatusLabelFromRow()}}},
						{Label: "Date & Time", Name: "Datetime", Children: []components.PageInterface{components.FieldDatetime{Getter: getters.Key[time.Time]("$row.Datetime")}}},
						{Label: "Created By", Name: "CreatedBy", Children: []components.PageInterface{components.FieldText{Getter: getters.ForeignKey[p_users.User, uint, string](getters.Key[uint]("$row.CreatedByID"), "Name")}}},
					},
				},
			},
		}},
	}
}

func registerSelectionPages() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "appointments.UserSelectionTable", Value: components.Modal{
			UID: "appointment-user-selection-modal",
			Children: []components.PageInterface{
				components.DataTable[p_users.User]{
					UID:   "appointment-user-selection-table",
					Title: "Select User",
					Data:  getters.Key[components.ObjectList[p_users.User]]("users"),
					RowAttr: getters.RowAttrSelectNamed(
						getters.IfOrElse(getters.Key[string]("$get.target_input"), getters.Static("CreatedByID")),
						getters.Key[uint]("$row.ID"),
						getters.Key[string]("$row.Name"),
					),
					Columns: []components.TableColumn{
						{Label: "Name", Name: "Name", Children: []components.PageInterface{components.FieldText{Getter: getters.Key[string]("$row.Name")}}},
						{Label: "Email", Name: "Email", Children: []components.PageInterface{components.FieldText{Getter: getters.Key[string]("$row.Email")}}},
					},
				},
			},
		}},
		{Key: "appointments.AppointmentSelectionTable", Value: components.Modal{
			UID: "appointment-selection-modal",
			Children: []components.PageInterface{
				components.DataTable[Appointment]{
					UID:     "appointment-selection-table",
					Title:   "Select Appointment",
					Data:    getters.Key[components.ObjectList[Appointment]]("appointments"),
					RowAttr: getters.RowAttrSelect("appointment", getters.Key[uint]("$row.ID"), getters.ForeignKey[p_totschool_clients.Client, uint, string](getters.Key[uint]("$row.ClientID"), "Name")),
					Actions: []components.PageInterface{
						&components.TableButtonFilter{Child: lamu.DynamicPage{Name: "appointments.AppointmentFilter"}},
					},
					Columns: []components.TableColumn{
						{Label: "Client", Name: "Client", Children: []components.PageInterface{components.FieldText{Getter: getters.ForeignKey[p_totschool_clients.Client, uint, string](getters.Key[uint]("$row.ClientID"), "Name")}}},
						{Label: "Location", Name: "Location", Children: []components.PageInterface{components.FieldText{Getter: getters.Key[string]("$row.Location")}}},
						{Label: "Date & Time", Name: "Datetime", Children: []components.PageInterface{components.FieldDatetime{Getter: getters.Key[time.Time]("$row.Datetime")}}},
						{Label: "Status", Name: "Status", Children: []components.PageInterface{components.FieldText{Getter: appointmentStatusLabelFromRow()}}},
					},
				},
			},
		}},
		{Key: "appointments.AppointmentCardTimeline", Value: components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "clients.ClientMenu"}},
			Children: []components.PageInterface{
				components.FormComponent[Appointment]{
					Classes: "max-w-xs mb-4",
					Attr:    appointmentTimelineDateFilterAttr(),
					ChildrenInput: []components.PageInterface{
						components.ContainerError{
							Error: getters.Key[error]("$error.Date"),
							Children: []components.PageInterface{
								components.InputDate{Label: "Date", Name: "Date", Getter: appointmentTimelineDateGetter()},
							},
						},
					},
				},
				components.Timeline[Appointment]{
					UID:     "appointment-timeline",
					Title:   "Appointments Timeline",
					Data:    getters.Key[components.ObjectList[Appointment]]("appointments"),
					OnClick: lamu.RoutePath("appointments.DetailRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("$row.ID"))}),
					Children: []components.PageInterface{
						components.ContainerColumn{
							Children: []components.PageInterface{
								components.FieldText{Classes: "font-bold", Getter: getters.ForeignKey[p_totschool_clients.Client, uint, string](getters.Key[uint]("$row.ClientID"), "Name")},
								components.FieldDatetime{Getter: getters.Key[time.Time]("$row.Datetime"), Classes: "text-sm font-medium whitespace-nowrap"},
								components.ShowIf{Getter: getters.Any(getters.Key[string]("$row.Location")), Children: []components.PageInterface{
									components.FieldText{Getter: getters.Key[string]("$row.Location"), Classes: "text-sm"},
								}},
								components.FieldText{Classes: "text-sm", Getter: clientAddressFromRow()},
								components.FieldPhone{Classes: "text-sm", Getter: clientPhoneFromRow()},
								components.FieldText{Classes: "text-sm", Getter: appointmentStatusLabelFromRow()},
								components.ShowIf{Getter: getters.Any(getters.Key[string]("$row.Remarks")), Children: []components.PageInterface{
									components.FieldText{Getter: getters.Key[string]("$row.Remarks"), Classes: "text-sm italic"},
								}},
							},
						},
					},
				},
			},
		}},
	}
}
