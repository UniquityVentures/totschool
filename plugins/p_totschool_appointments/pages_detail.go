package p_totschool_appointments

import (
	"time"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/registry"
)

func registerDetail() []registry.Pair[string, components.PageInterface] {
	generatedSection := []components.PageInterface{
		components.ContainerColumn{Classes: "mt-2 p-4 card card-body border rounded-box border-base-300", Children: []components.PageInterface{
			components.ContainerRow{Classes: "flex flex-wrap justify-between items-center gap-4 mb-4", Children: []components.PageInterface{
				components.FieldTitle{Getter: getters.Static("Generated Letter")},
				components.ContainerColumn{Classes: "flex gap-2", Children: []components.PageInterface{
					components.ButtonLink{Classes: "btn-outline btn-success btn-sm", Label: "Send via WhatsApp", Link: getters.Format("https://wa.me/%v?text=%v", getters.Any(clientPhoneFromIn()), getters.Any(getters.QueryEscape(getters.Key[string]("$in.GeneratedLetter"))))},
					components.ButtonModalForm{Label: "Edit with AI", Name: getters.Static("appointments.AiEditModal"), Url: lamu.RoutePath("appointments.AiEditFormRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("appointment.ID"))}), FormPostURL: lamu.RoutePath("appointments.AiEditRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("appointment.ID"))}), ModalUID: "ai-edit-modal", Classes: "btn-outline btn-secondary btn-sm"},
					components.ButtonPost{Label: "Regenerate Letter", URL: lamu.RoutePath("appointments.GenerateRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("appointment.ID"))}), Classes: "btn-outline btn-primary btn-sm"},
				}},
			}},
			components.FieldMarkdown{Getter: getters.Key[string]("$in.GeneratedLetter")},
		}},
	}

	pendingSection := []components.PageInterface{
		components.HTMXPolling{
			URL: lamu.RoutePath("appointments.DetailRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("appointment.ID"))}),
			Children: []components.PageInterface{
				components.ContainerRow{Classes: "flex gap-2 items-center", Children: []components.PageInterface{
					components.FieldText{Getter: getters.Static("Generating...")},
					components.ButtonPost{
						Label:   "Cancel Generation",
						URL:     lamu.RoutePath("appointments.CancelRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("appointment.ID"))}),
						Classes: "btn-outline btn-error btn-sm",
					},
				}},
			},
		},
	}

	idleSection := []components.PageInterface{
		components.ButtonPost{Label: "Generate Letter with AI", URL: lamu.RoutePath("appointments.GenerateRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("appointment.ID"))}), Classes: "btn-primary"},
	}

	return []registry.Pair[string, components.PageInterface]{
		{Key: "appointments.AppointmentDetail", Value: components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "appointments.AppointmentDetailMenu"}},
			Children: []components.PageInterface{
				components.Detail[Appointment]{
					Getter: getters.Key[Appointment]("appointment"),
					Children: []components.PageInterface{
						components.ContainerColumn{Children: []components.PageInterface{
							components.ShowIf{Getter: getters.Any(getters.Key[bool]("OverlapWarning")), Children: []components.PageInterface{
								components.ContainerColumn{Classes: "bg-warning rounded-box border border-base-300 mb-4 shadow-sm gap-4 p-4", Children: []components.PageInterface{
									components.ContainerRow{Classes: "flex items-center gap-2 font-semibold", Children: []components.PageInterface{
										components.Icon{Name: "exclamation-triangle", Classes: "w-5 h-5"},
										components.FieldText{Getter: getters.Static("Overlapping Appointments:")},
									}},
									components.FieldList[map[string]any]{
										Getter:  getters.Key[[]map[string]any]("OverlapWarningList"),
										Classes: "flex flex-col gap-2 pl-4",
										Children: []components.PageInterface{
											components.ContainerRow{Classes: "items-center gap-2", Children: []components.PageInterface{
												components.ButtonLink{GetterLabel: getters.Key[string]("$row.Name"), Link: lamu.RoutePath("appointments.DetailRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("$row.ID"))})},
												components.FieldText{Getter: getters.Static(" — ")},
												components.FieldDatetime{Getter: getters.Key[time.Time]("$row.Date")},
											}},
										},
									},
								}},
							}},
							components.ContainerRow{Classes: "flex justify-between items-start", Children: []components.PageInterface{
								components.ContainerColumn{Children: []components.PageInterface{
									components.FieldTitle{Getter: getters.Key[string]("$in.Client.Name")},
									components.FieldSubtitle{Getter: getters.Deref(getters.Key[*string]("$in.Client.Address"))},
								}},
							}},
							components.LabelInline{Title: "Phone", Children: []components.PageInterface{components.FieldPhone{Getter: clientPhoneFromIn()}}},
							components.LabelInline{Title: "Status", Children: []components.PageInterface{components.FieldText{Getter: appointmentStatusLabelFromIn()}}},
							components.LabelInline{Title: "Date & Time", Children: []components.PageInterface{components.FieldDatetime{Getter: getters.Key[time.Time]("$in.Datetime")}}},
							components.LabelInline{Title: "Remarks", Children: []components.PageInterface{components.FieldText{Getter: getters.Key[string]("$in.Remarks")}}},
							components.LabelInline{Title: "Extra Info", Children: []components.PageInterface{components.FieldText{Getter: getters.Key[string]("$in.ExtraInfo")}}},
							components.LabelInline{Title: "Created By", Children: []components.PageInterface{components.FieldText{Getter: getters.ForeignKey[p_users.User, uint, string](getters.Key[uint]("$in.CreatedByID"), "Name")}}},
							components.LabelInline{Title: "Created At", Children: []components.PageInterface{components.FieldDatetime{Getter: getters.Key[time.Time]("$in.CreatedAt")}}},
							components.ContainerColumn{Classes: "mt-6", Children: []components.PageInterface{
								components.ShowIf{Getter: getters.Any(getterGenerated()), Children: generatedSection},
								components.ShowIf{Getter: getters.Any(getterGenerationPending()), Children: pendingSection},
								components.ShowIf{Getter: getters.Any(getterIdleGeneration()), Children: idleSection},
							}},
						}},
					},
				},
			},
		}},
	}
}

func registerModal() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "appointments.AiEditModal", Value: components.Modal{
			UID: "ai-edit-modal",
			Children: []components.PageInterface{
				components.FormComponent[Appointment]{
					Getter: getters.Key[Appointment]("appointment"),
					Attr:   getters.FormBubbling(getters.Key[string]("$get.name")),

					Title: "Edit with AI",
					ChildrenInput: []components.PageInterface{
						components.ContainerError{
							Error: getters.Key[error]("$error.generated_letter"),
							Children: []components.PageInterface{
								components.InputTextarea{Name: "generated_letter", Label: "Current Letter Content", Getter: getters.Key[string]("$in.GeneratedLetter"), Rows: 8},
							},
						},
						components.ContainerError{
							Error: getters.Key[error]("$error.instructions"),
							Children: []components.PageInterface{
								components.InputTextarea{Name: "instructions", Label: "Instructions for AI", Rows: 4, Required: true},
							},
						},
					},
					ChildrenAction: []components.PageInterface{
						components.ContainerRow{Classes: "flex justify-end gap-2", Children: []components.PageInterface{
							components.ButtonSubmit{Label: "Generate", Classes: "btn-secondary"},
						}},
					},
				},
			},
		}},
	}
}

func registerDelete() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "appointments.AppointmentDeleteForm", Value: components.Modal{
			UID: "appointment-delete-modal",
			Children: []components.PageInterface{
				components.DeleteConfirmation{
					Title:   "Confirm Deletion",
					Message: "Are you sure you want to delete this appointment?",
					Attr:    getters.FormBubbling(getters.Key[string]("$get.name")),
				},
			},
		}},
	}
}

// --- Selection Tables ---
