package p_totschool_proposals

import (
	"context"
	"time"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
	"github.com/UniquityVentures/totschool/plugins/p_totschool_clients"
	"gorm.io/datatypes"
)

func registerMenus() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "proposals.ProposalDetailMenu", Value: components.SidebarMenu{
			Title: getters.Format("Proposal: %s", getters.Any(getters.Key[string]("proposal.Title"))),
			Back:  proposalDetailBackItem(),
			Children: []components.PageInterface{
				components.SidebarMenuItem{Title: getters.Static("Proposal Detail"), Url: lamu.RoutePath("proposals.DetailRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("proposal.ID"))})},
			},
		}},
	}
}

func proposalDetailBackItem() *components.SidebarMenuItem {
	return &components.SidebarMenuItem{
		Title: proposalDetailBackTitle(),
		Url:   proposalDetailBackURL(),
	}
}

func proposalDetailBackTitle() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		clientID, err := getters.Deref(getters.Key[*uint]("proposal.ClientID"))(ctx)
		if err == nil && clientID != 0 {
			return "Back to Client", nil
		}
		return "Back to Unassigned Proposals", nil
	}
}

func proposalDetailBackURL() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		clientID, err := getters.Deref(getters.Key[*uint]("proposal.ClientID"))(ctx)
		if err == nil && clientID != 0 {
			return lamu.RoutePath("clients.DetailRoute", map[string]getters.Getter[any]{
				"id": getters.Any(getters.Static(clientID)),
			})(ctx)
		}
		return lamu.RoutePath("proposals.ListRoute", nil)(ctx)
	}
}

func registerFilter() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{{Key: "proposals.ProposalFilter", Value: components.FormComponent[Proposal]{
		Attr: getters.FormBoostedGet(lamu.RoutePath("proposals.ListRoute", nil)),

		ChildrenInput: []components.PageInterface{
			components.InputText{Label: "Title", Name: "Title", Getter: getters.Key[string]("$get.Title")},
		},
		ChildrenAction: []components.PageInterface{
			components.ContainerRow{Classes: "flex gap-2", Children: []components.PageInterface{
				components.ButtonSubmit{Label: "Apply Filters"},
				components.ButtonClear{Label: "Clear"},
			}},
		},
	}}}
}

func proposalClientPickerField() components.PageInterface {
	return components.ContainerError{
		Error: getters.Key[error]("$error.ClientID"),
		Children: []components.PageInterface{
			components.InputForeignKey[p_totschool_clients.Client]{
				Name:        "ClientID",
				Label:       "Client",
				Url:         getters.Format("%s?without_proposal=1", getters.Any(lamu.RoutePath("clients.SelectRoute", nil))),
				Display:     getters.Key[string]("$in.Name"),
				Placeholder: "Select a client...",
				Getter:      getters.Association[p_totschool_clients.Client](proposalFormClientID()),
			},
		},
	}
}

func proposalClientHiddenField() components.PageInterface {
	return components.InputForeignKey[p_totschool_clients.Client]{
		Hidden: true,
		Name:   "ClientID",
		Getter: getters.Association[p_totschool_clients.Client](proposalFormClientID()),
	}
}

func proposalClientAssignmentFields() []components.PageInterface {
	return []components.PageInterface{
		components.ShowIf{
			Getter:   getters.Any(getterProposalUnassigned()),
			Children: []components.PageInterface{proposalClientPickerField()},
		},
		components.ShowIf{
			Getter:   getters.Any(getterProposalAssigned()),
			Children: []components.PageInterface{proposalClientHiddenField()},
		},
	}
}

func proposalQuestionnaireCoreFields() []components.PageInterface {
	return []components.PageInterface{
		components.InputText{Label: "Proposal Title", Name: "Title", Required: true, Getter: getters.Key[string]("$in.Title")},
		components.InputKeyValue{Getter: getters.Key[datatypes.JSON]("$in.Answers"), Keys: getters.Static(QUESTIONS), Name: "Answers"},
	}
}

