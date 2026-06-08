package p_totschool_proposals

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/views"
	"github.com/UniquityVentures/totschool/plugins/p_totschool_clients"
	"gorm.io/gorm"
	"maragu.dev/gomponents"
	g_html "maragu.dev/gomponents/html"
)

func clientDetailModalButtonAttr(sectionSelector string) getters.Getter[gomponents.Node] {
	refresh := getters.ModalRefreshList(getters.Static(""), getters.Static(sectionSelector))
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
	clientDetailProposalContextKey = "client_proposal"
	clientDetailHasProposalKey     = "client_has_proposal"
	clientDetailProposalsLayerKey  = "proposals.client_detail"
	clientDetailProposalSectionKey = "proposals.ClientDetailProposalSection"
	clientDetailProposalSectionID  = "#client-detail-proposal-section"
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

		query := gorm.G[Proposal](db).Where("client_id = ?", client.ID)
		query = scopeProposalsQueryToCurrentUser(r, query)

		proposal, err := query.First(r.Context())
		hasProposal := err == nil && proposal.ID != 0
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("clientProposalsContextLayer: query failed", "error", err)
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), clientDetailProposalContextKey, proposal)
		ctx = context.WithValue(ctx, clientDetailHasProposalKey, hasProposal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func clientDetailNoProposal() getters.Getter[bool] {
	return func(ctx context.Context) (bool, error) {
		has, err := getters.Key[bool](clientDetailHasProposalKey)(ctx)
		if err != nil {
			return false, err
		}
		return !has, nil
	}
}

func clientDetailProposalCreateButton() components.PageInterface {
	createFormName := getters.Static("proposals.ProposalCreateForm")
	return &components.ButtonModalForm{
		Label: "Create Proposal",
		Icon:  "plus",
		Name:  createFormName,
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
		Classes:  "btn-outline btn-sm",
		Attr:     getters.ModalRefreshList(getters.Static(""), getters.Static(clientDetailProposalSectionID)),
	}
}

func clientDetailProposalCard() components.PageInterface {
	return components.ContainerColumn{
		Classes: "border border-base-300 rounded-box bg-base-100 p-2",
		Children: []components.PageInterface{
			components.Detail[Proposal]{
				Getter: getters.Key[Proposal](clientDetailProposalContextKey),
				Children: []components.PageInterface{
					components.ContainerRow{
						Classes: "flex flex-wrap justify-between items-start",
						Children: []components.PageInterface{
							components.ContainerColumn{
								Classes: "min-w-0",
								Children: []components.PageInterface{
									components.FieldText{
										Classes: "font-semibold text-lg",
										Getter:  getters.Key[string]("$in.Title"),
									},
									components.LabelInline{
										Title: "Created",
										Children: []components.PageInterface{
											components.FieldDatetime{Getter: getters.Key[time.Time]("$in.CreatedAt")},
										},
									},
								},
							},
							components.ContainerRow{
								Classes: "flex flex-wrap gap-2 shrink-0",
								Children: []components.PageInterface{
									components.ButtonLink{
										Label: "View Proposal",
										Link: lamu.RoutePath("proposals.DetailRoute", map[string]getters.Getter[any]{
											"id": getters.Any(getters.Key[uint]("$in.ID")),
										}),
										Classes: "btn-outline btn-sm",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func clientDetailProposalSection() components.PageInterface {
	return containerWithID{
		Page:    components.Page{Key: clientDetailProposalSectionKey},
		ID:      "client-detail-proposal-section",
		Classes: "w-full mt-6 gap-3 flex flex-col",
		Children: []components.PageInterface{
			components.ContainerRow{
				Classes: "flex flex-wrap items-center justify-between",
				Children: []components.PageInterface{
					components.FieldTitle{Getter: getters.Static("Proposal")},
				},
			},
			components.ShowIf{
				Getter:   getters.Any(getters.Key[bool](clientDetailHasProposalKey)),
				Children: []components.PageInterface{clientDetailProposalCard()},
			},
			components.ShowIf{
				Getter: getters.Any(clientDetailNoProposal()),
				Children: []components.PageInterface{
					components.ContainerColumn{
						Classes: "border border-dashed border-base-300 rounded-box bg-base-100 p-6 gap-3 items-center text-center",
						Children: []components.PageInterface{
							components.FieldText{
								Classes: "text-base-content/70",
								Getter:  getters.Static("No proposal yet for this client."),
							},
							clientDetailProposalCreateButton(),
						},
					},
				},
			},
		},
	}
}

type containerWithID struct {
	components.Page
	ID       string
	Classes  string
	Children []components.PageInterface
}

func (e containerWithID) Build(ctx context.Context) gomponents.Node {
	group := gomponents.Group{}
	for _, child := range e.Children {
		group = append(group, components.Render(child, ctx))
	}
	return g_html.Div(
		g_html.ID(e.ID),
		g_html.Class(e.Classes),
		group,
	)
}

func (e containerWithID) GetKey() string {
	return e.Key
}

func (e containerWithID) GetRoles() []string {
	return e.Roles
}

func (e containerWithID) GetChildren() []components.PageInterface {
	return e.Children
}

func (e *containerWithID) SetChildren(children []components.PageInterface) {
	e.Children = children
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
	if containerColumnHasChildKey(column, clientDetailProposalSectionKey) {
		return column
	}
	column.Children = append(column.Children, clientDetailProposalSection())
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
