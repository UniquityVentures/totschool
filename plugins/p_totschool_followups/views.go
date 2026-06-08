package p_totschool_followups

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/registry"
	"github.com/UniquityVentures/lamu/views"
	"github.com/UniquityVentures/totschool/plugins/p_totschool_proposals"
	"github.com/alnah/go-md2pdf"
	"gorm.io/gorm"
)

func followupDB(r *http.Request, op string) *gorm.DB {
	db, err := getters.DBFromContext(r.Context())
	if err != nil {
		slog.Error(op+": db from context", "error", err)
		return nil
	}
	return db
}

func scopeFollowupsQueryToCurrentUser(r *http.Request, query gorm.ChainInterface[Followup]) gorm.ChainInterface[Followup] {
	user, role := p_users.UserAndRoleFromContext(r.Context(), "scopeFollowupsQueryToCurrentUser")
	if user.IsSuperuser || role == "totschool_admin" {
		return query
	}
	return query.Where("created_by_id = ?", user.ID)
}

type followupDetailCtxLayer struct{}

func (followupDetailCtxLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		followup, ok := ctx.Value("followup").(Followup)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		ctx = context.WithValue(ctx, "GenerationPending", followup.GenerationID != nil)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type followupFormPatcher struct{}

func (followupFormPatcher) Patch(_ views.View, r *http.Request, formData map[string]any, formErrors map[string]error) (map[string]any, map[string]error) {
	user := p_users.UserFromContext(r.Context(), "followupFormPatcher")
	formData["CreatedByID"] = user.ID
	return formData, formErrors
}

type followupListQueryPatcher struct{}

func (followupListQueryPatcher) Patch(_ views.View, r *http.Request, query gorm.ChainInterface[Followup]) gorm.ChainInterface[Followup] {
	query = scopeFollowupsQueryToCurrentUser(r, query)
	if get, ok := r.Context().Value("$get").(map[string]any); ok {
		if raw, exists := get["Title"]; exists {
			if title, ok := raw.(string); ok && title != "" {
				query = query.Where("title ILIKE ?", "%"+title+"%")
			}
		}
		if raw, exists := get["ClientID"]; exists {
			query = addUintFilter(query, "client_id", raw)
		}
		if raw, exists := get["CreatedByID"]; exists {
			query = addUintFilter(query, "created_by_id", raw)
		}
	}
	return query
}

func addUintFilter(query gorm.ChainInterface[Followup], column string, raw any) gorm.ChainInterface[Followup] {
	switch v := raw.(type) {
	case uint:
		if v != 0 {
			return query.Where(column+" = ?", v)
		}
	case int:
		if v != 0 {
			return query.Where(column+" = ?", v)
		}
	case string:
		if v != "" {
			if id, err := strconv.ParseUint(v, 10, 32); err == nil && id != 0 {
				return query.Where(column+" = ?", uint(id))
			}
		}
	}
	return query
}

func followupPreloadPatchers() views.QueryPatchers[Followup] {
	return views.QueryPatchers[Followup]{
		{Key: "followups.preload_client", Value: views.QueryPatcherPreload[Followup]{Fields: []string{"Client"}}},
	}
}

func redirectFollowupDetail(w http.ResponseWriter, r *http.Request, idStr string) bool {
	url, err := getters.IfOr(lamu.RoutePath("followups.DetailRoute", map[string]getters.Getter[any]{
		"id": getters.Any(getters.Static(idStr)),
	}), r.Context(), "")
	if err != nil || url == "" {
		http.NotFound(w, r)
		return false
	}
	views.HtmxRedirect(w, r, url, http.StatusMovedPermanently)
	return true
}

func writeGenerationAlert(w http.ResponseWriter, r *http.Request, message string) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`<div class="alert alert-error">` + message + `</div>`))
		return
	}
	http.Error(w, message, http.StatusBadRequest)
}

func generateHandler(v *views.View) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := r.PathValue("id")
		db := followupDB(r, "generateHandler")
		if db == nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		user := p_users.UserFromContext(r.Context(), "followups.generateHandler")

		followup, err := gorm.G[Followup](db).Preload("Client", nil).Where("id = ?", idStr).First(r.Context())
		if err != nil {
			http.NotFound(w, r)
			return
		}

		proposal, err := gorm.G[p_totschool_proposals.Proposal](db).Where("client_id = ?", followup.ClientID).First(r.Context())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				writeGenerationAlert(w, r, "This client has no proposal. Create a proposal first to generate a follow-up letter.")
				return
			}
			slog.Error("generateHandler: proposal lookup failed", "error", err, "followupID", followup.ID, "clientID", followup.ClientID, "pageName", v.PageName)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		answersText, err := proposal.FormatAnswersForAI()
		if err != nil {
			slog.Error("generateHandler: FormatAnswersForAI failed", "error", err, "proposalID", proposal.ID, "pageName", v.PageName)
		}
		if err != nil || answersText == "" || len(answersText) < 10 {
			writeGenerationAlert(w, r, "This client's proposal has no usable answers. Please fill in the proposal questionnaire first.")
			return
		}

		currentDate := time.Now().Format("January 02, 2006")
		clientCity := ""
		if followup.Client.Address != nil {
			clientCity = *followup.Client.Address
		}
		userPrompt := fmt.Sprintf(`Generate a financial advisory initial follow-up letter.

FOLLOW-UP TITLE:
%s

ADVISOR:
%s

CURRENT DATE:
%s

CLIENT NAME:
%s

CLIENT CITY / ADDRESS:
%s

EXTRA CONTEXT FOR THIS FOLLOW-UP:
%s

PROPOSAL QUESTIONNAIRE RESPONSES:
%s

Use the system template and calculations. Output only the final markdown letter.`, followup.Title, user.Name, currentDate, followup.Client.Name, clientCity, followup.ExtraInfo, answersText)

		Generate(db, followup.ID, userPrompt, followupLetterSystemPrompt)
		redirectFollowupDetail(w, r, fmt.Sprintf("%d", followup.ID))
	})
}

