package p_totschool_followups

import (
	"context"
	"time"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/registry"
	"github.com/UniquityVentures/totschool/plugins/p_totschool_clients"
)

func registerMenus() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "followups.FollowupDetailMenu", Value: components.SidebarMenu{
			Title: getters.Format("Follow-up: %s", getters.Any(getters.Key[string]("followup.Title"))),
			Back:  followupDetailBackItem(),
			Children: []components.PageInterface{
				components.SidebarMenuItem{Title: getters.Static("Follow-up Detail"), Url: lamu.RoutePath("followups.DetailRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("followup.ID"))})},
			},
		}},
	}
}

func followupDetailBackItem() *components.SidebarMenuItem {
	return &components.SidebarMenuItem{
		Title: followupDetailBackTitle(),
		Url:   followupDetailBackURL(),
	}
}

func followupDetailBackTitle() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		clientID, err := getters.Key[uint]("followup.ClientID")(ctx)
		if err == nil && clientID != 0 {
			return "Back to Client", nil
		}
		return "Back to Follow-ups", nil
	}
}

func followupDetailBackURL() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		clientID, err := getters.Key[uint]("followup.ClientID")(ctx)
		if err == nil && clientID != 0 {
			return lamu.RoutePath("clients.DetailRoute", map[string]getters.Getter[any]{
				"id": getters.Any(getters.Static(clientID)),
			})(ctx)
		}
		return lamu.RoutePath("followups.ListRoute", nil)(ctx)
	}
}

func followupFormFields(includeClientPicker bool) []components.PageInterface {
	fields := []components.PageInterface{}
	if includeClientPicker {
		fields = append(fields, components.ContainerError{
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
		})
	} else {
		fields = append(fields, components.InputForeignKey[p_totschool_clients.Client]{
			Hidden: true,
			Name:   "ClientID",
			Getter: getters.Association[p_totschool_clients.Client](getters.Key[uint]("$in.ClientID")),
		})
	}
	fields = append(fields,
		components.ContainerError{
			Error: getters.Key[error]("$error.Title"),
			Children: []components.PageInterface{
				components.InputText{Label: "Follow-up Title", Name: "Title", Required: true, Getter: getters.Key[string]("$in.Title")},
			},
		},
		components.ContainerError{
			Error: getters.Key[error]("$error.ExtraInfo"),
			Children: []components.PageInterface{
				components.InputTextarea{Label: "Extra Info (For AI)", Name: "ExtraInfo", Getter: getters.Key[string]("$in.ExtraInfo"), Rows: 3},
			},
		},
	)
	return fields
}

