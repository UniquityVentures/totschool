package p_totschool_appointments

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/registry"
	"github.com/UniquityVentures/lamu/views"
	"gorm.io/gorm"
)

// AppointmentDetailCtxLayer enriches detail context after LayerDetail loads "appointment".
type AppointmentDetailCtxLayer struct{}

func (AppointmentDetailCtxLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		rawAppt := ctx.Value("appointment")
		appointment, ok := rawAppt.(Appointment)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		db, dberr := getters.DBFromContext(ctx)
		if dberr != nil {
			slog.Error("AppointmentDetailCtxLayer: db from context", "error", dberr)
			next.ServeHTTP(w, r)
			return
		}

		if appointment.GenerationID != nil {
			ctx = context.WithValue(ctx, "GenerationPending", true)
		} else {
			ctx = context.WithValue(ctx, "GenerationPending", false)
		}

		overlapping := appointment.GetOverlappingAppointments(db)
		if len(overlapping) > 0 {
			overlapList := []map[string]any{}
			for _, o := range overlapping {
				overlapList = append(overlapList, map[string]any{
					"ID":   o.ID,
					"Name": o.Name,
					"Date": o.Datetime,
				})
			}
			ctx = context.WithValue(ctx, "OverlapWarningList", overlapList)
			ctx = context.WithValue(ctx, "OverlapWarning", true)
		} else {
			ctx = context.WithValue(ctx, "OverlapWarning", false)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func redirectAppointmentDetail(w http.ResponseWriter, r *http.Request, idStr string) bool {
	url, err := getters.IfOr(lamu.RoutePath("appointments.DetailRoute", map[string]getters.Getter[any]{
		"id": getters.Any(getters.Static(idStr)),
	}), r.Context(), "")
	if err != nil || url == "" {
		http.NotFound(w, r)
		return false
	}
	views.HtmxRedirect(w, r, url, http.StatusMovedPermanently)
	return true
}

func generateHandler(v *views.View) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := r.PathValue("id")
		db, dberr := getters.DBFromContext(r.Context())
		if dberr != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		user := p_users.UserFromContext(r.Context(), "appointments.generateHandler")

		appointment, err := gorm.G[Appointment](db).Where("id = ?", idStr).First(r.Context())
		if err != nil {
			http.NotFound(w, r)
			return
		}

		content, systemPrompt := buildLetterContent(db, r.Context(), &appointment, user.Name)
		Generate(db, appointment.ID, content, systemPrompt)

		redirectAppointmentDetail(w, r, idStr)
	})
}

func cancelHandler(v *views.View) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := r.PathValue("id")
		db, dberr := getters.DBFromContext(r.Context())
		if dberr != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		appointment, err := gorm.G[Appointment](db).Where("id = ?", idStr).First(r.Context())
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if appointment.GenerationID != nil {
			CancelGeneration(db, appointment.ID)
		}

		redirectAppointmentDetail(w, r, idStr)
	})
}

func aiEditFormHandler(v *views.View) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		db, dberr := getters.DBFromContext(r.Context())
		if dberr != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		appointment, err := gorm.G[Appointment](db).Where("id = ?", idStr).First(r.Context())
		if err != nil {
			http.NotFound(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "appointment", appointment)
		v.RenderPage(w, r.WithContext(ctx))
	})
}

func aiEditHandler(v *views.View) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := r.PathValue("id")
		db, dberr := getters.DBFromContext(r.Context())
		if dberr != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		appointment, err := gorm.G[Appointment](db).Where("id = ?", idStr).First(r.Context())
		if err != nil {
			http.NotFound(w, r)
			return
		}

		content := r.FormValue("generated_letter")
		instructions := r.FormValue("instructions")
		if content == "" || instructions == "" {
			http.Error(w, "Missing content or instructions", http.StatusBadRequest)
			return
		}

		userPrompt := "Here is the current letter content:\n\n" + content + "\n\nPlease edit this letter according to these instructions: " + instructions + "\n\nOutput only the edited text, nothing else."
		Generate(db, appointment.ID, userPrompt, letterEditorSystemPrompt)

		redirectAppointmentDetail(w, r, idStr)
	})
}

type appointmentFormCreatedByPatcher struct{}

func (appointmentFormCreatedByPatcher) Patch(_ views.View, r *http.Request, formData map[string]any, formErrors map[string]error) (map[string]any, map[string]error) {
	user := p_users.UserFromContext(r.Context(), "appointmentFormCreatedByPatcher")
	formData["CreatedByID"] = user.ID
	return formData, formErrors
}

func scopeAppointmentsQueryToCurrentUser(r *http.Request, query gorm.ChainInterface[Appointment]) gorm.ChainInterface[Appointment] {
	user, role := p_users.UserAndRoleFromContext(r.Context(), "scopeAppointmentsQueryToCurrentUser")
	if user.IsSuperuser || role == "totschool_admin" {
		return query
	}
	return query.Where("created_by_id = ?", user.ID)
}

type appointmentListQueryPatcher struct{}

func (appointmentListQueryPatcher) Patch(_ views.View, r *http.Request, query gorm.ChainInterface[Appointment]) gorm.ChainInterface[Appointment] {
	ctx := r.Context()
	query = scopeAppointmentsQueryToCurrentUser(r, query)

	if get, ok := ctx.Value("$get").(map[string]any); ok {
		if val, exists := get["Overlapping"]; exists {
			if b, ok := val.(bool); ok && b {
				query = WithOverlappingFilterChain(query)
			}
		}
		if raw, exists := get["Date"]; exists && raw != nil {
			query = applyDateFilterChain(raw, query)
		}
	}

	return query
}

