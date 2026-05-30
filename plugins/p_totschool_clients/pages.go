package p_totschool_clients

import (
	"context"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/registry"
)

var clientAdminRoles = []string{"totschool_admin", "superuser"}

func clientFormFields() []components.PageInterface {
	return []components.PageInterface{
		components.ContainerError{
			Error: getters.Key[error]("$error.Name"),
			Children: []components.PageInterface{
				components.InputText{Label: "Name", Name: "Name", Required: true, Getter: getters.Key[string]("$in.Name")},
			},
		},
		components.ContainerError{
			Error: getters.Key[error]("$error.Address"),
			Children: []components.PageInterface{
				components.InputTextarea{Label: "Address", Name: "Address", Getter: getters.Deref(getters.Key[*string]("$in.Address")), Rows: 2},
			},
		},
		components.ContainerError{
			Error: getters.Key[error]("$error.Phone"),
			Children: []components.PageInterface{
				components.InputPhone{Label: "Phone", Name: "Phone", Getter: getters.Deref(getters.Key[*string]("$in.Phone"))},
			},
		},
		components.ContainerError{
			Error: getters.Key[error]("$error.Remarks"),
			Children: []components.PageInterface{
				components.InputTextarea{Label: "Remarks", Name: "Remarks", Getter: getters.Deref(getters.Key[*string]("$in.Remarks")), Rows: 3},
			},
		},
	}
}