func registerFilter() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "followups.FollowupFilter", Value: components.FormComponent[Followup]{
			Attr: getters.FormBoostedGet(lamu.RoutePath("followups.ListRoute", nil)),
			ChildrenInput: []components.PageInterface{
				components.ContainerError{
					Error: getters.Key[error]("$error.ClientID"),
					Children: []components.PageInterface{
						components.InputForeignKey[p_totschool_clients.Client]{
							Label:       "Client",
							Name:        "ClientID",
							Url:         lamu.RoutePath("clients.SelectRoute", nil),
							Placeholder: "Filter by client...",
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
							Url:         lamu.RoutePath("followups.UserSelectRoute", nil),
							Placeholder: "Filter by user...",
							Display:     getters.Key[string]("$in.Name"),
							Getter:      getters.Association[p_users.User](getters.Key[uint]("$get.CreatedByID")),
						},
					},
				},
				components.InputText{Label: "Title", Name: "Title", Getter: getters.Key[string]("$get.Title")},
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

func registerForms() []registry.Pair[string, components.PageInterface] {
	createFormName := getters.Static("followups.FollowupCreateForm")
	updateFormName := getters.Static("followups.FollowupUpdateForm")
	deleteFormName := getters.Static("followups.FollowupDeleteForm")
	return []registry.Pair[string, components.PageInterface]{
		{Key: "followups.FollowupCreateForm", Value: components.Modal{
			UID: "followup-create-modal",
			Children: []components.PageInterface{
				components.FormComponent[Followup]{
					Attr:           getters.FormBubbling(createFormName),
					Title:          "Create Follow-up Letter",
					Subtitle:       "Create a follow-up letter linked to a client proposal",
					ChildrenInput:  followupFormFields(false),
					ChildrenAction: []components.PageInterface{components.ButtonSubmit{Label: "Save Follow-up"}},
				},
			},
		}},
		{Key: "followups.FollowupUpdateForm", Value: components.Modal{
			UID: "followup-update-modal",
			Children: []components.PageInterface{
				components.FormComponent[Followup]{
					Getter:        getters.Key[Followup]("followup"),
					Attr:          getters.FormBubbling(updateFormName),
					Title:         "Edit Follow-up Letter",
					Subtitle:      "Update follow-up details",
					ChildrenInput: followupFormFields(false),
					ChildrenAction: []components.PageInterface{
						components.ContainerRow{
							Classes: "flex flex-wrap justify-end gap-2 mt-2",
							Children: []components.PageInterface{
								components.ButtonSubmit{Label: "Save Follow-up"},
								components.ButtonModalForm{
									Label:       "Delete",
									Icon:        "trash",
									Name:        deleteFormName,
									Url:         lamu.RoutePath("followups.DeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("followup.ID"))}),
									FormPostURL: lamu.RoutePath("followups.DeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("followup.ID"))}),
									ModalUID:    "followup-delete-modal",
									Classes:     "btn-error",
								},
							},
						},
					},
				},
			},
		}},
	}
}

func registerTable() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "followups.FollowupTable", Value: components.ShellScaffold{
			Page:    components.Page{Roles: []string{"totschool_admin", "superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "clients.ClientMenu"}},
			Children: []components.PageInterface{
				components.DataTable[Followup]{
					UID:      "followup-table",
					Data:     getters.Key[components.ObjectList[Followup]]("followups"),
					Title:    "Follow-up Letters",
					Subtitle: "List of generated follow-up letters",
					Actions: []components.PageInterface{
						&components.TableButtonFilter{Child: lamu.DynamicPage{Name: "followups.FollowupFilter"}},
					},
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("followups.DetailRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("$row.ID"))})),
					Columns: []components.TableColumn{
						{Label: "Title", Name: "Title", Children: []components.PageInterface{components.FieldText{Getter: getters.Key[string]("$row.Title")}}},
						{Label: "Client", Name: "Client", Children: []components.PageInterface{components.FieldText{Getter: getters.ForeignKey[p_totschool_clients.Client, uint, string](getters.Key[uint]("$row.ClientID"), "Name")}}},
						{Label: "Phone", Name: "Phone", Children: []components.PageInterface{components.FieldPhone{Getter: getters.Deref(getters.Key[*string]("$row.Client.Phone"))}}},
						{Label: "Created By", Name: "CreatedBy", Children: []components.PageInterface{components.FieldText{Getter: getters.ForeignKey[p_users.User, uint, string](getters.Key[uint]("$row.CreatedByID"), "Name")}}},
						{Label: "Created At", Name: "CreatedAt", Children: []components.PageInterface{components.FieldDatetime{Getter: getters.Key[time.Time]("$row.CreatedAt")}}},
					},
				},
			},
		}},
	}
}

func followupDetailEditButton() components.PageInterface {
	updateFormName := getters.Static("followups.FollowupUpdateForm")
	return components.ButtonModalForm{
		Label: "Edit Follow-up",
		Icon:  "pencil",
		Name:  updateFormName,
		Url: lamu.RoutePath("followups.UpdateRoute", map[string]getters.Getter[any]{
			"id": getters.Any(getters.Key[uint]("followup.ID")),
		}),
		FormPostURL: lamu.RoutePath("followups.UpdateRoute", map[string]getters.Getter[any]{
			"id": getters.Any(getters.Key[uint]("followup.ID")),
		}),
		ModalUID: "followup-update-modal",
		Classes:  "btn-outline",
	}
}

func registerDetail() []registry.Pair[string, components.PageInterface] {
	generatedSection := []components.PageInterface{
		components.ContainerColumn{Classes: "mt-2 p-4 card card-body border rounded-box border-base-300", Children: []components.PageInterface{
			components.ContainerRow{Classes: "flex flex-wrap justify-between items-center gap-4 mb-4", Children: []components.PageInterface{
				components.FieldTitle{Getter: getters.Static("Generated Follow-up Letter")},
				components.ContainerColumn{Classes: "flex flex-wrap gap-2", Children: []components.PageInterface{
					components.ButtonLink{Classes: "btn-outline btn-success btn-sm", Label: "Send via WhatsApp", Link: getters.Format("https://wa.me/%v?text=%v", getters.Any(getters.Deref(getters.Key[*string]("$in.Client.Phone"))), getters.Any(getters.QueryEscape(getters.Key[string]("$in.GeneratedLetter"))))},
					components.ButtonDownload{Label: "Export to PDF", Link: lamu.RoutePath("followups.ExportPdfRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("followup.ID"))}), Classes: "btn-outline btn-secondary btn-sm"},
					components.ButtonDownload{Label: "Export to Word", Link: lamu.RoutePath("followups.ExportDocxRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("followup.ID"))}), Classes: "btn-outline btn-secondary btn-sm"},
					components.ButtonModalForm{Label: "Edit with AI", Name: getters.Static("followups.AiEditModal"), Url: lamu.RoutePath("followups.AiEditFormRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("followup.ID"))}), FormPostURL: lamu.RoutePath("followups.AiEditRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("followup.ID"))}), ModalUID: "followup-ai-edit-modal", Classes: "btn-outline btn-secondary btn-sm"},
					components.ButtonPost{Label: "Regenerate Letter", URL: lamu.RoutePath("followups.GenerateRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("followup.ID"))}), Classes: "btn-outline btn-primary btn-sm"},
				}},
			}},
			components.FieldMarkdown{Getter: getters.Key[string]("$in.GeneratedLetter")},
		}},
	}

	pendingSection := []components.PageInterface{
		components.HTMXPolling{
			URL: lamu.RoutePath("followups.DetailRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("followup.ID"))}),
			Children: []components.PageInterface{
				components.ContainerRow{Classes: "flex gap-2 items-center", Children: []components.PageInterface{
					components.FieldText{Getter: getters.Static("Generating...")},
					components.ButtonPost{
						Label:   "Cancel Generation",
						URL:     lamu.RoutePath("followups.CancelRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("followup.ID"))}),
						Classes: "btn-outline btn-error btn-sm",
					},
				}},
			},
		},
	}

	idleSection := []components.PageInterface{
		components.ButtonPost{Label: "Generate Letter with AI", URL: lamu.RoutePath("followups.GenerateRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("followup.ID"))}), Classes: "btn-primary"},
	}

	return []registry.Pair[string, components.PageInterface]{
		{Key: "followups.FollowupDetail", Value: components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "followups.FollowupDetailMenu"}},
			Children: []components.PageInterface{
				components.Detail[Followup]{
					Getter: getters.Key[Followup]("followup"),
					Children: []components.PageInterface{
						components.ContainerColumn{Children: []components.PageInterface{
							components.FieldTitle{Getter: getters.Key[string]("$in.Title")},
							components.LabelInline{Title: "Client", Children: []components.PageInterface{components.FieldText{Getter: getters.Key[string]("$in.Client.Name")}}},
							components.LabelInline{Title: "Phone", Children: []components.PageInterface{components.FieldPhone{Getter: getters.Deref(getters.Key[*string]("$in.Client.Phone"))}}},
							components.LabelInline{Title: "Address", Children: []components.PageInterface{components.FieldText{Getter: getters.Deref(getters.Key[*string]("$in.Client.Address"))}}},
							components.LabelInline{Title: "Extra Info", Children: []components.PageInterface{components.FieldText{Getter: getters.Key[string]("$in.ExtraInfo")}}},
							components.LabelInline{Title: "Created By", Children: []components.PageInterface{components.FieldText{Getter: getters.ForeignKey[p_users.User, uint, string](getters.Key[uint]("$in.CreatedByID"), "Name")}}},
							components.LabelInline{Title: "Created At", Children: []components.PageInterface{components.FieldDatetime{Getter: getters.Key[time.Time]("$in.CreatedAt")}}},
							components.ContainerRow{Classes: "flex flex-wrap gap-2 my-2", Children: []components.PageInterface{
								components.ShowIf{Getter: getters.Any(getterIdleGeneration()), Children: idleSection},
								followupDetailEditButton(),
							}},
							components.ContainerColumn{Classes: "mt-6", Children: []components.PageInterface{
								components.ShowIf{Getter: getters.Any(getterGenerated()), Children: generatedSection},
								components.ShowIf{Getter: getters.Any(getterGenerationPending()), Children: pendingSection},
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
		{Key: "followups.AiEditModal", Value: components.Modal{
			UID: "followup-ai-edit-modal",
			Children: []components.PageInterface{
				components.FormComponent[Followup]{
					Getter: getters.Key[Followup]("followup"),
					Attr:   getters.FormBubbling(getters.Key[string]("$get.name")),
					Title:  "Edit with AI",
					ChildrenInput: []components.PageInterface{
						components.ContainerError{
							Error: getters.Key[error]("$error.GeneratedLetter"),
							Children: []components.PageInterface{
								components.InputTextarea{Name: "GeneratedLetter", Label: "Current Letter Markdown", Getter: getters.Key[string]("$in.GeneratedLetter"), Rows: 8},
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
		{Key: "followups.FollowupDeleteForm", Value: components.Modal{
			UID: "followup-delete-modal",
			Children: []components.PageInterface{
				components.DeleteConfirmation{
					Title:   "Confirm Deletion",
					Message: "Are you sure you want to delete this follow-up letter?",
					Attr:    getters.FormBubbling(getters.Key[string]("$get.name")),
				},
			},
		}},
	}
}

func registerSelectionPages() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "followups.UserSelectionTable", Value: components.Modal{
			UID: "followup-user-selection-modal",
			Children: []components.PageInterface{
				components.DataTable[p_users.User]{
					UID:   "followup-user-selection-table",
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
	}
}

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
