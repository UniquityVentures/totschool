package p_totschool_proposals

import (
	"time"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
	"gorm.io/datatypes"
)

func registerMenus() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "proposals.ProposalMenu", Value: components.SidebarMenu{
			Title: getters.Static("Proposals"),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("Back to All Apps"),
				Url:   lamu.RoutePath("dashboard.AppsPage", nil),
			},
			Children: []components.PageInterface{
				components.SidebarMenuItem{Title: getters.Static("All Proposals"), Url: lamu.RoutePath("proposals.ListRoute", nil)},
				components.SidebarMenuItem{Title: getters.Static("Create Proposal"), Url: lamu.RoutePath("proposals.CreateRoute", nil)},
			},
		}},
		{Key: "proposals.ProposalDetailMenu", Value: components.SidebarMenu{
			Title: getters.Format("Proposal: %s", getters.Any(getters.Key[string]("proposal.Title"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("Back to all Proposals"),
				Url:   lamu.RoutePath("proposals.ListRoute", nil),
			},
			Children: []components.PageInterface{
				components.SidebarMenuItem{Title: getters.Static("Proposal Detail"), Url: lamu.RoutePath("proposals.DetailRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("proposal.ID"))})},
				components.SidebarMenuItem{Title: getters.Static("Edit Proposal"), Url: lamu.RoutePath("proposals.UpdateRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("proposal.ID"))})},
			},
		}},
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

func registerForms() []registry.Pair[string, components.PageInterface] {
	createFormName := getters.Static("proposals.ProposalCreateForm")
	updateFormName := getters.Static("proposals.ProposalUpdateForm")
	deleteFormName := getters.Static("proposals.ProposalDeleteForm")
	return []registry.Pair[string, components.PageInterface]{
		{Key: "proposals.ProposalFormFields", Value: components.ContainerColumn{
			Children: []components.PageInterface{
				components.ContainerColumn{Children: []components.PageInterface{components.InputText{Label: "Proposal Title", Name: "Title", Required: true, Getter: getters.Key[string]("$in.Title")}, components.InputKeyValue{Getter: getters.Key[datatypes.JSON]("$in.Answers"), Keys: getters.Static(QUESTIONS), Name: "Answers"}}},
			},
		}},
		{Key: "proposals.ProposalCreateForm", Value: components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "proposals.ProposalMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      createFormName,
					ActionURL: lamu.RoutePath("proposals.CreateRoute", nil),
					Children: []components.PageInterface{
						components.FormComponent[Proposal]{
							Attr: getters.FormBubbling(createFormName),

							Title:          "Create Proposal",
							Subtitle:       "Fill in the questionnaire answers",
							ChildrenInput:  []components.PageInterface{components.InputText{Label: "Proposal Title", Name: "Title", Required: true, Getter: getters.Key[string]("$in.Title")}, components.InputKeyValue{Getter: getters.Key[datatypes.JSON]("$in.Answers"), Keys: getters.Static(QUESTIONS), Name: "Answers"}},
							ChildrenAction: []components.PageInterface{components.ButtonSubmit{Label: "Save Proposal"}},
						},
					},
				},
			},
		}},
		{Key: "proposals.ProposalUpdateForm", Value: components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "proposals.ProposalDetailMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      updateFormName,
					ActionURL: lamu.RoutePath("proposals.UpdateRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("proposal.ID"))}),
					Children: []components.PageInterface{
						components.FormComponent[Proposal]{
							Getter: getters.Key[Proposal]("proposal"),
							Attr:   getters.FormBubbling(updateFormName),

							Title:    "Edit Proposal",
							Subtitle: "Update questionnaire answers",
							ChildrenInput: []components.PageInterface{
								components.InputText{Label: "Title", Name: "Title", Getter: getters.Key[string]("$in.Title")},
								components.InputKeyValue{
									Getter: getters.Key[datatypes.JSON]("$in.Answers"),
									Keys:   getters.Static(QUESTIONS),
									Name:   "Answers",
								},
							},
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
	return []registry.Pair[string, components.PageInterface]{{Key: "proposals.ProposalTable", Value: components.ShellScaffold{
		Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "proposals.ProposalMenu"}},
		Children: []components.PageInterface{
			components.DataTable[Proposal]{
				UID:      "proposal-table",
				Data:     getters.Key[components.ObjectList[Proposal]]("proposals"),
				Title:    "Proposals",
				Subtitle: "List of financial proposals",
				Actions: []components.PageInterface{
					&components.TableButtonFilter{Child: lamu.DynamicPage{Name: "proposals.ProposalFilter"}},
					&components.TableButtonCreate{Link: lamu.RoutePath("proposals.CreateRoute", nil)},
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
			Label: "Generate Proposal with AI",
			URL:   lamu.RoutePath("proposals.GenerateRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("proposal.ID"))}), Classes: "btn-primary",
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
	return lamu.PluginFeatures[components.PageInterface]{Entries: entries}
}
