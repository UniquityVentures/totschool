package p_totschool_followups

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
	g_html "maragu.dev/gomponents/html"
)

const (
	clientDetailFollowupsContextKey = "client_followups_table"
	clientDetailFollowupsLayerKey   = "followups.client_detail"
	clientDetailFollowupsTableKey   = "followups.ClientDetailFollowupsTable"
	clientDetailFollowupsTableID    = "#client-detail-followups-table"
)

func init() {
	registerClientDetailFollowupsPatch()
}

type clientFollowupsContextLayer struct{}

func (clientFollowupsContextLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client, ok := r.Context().Value("client").(p_totschool_clients.Client)
		if !ok || client.ID == 0 {
			next.ServeHTTP(w, r)
			return
		}

		db, dberr := getters.DBFromContext(r.Context())
		if dberr != nil {
			slog.Error("clientFollowupsContextLayer: db from context", "error", dberr)
			next.ServeHTTP(w, r)
			return
		}

		query := gorm.G[Followup](db).
			Preload("Client", nil).
			Where("client_id = ?", client.ID).
			Order("created_at DESC").
			Order("id DESC")
		query = scopeFollowupsQueryToCurrentUser(r, query)

		rows, err := query.Find(r.Context())
		if err != nil {
			slog.Error("clientFollowupsContextLayer: query failed", "error", err)
			next.ServeHTTP(w, r)
			return
		}

		ol := components.ObjectList[Followup]{
			Items:    rows,
			Number:   1,
			NumPages: 1,
			Total:    uint64(len(rows)),
		}
		ctx := context.WithValue(r.Context(), clientDetailFollowupsContextKey, ol)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func clientDetailFollowupColumns() []components.TableColumn {
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
	}
}

func clientDetailFollowupsSection() components.PageInterface {
	createFormName := getters.Static("followups.FollowupCreateForm")
	return containerWithID{
		Page:    components.Page{Key: clientDetailFollowupsTableKey},
		ID:      "client-detail-followups-table",
		Classes: "w-full mt-4",
		Children: []components.PageInterface{
			&components.DataTable[Followup]{
				UID:         "client-detail-followups-data-table",
				Title:       "Follow-up Letters",
				Data:        getters.Key[components.ObjectList[Followup]](clientDetailFollowupsContextKey),
				DefaultView: "Grid",
				RowAttr: getters.RowAttrNavigate(lamu.RoutePath("followups.DetailRoute", map[string]getters.Getter[any]{
					"id": getters.Any(getters.Key[uint]("$row.ID")),
				})),
				Actions: []components.PageInterface{
					&components.ButtonModalForm{
						Name: createFormName,
						Url: getters.Format(
							"%s?ClientID=%d&return=client",
							getters.Any(lamu.RoutePath("followups.CreateRoute", nil)),
							getters.Any(getters.Key[uint]("client.ID")),
						),
						FormPostURL: getters.Format(
							"%s?ClientID=%d&return=client",
							getters.Any(lamu.RoutePath("followups.CreateRoute", nil)),
							getters.Any(getters.Key[uint]("client.ID")),
						),
						ModalUID: "followup-create-modal",
						Icon:     "plus",
						Classes:  "btn-square btn-outline btn-sm",
						Attr:     getters.ModalRefreshList(getters.Static(""), getters.Static(clientDetailFollowupsTableID)),
					},
				},
				Columns: clientDetailFollowupColumns(),
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
	if containerColumnHasChildKey(column, clientDetailFollowupsTableKey) {
		return column
	}
	column.Children = append(column.Children, clientDetailFollowupsSection())
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

func registerClientDetailFollowupsPatch() {
	patchPluginView("clients.DetailView", func(v *views.View) *views.View {
		if viewHasLayer(v, clientDetailFollowupsLayerKey) {
			return v
		}
		return v.InsertLayerAfter("clients.detail", clientDetailFollowupsLayerKey, clientFollowupsContextLayer{})
	})

	patchPluginPage("clients.ClientDetail", patchClientDetailPage)
}