type appointmentTimelineQueryPatcher struct{}

func (appointmentTimelineQueryPatcher) Patch(_ views.View, r *http.Request, query gorm.ChainInterface[Appointment]) gorm.ChainInterface[Appointment] {
	ctx := r.Context()
	query = scopeAppointmentsQueryToCurrentUser(r, query)

	if get, ok := ctx.Value("$get").(map[string]any); ok {
		if raw, exists := get["Date"]; exists && raw != nil {
			switch d := raw.(type) {
			case time.Time:
				if !d.IsZero() {
					return applyDateFilterChain(raw, query)
				}
			case string:
				if d != "" {
					return applyDateFilterChain(raw, query)
				}
			}
		}
	}
	return applyDateFilterChain(time.Now(), query)
}

func applyDateFilterChain(raw any, query gorm.ChainInterface[Appointment]) gorm.ChainInterface[Appointment] {
	switch d := raw.(type) {
	case time.Time:
		start := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
		end := start.Add(24 * time.Hour)
		query = query.Where("datetime >= ? AND datetime < ?", start, end)
	case string:
		if d != "" {
			if parsed, err := time.Parse("2006-01-02", d); err == nil {
				start := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, parsed.Location())
				end := start.Add(24 * time.Hour)
				query = query.Where("datetime >= ? AND datetime < ?", start, end)
			}
		}
	}
	return query
}

func pluginViews() lamu.PluginFeatures[*views.View] {
	return lamu.PluginFeatures[*views.View]{
		Entries: []registry.Pair[string, *views.View]{
			{
				Key: "appointments.ListView",
				Value: lamu.GetPageView("appointments.AppointmentTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("appointments.list", views.LayerList[Appointment]{
						Key: getters.Static("appointments"),
						QueryPatchers: views.QueryPatchers[Appointment]{
							{Key: "appointments.list", Value: appointmentListQueryPatcher{}},
						},
					}),
			},
			{
				Key: "appointments.DetailView",
				Value: lamu.GetPageView("appointments.AppointmentDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("appointments.detail", views.LayerDetail[Appointment]{
						Key:          getters.Static("appointment"),
						PathParamKey: getters.Static("id"),
					}).
					WithLayer("appointments.detail_ctx", AppointmentDetailCtxLayer{}),
			},
			{
				Key: "appointments.CreateView",
				Value: lamu.GetPageView("appointments.AppointmentCreateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("appointments.create", views.LayerCreate[Appointment]{
						SuccessURL: lamu.RoutePath("appointments.DetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
						FormPatchers: views.FormPatchers{
							{Key: "appointments.form", Value: appointmentFormCreatedByPatcher{}},
						},
					}),
			},
			{
				Key: "appointments.UpdateView",
				Value: lamu.GetPageView("appointments.AppointmentUpdateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("appointments.detail", views.LayerDetail[Appointment]{
						Key:          getters.Static("appointment"),
						PathParamKey: getters.Static("id"),
					}).
					WithLayer("appointments.update", views.LayerUpdate[Appointment]{
						Key: getters.Static("appointment"),
						SuccessURL: lamu.RoutePath("appointments.DetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("appointment.ID")),
						}),
						FormPatchers: views.FormPatchers{
							{Key: "appointments.form", Value: appointmentFormCreatedByPatcher{}},
						},
					}),
			},
			{
				Key: "appointments.DeleteView",
				Value: lamu.GetPageView("appointments.AppointmentDeleteForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("appointments.detail", views.LayerDetail[Appointment]{
						Key:          getters.Static("appointment"),
						PathParamKey: getters.Static("id"),
					}).
					WithLayer("appointments.delete", views.LayerDelete[Appointment]{
						Key:        getters.Static("appointment"),
						SuccessURL: lamu.RoutePath("appointments.ListRoute", nil),
					}),
			},
			{
				Key: "appointments.GenerateView",
				Value: lamu.GetPageView("appointments.AppointmentDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("appointments.generate", views.MethodLayer{
						Method:  http.MethodPost,
						Handler: generateHandler,
					}),
			},
			{
				Key: "appointments.CancelView",
				Value: lamu.GetPageView("appointments.AppointmentDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("appointments.cancel", views.MethodLayer{
						Method:  http.MethodPost,
						Handler: cancelHandler,
					}),
			},
			{
				Key: "appointments.AiEditFormView",
				Value: lamu.GetPageView("appointments.AiEditModal").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("appointments.ai_edit_form", views.MethodLayer{
						Method:  http.MethodGet,
						Handler: aiEditFormHandler,
					}),
			},
			{
				Key: "appointments.AiEditView",
				Value: lamu.GetPageView("appointments.AiEditModal").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("appointments.ai_edit", views.MethodLayer{
						Method:  http.MethodPost,
						Handler: aiEditHandler,
					}),
			},
			{
				Key: "appointments.SelectView",
				Value: lamu.GetPageView("appointments.AppointmentSelectionTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("appointments.select_list", views.LayerList[Appointment]{
						Key: getters.Static("appointments"),
					}),
			},
			{
				Key: "appointments.CardTimelineView",
				Value: lamu.GetPageView("appointments.AppointmentCardTimeline").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("appointments.timeline", views.LayerList[Appointment]{
						Key: getters.Static("appointments"),
						QueryPatchers: views.QueryPatchers[Appointment]{
							{Key: "appointments.timeline", Value: appointmentTimelineQueryPatcher{}},
							{Key: "appointments.timeline_order", Value: views.QueryPatcherOrderBy[Appointment]{Order: "datetime ASC"}},
						},
					}),
			},
		},
	}
}
