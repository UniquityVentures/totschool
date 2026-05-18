package p_totschool_appointments

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
	clientDetailAppointmentsContextKey = "client_appointments_table"
	clientDetailAppointmentsLayerKey   = "appointments.client_detail"
	clientDetailAppointmentsTableKey   = "appointments.ClientDetailAppointmentsTable"
)

func init() {
	registerClientDetailAppointmentsPatch()
}

type clientAppointmentsContextLayer struct{}

func (clientAppointmentsContextLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client, ok := r.Context().Value("client").(p_totschool_clients.Client)
		if !ok || client.ID == 0 {
			next.ServeHTTP(w, r)
			return
		}

		db, dberr := getters.DBFromContext(r.Context())
		if dberr != nil {
			slog.Error("clientAppointmentsContextLayer: db from context", "error", dberr)
			next.ServeHTTP(w, r)
			return
		}

		query := gorm.G[Appointment](db).
			Preload("Client", nil).
			Where("client_id = ?", client.ID).
			Order("datetime DESC").
			Order("id DESC")
		query = scopeAppointmentsQueryToCurrentUser(r, query)

		rows, err := query.Find(r.Context())
		if err != nil {
			slog.Error("clientAppointmentsContextLayer: query failed", "error", err)
			next.ServeHTTP(w, r)
			return
		}

		ol := components.ObjectList[Appointment]{
			Items:    rows,
			Number:   1,
			NumPages: 1,
			Total:    uint64(len(rows)),
		}
		ctx := context.WithValue(r.Context(), clientDetailAppointmentsContextKey, ol)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func clientDetailAppointmentColumns() []components.TableColumn {
	updateFormName := getters.Static("appointments.AppointmentUpdateForm")
	deleteFormName := getters.Static("appointments.AppointmentDeleteForm")

	return []components.TableColumn{
		{Label: "Status", Name: "Status", Children: []components.PageInterface{
			components.FieldText{Getter: appointmentStatusLabelFromRow()},
		}},
		{Label: "Date & Time", Name: "Datetime", Children: []components.PageInterface{
			components.FieldDatetime{Getter: getters.Key[time.Time]("$row.Datetime")},
		}},
		{Label: "Remarks", Name: "Remarks", Children: []components.PageInterface{
			components.FieldText{Getter: getters.Key[string]("$row.Remarks")},
		}},
		{Label: "Created By", Name: "CreatedBy", Children: []components.PageInterface{
			components.FieldText{Getter: getters.ForeignKey[p_users.User, uint, string](getters.Key[uint]("$row.CreatedByID"), "Name")},
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
								getters.Any(lamu.RoutePath("appointments.UpdateRoute", map[string]getters.Getter[any]{
									"id": getters.Any(getters.Key[uint]("$row.ID")),
								})),
							),
							FormPostURL: getters.Format(
								"%s?return=client",
								getters.Any(lamu.RoutePath("appointments.UpdateRoute", map[string]getters.Getter[any]{
									"id": getters.Any(getters.Key[uint]("$row.ID")),
								})),
							),
							ModalUID: "appointment-update-modal",
							Classes:  "btn-outline btn-sm",
							Attr:     clientDetailModalButtonAttr("#client-detail-appointments-table"),
						},
						components.ButtonModalForm{
							Label: "Delete",
							Icon:  "trash",
							Name:  deleteFormName,
							Url: getters.Format(
								"%s?return=client",
								getters.Any(lamu.RoutePath("appointments.DeleteRoute", map[string]getters.Getter[any]{
									"id": getters.Any(getters.Key[uint]("$row.ID")),
								})),
							),
							FormPostURL: getters.Format(
								"%s?return=client",
								getters.Any(lamu.RoutePath("appointments.DeleteRoute", map[string]getters.Getter[any]{
									"id": getters.Any(getters.Key[uint]("$row.ID")),
								})),
							),
							ModalUID: "appointment-delete-modal",
							Classes:  "btn-error btn-sm",
							Attr:     clientDetailModalButtonAttr("#client-detail-appointments-table"),
						},
					},
				},
			},
		},
	}
}

func clientDetailAppointmentsSection() components.PageInterface {
	createFormName := getters.Static("appointments.AppointmentCreateForm")
	return &components.DataTable[Appointment]{
		Page:        components.Page{Key: clientDetailAppointmentsTableKey},
		UID:         "client-detail-appointments-table",
		Title:       "Appointments",
		Classes:     "w-full mt-4",
		Data:        getters.Key[components.ObjectList[Appointment]](clientDetailAppointmentsContextKey),
		DefaultView: "Grid",
		RowAttr: getters.RowAttrNavigate(lamu.RoutePath("appointments.DetailRoute", map[string]getters.Getter[any]{
			"id": getters.Any(getters.Key[uint]("$row.ID")),
		})),
		Actions: []components.PageInterface{
			&components.ButtonModalForm{
				Name: createFormName,
				Url: getters.Format(
					"%s?ClientID=%d&return=client",
					getters.Any(lamu.RoutePath("appointments.CreateRoute", nil)),
					getters.Any(getters.Key[uint]("client.ID")),
				),
				FormPostURL: getters.Format(
					"%s?ClientID=%d&return=client",
					getters.Any(lamu.RoutePath("appointments.CreateRoute", nil)),
					getters.Any(getters.Key[uint]("client.ID")),
				),
				ModalUID: "appointment-create-modal",
				Icon:     "plus",
				Classes:  "btn-square btn-outline btn-sm",
				Attr:     getters.ModalRefreshList(getters.Static(""), getters.Static("#client-detail-appointments-table")),
			},
		},
		Columns: clientDetailAppointmentColumns(),
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
	if containerColumnHasChildKey(column, clientDetailAppointmentsTableKey) {
		return column
	}
	column.Children = append(column.Children, clientDetailAppointmentsSection())
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

func registerClientDetailAppointmentsPatch() {
	patchPluginView("clients.DetailView", func(v *views.View) *views.View {
		if viewHasLayer(v, clientDetailAppointmentsLayerKey) {
			return v
		}
		return v.InsertLayerAfter("clients.detail", clientDetailAppointmentsLayerKey, clientAppointmentsContextLayer{})
	})

	patchPluginPage("clients.ClientDetail", patchClientDetailPage)
}
