package p_totschool_proposals

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/views"
	"github.com/UniquityVentures/totschool/plugins/p_totschool_clients"
	"gorm.io/gorm"
	"maragu.dev/gomponents"
)

func clientDetailModalButtonAttr(tableSelector string) getters.Getter[gomponents.Node] {
	refresh := getters.ModalRefreshList(getters.Static(""), getters.Static(tableSelector))
	return func(ctx context.Context) (gomponents.Node, error) {
		nodes, err := refresh(ctx)
		if err != nil {
			return nil, err
		}
		return gomponents.Group{
			gomponents.Attr("@click.stop", ""),
			nodes,
		}, nil
	}
}

const (
	clientDetailProposalsContextKey = "client_proposals_table"
	clientDetailProposalsLayerKey   = "proposals.client_detail"
	clientDetailProposalsTableKey   = "proposals.ClientDetailProposalsTable"
)

func init() {
	registerClientDetailProposalsPatch()
}

type clientProposalsContextLayer struct{}

func (clientProposalsContextLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client, ok := r.Context().Value("client").(p_totschool_clients.Client)
		if !ok || client.ID == 0 {
			next.ServeHTTP(w, r)
			return
		}

		db, dberr := getters.DBFromContext(r.Context())
		if dberr != nil {
			slog.Error("clientProposalsContextLayer: db from context", "error", dberr)
			next.ServeHTTP(w, r)
			return
		}

		query := gorm.G[Proposal](db).
			Where("client_id = ?", client.ID).
			Order("created_at DESC").
			Order("id DESC")
		query = scopeProposalsQueryToCurrentUser(r, query)

		rows, err := query.Find(r.Context())
		if err != nil {
			slog.Error("clientProposalsContextLayer: query failed", "error", err)
			next.ServeHTTP(w, r)
			return
		}

		ol := components.ObjectList[Proposal]{
			Items:    rows,
			Number:   1,
			NumPages: 1,
			Total:    uint64(len(rows)),
		}
		ctx := context.WithValue(r.Context(), clientDetailProposalsContextKey, ol)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func clientDetailProposalColumns() []components.TableColumn {
	updateFormName := getters.Static("proposals.ProposalUpdateForm")

	return []components.TableColumn{
		{Label: "Title", Name: "Title", Children: []components.PageInterface{
			components.FieldText{Getter: getters.Key[string]("$row.Title")},
		}},
		{Label: "Created By", Name: "CreatedBy", Children: []components.PageInterface{
			components.FieldText{Getter: getters.ForeignKey[p_users.User, uint, string](getters.Key[uint]("$row.CreatedByID"), "Name")},
		}},
		{Label: "Created At", Name: "CreatedAt", Children: []components.PageInterface{
			components.FieldDatetime{Getter: getters.Key[time.Time]("$row.CreatedAt")},
		}},
		{
			Label: "",
			Name:  "Actions",
			Children: []components.PageInterface{
				components.ContainerRow{
					Classes: "flex gap-1",
					Children: []components.PageInterface{
						components.ButtonModalForm{
							Label: "Edit",
							Icon:  "pencil",
							Name:  updateFormName,
							Url: getters.Format(
								"%s?return=client",
								getters.Any(lamu.RoutePath("proposals.UpdateRoute", map[string]getters.Getter[any]{
									"id": getters.Any(getters.Key[uint]("$row.ID")),
								})),
							),
							FormPostURL: getters.Format(
								"%s?return=client",
								getters.Any(lamu.RoutePath("proposals.UpdateRoute", map[string]getters.Getter[any]{
									"id": getters.Any(getters.Key[uint]("$row.ID")),
								})),
							),
							ModalUID: "proposal-update-modal",
							Classes:  "btn-outline btn-sm m-2",
							Attr:     clientDetailModalButtonAttr("#client-detail-proposals-table"),
						},
					},
				},
			},
		},
	}
}

func clientDetailProposalsSection() components.PageInterface {
	createFormName := getters.Static("proposals.ProposalCreateForm")
	return &components.DataTable[Proposal]{
		Page:        components.Page{Key: clientDetailProposalsTableKey},
		UID:         "client-detail-proposals-table",
		Title:       "Proposals",
		Classes:     "w-full mt-4",
		Data:        getters.Key[components.ObjectList[Proposal]](clientDetailProposalsContextKey),
		DefaultView: "Grid",
		RowAttr: getters.RowAttrClickWithClass(
			getters.Format(
				"if (!$event.target.closest('button, a, input, select, textarea, .fk-modal-host')) { htmx.ajax('GET', '%s', {target: 'body', swap: 'outerHTML'}) }",
				getters.Any(lamu.RoutePath("proposals.DetailRoute", map[string]getters.Getter[any]{
					"id": getters.Any(getters.Key[uint]("$row.ID")),
				})),
			),
			nil,
		),
		Actions: []components.PageInterface{
			&components.ButtonModalForm{
				Name: createFormName,
				Url: getters.Format(
					"%s?ClientID=%d&return=client",
					getters.Any(lamu.RoutePath("proposals.CreateRoute", nil)),
					getters.Any(getters.Key[uint]("client.ID")),
				),
				FormPostURL: getters.Format(
					"%s?ClientID=%d&return=client",
					getters.Any(lamu.RoutePath("proposals.CreateRoute", nil)),
					getters.Any(getters.Key[uint]("client.ID")),
				),
				ModalUID: "proposal-create-modal",
				Icon:     "plus",
				Classes:  "btn-square btn-outline btn-sm",
				Attr:     getters.ModalRefreshList(getters.Static(""), getters.Static("#client-detail-proposals-table")),
			},
		},
		Columns: clientDetailProposalColumns(),
	}
}

func viewHasLayer(v *views.View, name string) bool {
	for _, layer := range v.Layers {
		if layer.Key == name {
			return true
		}
	}
	return false
}

func containerColumnHasChildKey(column components.ContainerColumn, key string) bool {
	for _, child := range column.Children {
		if child.GetKey() == key {
			return true
		}
	}
	return false
}

func patchClientDetailContentColumn(column components.ContainerColumn) components.ContainerColumn {
	if containerColumnHasChildKey(column, clientDetailProposalsTableKey) {
		return column
	}
	column.Children = append(column.Children, clientDetailProposalsSection())
	return column
}

func patchClientDetailPage(page components.PageInterface) components.PageInterface {
	switch p := page.(type) {
	case *components.ShellScaffold:
		s := *p
		components.ReplaceChild(&s, "clients.ClientDetailContent", patchClientDetailContentColumn)
		return &s
	case components.ShellScaffold:
		s := p
		components.ReplaceChild(&s, "clients.ClientDetailContent", patchClientDetailContentColumn)
		return s
	default:
		log.Panic("clients.ClientDetail was not ShellScaffold")
		return page
	}
}

func registerClientDetailProposalsPatch() {
	patchPluginView("clients.DetailView", func(v *views.View) *views.View {
		if viewHasLayer(v, clientDetailProposalsLayerKey) {
			return v
		}
		return v.InsertLayerAfter("clients.detail", clientDetailProposalsLayerKey, clientProposalsContextLayer{})
	})

	patchPluginPage("clients.ClientDetail", patchClientDetailPage)
}