func proposalQuestionnaireFields(includeClientPicker bool) []components.PageInterface {
	fields := []components.PageInterface{}
	if includeClientPicker {
		fields = append(fields, proposalClientPickerField())
	} else {
		fields = append(fields, proposalClientHiddenField())
	}
	fields = append(fields, proposalQuestionnaireCoreFields()...)
	return fields
}

func proposalUpdateFormFields() []components.PageInterface {
	fields := proposalClientAssignmentFields()
	fields = append(fields, proposalQuestionnaireCoreFields()...)
	return fields
}

func registerForms() []registry.Pair[string, components.PageInterface] {
	createFormName := getters.Static("proposals.ProposalCreateForm")
	updateFormName := getters.Static("proposals.ProposalUpdateForm")
	deleteFormName := getters.Static("proposals.ProposalDeleteForm")
	return []registry.Pair[string, components.PageInterface]{
		{Key: "proposals.ProposalCreateForm", Value: components.Modal{
			UID: "proposal-create-modal",
			Children: []components.PageInterface{
				components.FormComponent[Proposal]{
					Attr:           getters.FormBubbling(createFormName),
					Title:          "Create Proposal",
					Subtitle:       "Fill in the questionnaire answers",
					ChildrenInput:  proposalQuestionnaireFields(false),
					ChildrenAction: []components.PageInterface{components.ButtonSubmit{Label: "Save Proposal"}},
				},
			},
		}},
		{Key: "proposals.ProposalUpdateForm", Value: components.Modal{
			UID: "proposal-update-modal",
			Children: []components.PageInterface{
				components.FormComponent[Proposal]{
					Getter:        getters.Key[Proposal]("proposal"),
					Attr:          getters.FormBubbling(updateFormName),
					Title:         "Edit Proposal",
					Subtitle:      "Update questionnaire answers",
					ChildrenInput: proposalUpdateFormFields(),
					ChildrenAction: []components.PageInterface{
						components.ContainerRow{
							Classes: "flex flex-wrap justify-end gap-2 mt-2",
							Children: []components.PageInterface{
								components.ButtonSubmit{Label: "Save Proposal"},
								components.ButtonModalForm{
									Label:       "Delete",
									Icon:        "trash",
									Name:        deleteFormName,
									Url:         lamu.RoutePath("proposals.DeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("proposal.ID"))}),
									FormPostURL: lamu.RoutePath("proposals.DeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("proposal.ID"))}),
									ModalUID:    "proposal-delete-modal",
									Classes:     "btn-error",
								},
							},
						},
					},
				},
			},
		}},
		{Key: "proposals.ProposalUpdatePageForm", Value: components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "proposals.ProposalDetailMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      updateFormName,
					ActionURL: lamu.RoutePath("proposals.UpdateRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("proposal.ID"))}),
					Children: []components.PageInterface{
						components.FormComponent[Proposal]{
							Getter: getters.Key[Proposal]("proposal"),
							Attr:   getters.FormBubbling(updateFormName),

							Title:         "Edit Proposal",
							Subtitle:      "Update questionnaire answers",
							ChildrenInput: proposalUpdateFormFields(),
							ChildrenAction: []components.PageInterface{
								components.ContainerRow{
									Classes: "flex flex-wrap justify-between gap-2 mt-2 items-center",
									Children: []components.PageInterface{
										components.ContainerRow{
											Classes: "flex justify-end gap-2",
											Children: []components.PageInterface{
												components.ButtonSubmit{Label: "Save Proposal"},
												components.ButtonModalForm{
													Label:       "Delete",
													Icon:        "trash",
													Name:        deleteFormName,
													Url:         lamu.RoutePath("proposals.DeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("proposal.ID"))}),
													FormPostURL: lamu.RoutePath("proposals.DeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("proposal.ID"))}),
													ModalUID:    "proposal-delete-modal",
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
		{
			Key: "proposals.ProposalTable", Value: components.ShellScaffold{
				Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "clients.ClientMenu"}},
				Children: []components.PageInterface{
					components.DataTable[Proposal]{
						UID:      "proposal-table",
						Data:     getters.Key[components.ObjectList[Proposal]]("proposals"),
						Title:    "Unassigned Proposals",
						Subtitle: "Proposals not yet linked to a client",
						Actions: []components.PageInterface{
							&components.TableButtonFilter{Child: lamu.DynamicPage{Name: "proposals.ProposalFilter"}},
						},
						RowAttr: getters.RowAttrNavigate(lamu.RoutePath("proposals.DetailRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("$row.ID"))})),
						Columns: []components.TableColumn{
							{Label: "Title", Name: "Title", Children: []components.PageInterface{components.FieldText{Getter: getters.Key[string]("$row.Title")}}},
							{Label: "Created At", Name: "CreatedAt", Children: []components.PageInterface{components.FieldDatetime{Getter: getters.Key[time.Time]("$row.CreatedAt")}}},
							{Label: "Updated At", Name: "UpdatedAt", Children: []components.PageInterface{components.FieldDatetime{Getter: getters.Key[time.Time]("$row.UpdatedAt")}}},
						},
					},
				},
			},
		},
	}
}

func proposalDetailEditButton() components.PageInterface {
	updateFormName := getters.Static("proposals.ProposalUpdateForm")
	return components.ButtonModalForm{
		Label: "Edit Proposal",
		Icon:  "pencil",
		Name:  updateFormName,
		Url: lamu.RoutePath("proposals.UpdateRoute", map[string]getters.Getter[any]{
			"id": getters.Any(getters.Key[uint]("proposal.ID")),
		}),
		FormPostURL: lamu.RoutePath("proposals.UpdateRoute", map[string]getters.Getter[any]{
			"id": getters.Any(getters.Key[uint]("proposal.ID")),
		}),
		ModalUID: "proposal-update-modal",
		Classes:  "btn-outline",
	}
}

func registerDetail() []registry.Pair[string, components.PageInterface] {
	generatedSection := []components.PageInterface{
		components.Accordion{
			Classes: "my-2",
			Items: []components.AccordionItem{
				{
					Title: components.FieldText{Classes: "font-semibold", Getter: getters.Static("Generated Proposal")},
					Children: []components.PageInterface{
						components.ContainerColumn{Classes: "my-2", Children: []components.PageInterface{
							components.ContainerRow{Classes: "flex flex-wrap justify-between items-center gap-4 mb-4", Children: []components.PageInterface{
								components.ContainerColumn{Classes: "flex flex-wrap gap-2", Children: []components.PageInterface{
									components.ButtonDownload{Label: "Export to PDF", Link: lamu.RoutePath("proposals.ExportPdfRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("proposal.ID"))}), Classes: "btn-outline btn-secondary btn-sm"},
									components.ButtonDownload{Label: "Export to Word", Link: lamu.RoutePath("proposals.ExportDocxRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("proposal.ID"))}), Classes: "btn-outline btn-secondary btn-sm"},
									components.ButtonModalForm{Label: "Edit with AI", Name: getters.Static("proposals.AiEditModal"), Url: lamu.RoutePath("proposals.AiEditFormRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("proposal.ID"))}), FormPostURL: lamu.RoutePath("proposals.AiEditRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("proposal.ID"))}), ModalUID: "ai-edit-modal", Classes: "btn-outline btn-secondary btn-sm"},
									components.ButtonPost{Label: "Regenerate Proposal", URL: lamu.RoutePath("proposals.GenerateRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("proposal.ID"))}), Classes: "btn-outline btn-primary btn-sm"},
								}},
							}},
							components.FieldMarkdown{Classes: "ml-2", Getter: getters.Key[string]("$in.GeneratedContent")},
						}},
					},
				},
			},
		},
	}

	pendingSection := []components.PageInterface{
		components.HTMXPolling{
			URL: lamu.RoutePath("proposals.DetailRoute", map[string]getters.Getter[any]{
				"id": getters.Any(getters.Key[uint]("proposal.ID")),
			}),
			Children: []components.PageInterface{
				components.ContainerRow{Classes: "flex gap-2 items-center my-2", Children: []components.PageInterface{
					components.FieldText{Getter: getters.Static("Generating...")},
					components.ButtonPost{
						Label:   "Cancel Generation",
						URL:     lamu.RoutePath("proposals.CancelRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("proposal.ID"))}),
						Classes: "btn-outline btn-error btn-sm",
					},
				}},
			},
		},
	}

	idleSection := []components.PageInterface{
		components.ButtonPost{
			Page: components.Page{
				Key: "proposals.GenerateProposalWithAi",
			},
			Label:   "Generate Proposal with AI",
			URL:     lamu.RoutePath("proposals.GenerateRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("proposal.ID"))}),
			Classes: "btn-primary",
		},
	}

	return []registry.Pair[string, components.PageInterface]{
		{Key: "proposals.ProposalDetail", Value: components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "proposals.ProposalDetailMenu"}},
			Children: []components.PageInterface{
				components.Detail[Proposal]{
					Getter: getters.Key[Proposal]("proposal"),
					Children: []components.PageInterface{
						components.ContainerColumn{Children: []components.PageInterface{
							components.FieldTitle{Getter: getters.Key[string]("$in.Title")},
							components.LabelInline{Title: "Created At", Children: []components.PageInterface{components.FieldDatetime{Getter: getters.Key[time.Time]("$in.CreatedAt")}}},
							components.Accordion{Classes: "my-2", Items: []components.AccordionItem{
								{
									Title: components.FieldText{Classes: "font-semibold", Getter: getters.Static("Questionnaire Answers")},
									Children: []components.PageInterface{
										components.FieldKeyValue{Getter: getters.Key[datatypes.JSON]("$in.Answers")},
									},
								},
							}},
							components.ContainerColumn{Children: []components.PageInterface{
								components.ContainerRow{
									Classes: "flex flex-wrap gap-2 my-2",
									Children: []components.PageInterface{
										components.ShowIf{
											Getter:   getters.Any(getterIdleGeneration()),
											Children: idleSection,
										},
										proposalDetailEditButton(),
									},
								},
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
	return []registry.Pair[string, components.PageInterface]{{Key: "proposals.AiEditModal", Value: components.Modal{
		UID: "ai-edit-modal",
		Children: []components.PageInterface{
			components.FormComponent[Proposal]{
				Getter: getters.Key[Proposal]("proposal"),
				Attr:   getters.FormBubbling(getters.Key[string]("$get.name")),

				Title: "Edit with AI",
				ChildrenInput: []components.PageInterface{
					components.InputTextarea{Name: "GeneratedContent", Label: "Current Proposal Markdown", Getter: getters.Key[string]("$in.GeneratedContent"), Rows: 8},
					components.InputTextarea{Name: "instructions", Label: "Instructions for AI", Getter: getters.Key[string]("$in.instructions"), Rows: 4, Required: true},
				},
				ChildrenAction: []components.PageInterface{
					components.ContainerRow{Classes: "flex justify-end gap-2", Children: []components.PageInterface{
						components.ButtonSubmit{Label: "Generate", Classes: "btn-secondary"},
					}},
				},
			},
		},
	}}}
}

func registerDelete() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{{Key: "proposals.ProposalDeleteForm", Value: components.Modal{
		UID: "proposal-delete-modal",
		Children: []components.PageInterface{
			components.DeleteConfirmation{
				Title:   "Confirm Deletion",
				Message: "Are you sure you want to delete this proposal?",
				Attr:    getters.FormBubbling(getters.Key[string]("$get.name")),
			},
		},
	}}}
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
	return pluginPagesWithPatches(entries)
}