func cancelHandler(v *views.View) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := r.PathValue("id")
		db := followupDB(r, "cancelHandler")
		if db == nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		followup, err := gorm.G[Followup](db).Where("id = ?", idStr).First(r.Context())
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if followup.GenerationID != nil {
			CancelGeneration(db, followup.ID)
		}
		redirectFollowupDetail(w, r, idStr)
	})
}

func aiEditFormHandler(v *views.View) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		db := followupDB(r, "aiEditFormHandler")
		if db == nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		followup, err := gorm.G[Followup](db).Preload("Client", nil).Where("id = ?", idStr).First(r.Context())
		if err != nil {
			http.NotFound(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), "followup", followup)
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
		db := followupDB(r, "aiEditHandler")
		if db == nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		followup, err := gorm.G[Followup](db).Where("id = ?", idStr).First(r.Context())
		if err != nil {
			http.NotFound(w, r)
			return
		}

		content := r.FormValue("GeneratedLetter")
		if content == "" {
			content = r.FormValue("generated_letter")
		}
		instructions := r.FormValue("instructions")
		if content == "" || instructions == "" {
			http.Error(w, "Missing content or instructions", http.StatusBadRequest)
			return
		}

		userPrompt := fmt.Sprintf("Here is the current follow-up letter markdown:\n\n---\n%s\n---\n\nPlease edit this follow-up letter according to these instructions: %s\n\nOutput only the edited markdown, nothing else.", content, instructions)
		Generate(db, followup.ID, userPrompt, followupLetterEditorSystemPrompt)
		redirectFollowupDetail(w, r, idStr)
	})
}

func exportDocxHandler(v *views.View) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		db := followupDB(r, "exportDocxHandler")
		if db == nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		followup, err := gorm.G[Followup](db).Where("id = ?", idStr).First(r.Context())
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if followup.GeneratedLetter == "" {
			http.Error(w, "No follow-up letter to export. Please generate the letter first.", http.StatusUnprocessableEntity)
			return
		}
		pandoc := exec.CommandContext(r.Context(), "pandoc", "-s", "-f", "markdown", "-t", "docx", "-o", "-")
		pandoc.Stdin = strings.NewReader(followup.GeneratedLetter)
		var docxOut, docxErr bytes.Buffer
		pandoc.Stdout = &docxOut
		pandoc.Stderr = &docxErr
		if err := pandoc.Run(); err != nil {
			slog.Error("exportDocxHandler: pandoc failed", "error", err, "stderr", docxErr.String(), "followupID", followup.ID, "pageName", v.PageName)
			http.Error(w, "Failed to export follow-up letter (is pandoc installed?)", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.docx"`, followup.Title))
		if _, err := w.Write(docxOut.Bytes()); err != nil {
			slog.Error("exportDocxHandler: failed to write DOCX response", "error", err, "followupID", followup.ID)
		}
	})
}