func registerMenus() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "clients.ClientMenu", Value: components.SidebarMenu{
			Title: getters.Static("Clients"),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("Back to All Apps"),
				Url:   lamu.RoutePath("dashboard.AppsPage", nil),
			},
			Children: []components.PageInterface{
				components.SidebarMenuItem{Title: getters.Static("All Clients"), Url: lamu.RoutePath("clients.ListRoute", nil)},
				components.SidebarMenuItem{Title: getters.Static("Appointments Timeline"), Url: lamu.RoutePath("appointments.CardTimelineRoute", nil)},
				components.SidebarMenuItem{Title: getters.Static("Old Proposals"), Url: lamu.RoutePath("proposals.ListRoute", nil)},
				components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"totschool_admin", "superuser"}},
					Title: getters.Static("All Appointments"),
					Url:   lamu.RoutePath("appointments.ListRoute", nil),
				},
			},
		}},
		{Key: "clients.ClientDetailMenu", Value: components.SidebarMenu{
			Title: getters.Format("Client: %s", getters.Any(getters.Key[string]("client.Name"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("Back to all Clients"),
				Url:   lamu.RoutePath("clients.ListRoute", nil),
			},
			Children: []components.PageInterface{
				components.SidebarMenuItem{Title: getters.Static("Client Detail"), Url: lamu.RoutePath("clients.DetailRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("client.ID"))})},
				components.SidebarMenuItem{Title: getters.Static("Edit Client"), Url: lamu.RoutePath("clients.UpdateRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("client.ID"))})},
			},
		}},
	}
}

func registerFilter() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "clients.ClientFilter", Value: components.FormComponent[Client]{
			Attr: getters.FormBoostedGet(lamu.RoutePath("clients.ListRoute", nil)),
			ChildrenInput: []components.PageInterface{
				components.InputText{Label: "Name", Name: "Name", Getter: getters.Key[string]("$get.Name")},
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
	createFormName := getters.Static("clients.ClientCreateForm")
	updateFormName := getters.Static("clients.ClientUpdateForm")
	deleteFormName := getters.Static("clients.ClientDeleteForm")
	return []registry.Pair[string, components.PageInterface]{
		{Key: "clients.ClientCreateForm", Value: components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "clients.ClientMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      createFormName,
					ActionURL: lamu.RoutePath("clients.CreateRoute", nil),
					Children: []components.PageInterface{
						components.FormComponent[Client]{
							Attr:           getters.FormBubbling(createFormName),
							Title:          "Create Client",
							Subtitle:       "Add a new client",
							ChildrenInput:  clientFormFields(),
							ChildrenAction: []components.PageInterface{components.ButtonSubmit{Label: "Save Client"}},
						},
					},
				},
			},
		}},
		{Key: "clients.ClientUpdateForm", Value: components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "clients.ClientDetailMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      updateFormName,
					ActionURL: lamu.RoutePath("clients.UpdateRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("client.ID"))}),
					Children: []components.PageInterface{
						components.FormComponent[Client]{
							Getter:        getters.Key[Client]("client"),
							Attr:          getters.FormBubbling(updateFormName),
							Title:         "Edit Client",
							Subtitle:      "Update client details",
							ChildrenInput: clientFormFields(),
							ChildrenAction: []components.PageInterface{
								components.ContainerRow{
									Classes: "flex flex-wrap justify-between gap-2 mt-2 items-center",
									Children: []components.PageInterface{
										components.ContainerRow{
											Classes: "flex justify-end gap-2",
											Children: []components.PageInterface{
												components.ButtonSubmit{Label: "Save Client"},
												components.ButtonModalForm{
													Label:       "Delete",
													Icon:        "trash",
													Name:        deleteFormName,
													Url:         lamu.RoutePath("clients.DeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("client.ID"))}),
													FormPostURL: lamu.RoutePath("clients.DeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("client.ID"))}),
													ModalUID:    "client-delete-modal",
													Classes:     "btn-error",
												},
											},
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

func registerTable() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "clients.ClientTable", Value: components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "clients.ClientMenu"}},
			Children: []components.PageInterface{
				components.DataTable[Client]{
					UID:      "client-table",
					Data:     getters.Key[components.ObjectList[Client]]("clients"),
					Title:    "Clients",
					Subtitle: "List of clients",
					Actions: []components.PageInterface{
						&components.TableButtonFilter{Child: lamu.DynamicPage{Name: "clients.ClientFilter"}},
						&components.TableButtonCreate{Link: lamu.RoutePath("clients.CreateRoute", nil)},
					},
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("clients.DetailRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("$row.ID"))})),
					EnabledColumns: getters.Map(getters.Key[string]("$role"), func(_ context.Context, role string) (map[string]bool, error) {
						if role == "totschool_admin" || role == "superuser" {
							return nil, nil
						}
						return map[string]bool{"Name": true, "Address": true, "Phone": true}, nil
					}),
					Columns: []components.TableColumn{
						{Label: "Name", Name: "Name", Children: []components.PageInterface{components.FieldText{Getter: getters.Key[string]("$row.Name")}}},
						{Label: "Address", Name: "Address", Children: []components.PageInterface{components.FieldText{Getter: getters.Deref(getters.Key[*string]("$row.Address"))}}},
						{Label: "Phone", Name: "Phone", Children: []components.PageInterface{components.FieldPhone{Getter: getters.Deref(getters.Key[*string]("$row.Phone"))}}},
						{Label: "Created By", Name: "CreatedBy", Children: []components.PageInterface{components.FieldText{Getter: getters.ForeignKey[p_users.User, uint, string](getters.Key[uint]("$row.CreatedByID"), "Name")}}},
					},
				},
			},
		}},
	}
}

func registerDetail() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "clients.ClientDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "clients.ClientDetailMenu"}},
			Children: []components.PageInterface{
				&components.Detail[Client]{
					Getter: getters.Key[Client]("client"),
					Children: []components.PageInterface{
						components.ContainerColumn{
							Page: components.Page{Key: "clients.ClientDetailContent"},
							Children: []components.PageInterface{
								components.FieldTitle{Getter: getters.Key[string]("$in.Name")},
								components.LabelInline{Title: "Address", Children: []components.PageInterface{components.FieldText{Getter: getters.Deref(getters.Key[*string]("$in.Address"))}}},
								components.LabelInline{Title: "Phone", Children: []components.PageInterface{components.FieldPhone{Getter: getters.Deref(getters.Key[*string]("$in.Phone"))}}},
								components.LabelInline{Title: "Remarks", Children: []components.PageInterface{components.FieldText{Getter: getters.Deref(getters.Key[*string]("$in.Remarks"))}}},
								components.LabelInline{
									Page:    components.Page{Roles: clientAdminRoles},
									Title:   "Created By",
									Children: []components.PageInterface{components.FieldText{Getter: getters.ForeignKey[p_users.User, uint, string](getters.Key[uint]("$in.CreatedByID"), "Name")}}},
							},
						},
					},
				},
			},
		}},
	}
}

func registerDelete() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "clients.ClientDeleteForm", Value: components.Modal{
			UID: "client-delete-modal",
			Children: []components.PageInterface{
				components.DeleteConfirmation{
					Title:   "Confirm Deletion",
					Message: "Are you sure you want to delete this client?",
					Attr:    getters.FormBubbling(getters.Key[string]("$get.name")),
				},
			},
		}},
	}
}

func registerSelection() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "clients.ClientSelectionTable", Value: components.Modal{
			UID: "client-selection-modal",
			Children: []components.PageInterface{
				components.DataTable[Client]{
					UID:   "client-selection-table",
					Title: "Select Client",
					Data:  getters.Key[components.ObjectList[Client]]("clients"),
					RowAttr: getters.RowAttrSelectNamed(
						getters.IfOrElse(getters.Key[string]("$get.target_input"), getters.Static("ClientID")),
						getters.Key[uint]("$row.ID"),
						getters.Key[string]("$row.Name"),
					),
					Columns: []components.TableColumn{
						{Label: "Name", Name: "Name", Children: []components.PageInterface{components.FieldText{Getter: getters.Key[string]("$row.Name")}}},
						{Label: "Phone", Name: "Phone", Children: []components.PageInterface{components.FieldPhone{Getter: getters.Deref(getters.Key[*string]("$row.Phone"))}}},
						{Label: "Address", Name: "Address", Children: []components.PageInterface{components.FieldText{Getter: getters.Deref(getters.Key[*string]("$row.Address"))}}},
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
	entries = append(entries, registerDelete()...)
	entries = append(entries, registerSelection()...)
	return lamu.PluginFeatures[components.PageInterface]{Entries: entries}
}