func exportPdfHandler(v *views.View) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		db := followupDB(r, "exportPdfHandler")
		if db == nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		followup, err := gorm.G[Followup](db).Where("id = ?", idStr).First(r.Context())
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if followup.GeneratedLetter == "" {
			http.Error(w, "No follow-up letter to export. Please generate the letter first.", http.StatusUnprocessableEntity)
			return
		}
		conv, err := md2pdf.NewConverter()
		if err != nil {
			slog.Error("exportPdfHandler: PDF converter unavailable", "error", err, "followupID", followup.ID, "pageName", v.PageName)
			http.Error(w, "PDF converter unavailable", http.StatusInternalServerError)
			return
		}
		defer conv.Close()
		result, err := conv.Convert(r.Context(), md2pdf.Input{
			Markdown: followup.GeneratedLetter,
			CSS: `
			@import url('https://fonts.googleapis.com/css2?family=Noto+Serif:ital,wght@0,100..900;1,100..900&family=Noto+Serif+Devanagari:wght@100..900&display=swap');
			html, body {
				font-family:
					"Noto Serif Devanagari",
					"Lohit Devanagari",
					"Noto Serif",
					"Noto Sans Devanagari",
					serif;
			}
			code, pre, kbd, samp {
				font-family: ui-monospace, "Roboto Mono", monospace;
			}
			`,
		})
		if err != nil {
			slog.Error("exportPdfHandler: PDF conversion failed", "error", err, "followupID", followup.ID, "contentLen", len(followup.GeneratedLetter), "pageName", v.PageName)
			http.Error(w, "PDF generation failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.pdf"`, followup.Title))
		if _, err := w.Write(result.PDF); err != nil {
			slog.Error("exportPdfHandler: failed to write PDF response", "error", err, "followupID", followup.ID, "pdfBytes", len(result.PDF))
		}
	})
}

func pluginViews() lamu.PluginFeatures[*views.View] {
	return pluginViewsWithPatches([]registry.Pair[string, *views.View]{
		{Key: "followups.ListView", Value: lamu.GetPageView("followups.FollowupTable").
			WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
			WithLayer("followups.list", views.LayerList[Followup]{
				Key: getters.Static("followups"),
				QueryPatchers: append(followupPreloadPatchers(),
					registry.Pair[string, views.QueryPatcher[Followup]]{Key: "followups.list", Value: followupListQueryPatcher{}},
					registry.Pair[string, views.QueryPatcher[Followup]]{Key: "followups.list_order", Value: views.QueryPatcherOrderBy[Followup]{Order: "created_at DESC, id DESC"}},
				),
			})},
		{Key: "followups.DetailView", Value: lamu.GetPageView("followups.FollowupDetail").
			WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
			WithLayer("followups.detail", views.LayerDetail[Followup]{
				Key:           getters.Static("followup"),
				PathParamKey:  getters.Static("id"),
				QueryPatchers: followupPreloadPatchers(),
			}).
			WithLayer("followups.detail_ctx", followupDetailCtxLayer{})},
		{Key: "followups.CreateView", Value: lamu.GetPageView("followups.FollowupCreateForm").
			WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
			WithLayer("followups.create_query_defaults", followupCreateQueryDefaultsLayer{}).
			WithLayer("followups.create", views.LayerCreate[Followup]{
				SuccessURL: followupCreateSuccessURL,
				FormPatchers: views.FormPatchers{
					{Key: "followups.form", Value: followupFormPatcher{}},
				},
			})},
		{Key: "followups.UpdateView", Value: lamu.GetPageView("followups.FollowupUpdateForm").
			WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
			WithLayer("followups.detail", views.LayerDetail[Followup]{
				Key:           getters.Static("followup"),
				PathParamKey:  getters.Static("id"),
				QueryPatchers: followupPreloadPatchers(),
			}).
			WithLayer("followups.update", views.LayerUpdate[Followup]{
				Key:        getters.Static("followup"),
				SuccessURL: followupUpdateSuccessURL,
				FormPatchers: views.FormPatchers{
					{Key: "followups.form", Value: followupFormPatcher{}},
				},
			})},
		{Key: "followups.DeleteView", Value: lamu.GetPageView("followups.FollowupDeleteForm").
			WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
			WithLayer("followups.detail", views.LayerDetail[Followup]{
				Key:           getters.Static("followup"),
				PathParamKey:  getters.Static("id"),
				QueryPatchers: followupPreloadPatchers(),
			}).
			WithLayer("followups.delete", views.LayerDelete[Followup]{
				Key:        getters.Static("followup"),
				SuccessURL: followupDeleteSuccessURL,
			})},
		{Key: "followups.GenerateView", Value: lamu.GetPageView("followups.FollowupDetail").
			WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
			WithLayer("followups.generate", views.MethodLayer{
				Method:  http.MethodPost,
				Handler: generateHandler,
			})},
		{Key: "followups.CancelView", Value: lamu.GetPageView("followups.FollowupDetail").
			WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
			WithLayer("followups.cancel", views.MethodLayer{
				Method:  http.MethodPost,
				Handler: cancelHandler,
			})},
		{Key: "followups.AiEditFormView", Value: lamu.GetPageView("followups.AiEditModal").
			WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
			WithLayer("followups.ai_edit_form", views.MethodLayer{
				Method:  http.MethodGet,
				Handler: aiEditFormHandler,
			})},
		{Key: "followups.AiEditView", Value: lamu.GetPageView("followups.AiEditModal").
			WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
			WithLayer("followups.ai_edit", views.MethodLayer{
				Method:  http.MethodPost,
				Handler: aiEditHandler,
			})},
		{Key: "followups.ExportPdfView", Value: lamu.GetPageView("followups.FollowupDetail").
			WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
			WithLayer("followups.export_pdf", views.MethodLayer{
				Method:  http.MethodGet,
				Handler: exportPdfHandler,
			})},
		{Key: "followups.ExportDocxView", Value: lamu.GetPageView("followups.FollowupDetail").
			WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
			WithLayer("followups.export_docx", views.MethodLayer{
				Method:  http.MethodGet,
				Handler: exportDocxHandler,
			})},
		{Key: "followups.UserSelectView", Value: lamu.GetPageView("followups.UserSelectionTable").
			WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
			WithLayer("followups.user_select", views.LayerList[p_users.User]{
				Key: getters.Static("users"),
			})},
	})
}
