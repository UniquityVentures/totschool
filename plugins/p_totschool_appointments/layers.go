package p_totschool_appointments

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/views"
	"github.com/UniquityVentures/totschool/plugins/p_totschool_clients"
	"gorm.io/gorm"
)

// appointmentCreateQueryDefaultsLayer merges ?ClientID= into $in on GET so the create
// form opened from client detail pre-fills the client and meeting location.
type appointmentCreateQueryDefaultsLayer struct{}

func appointmentCreateDefaultsFromClientID(ctx context.Context, clientID uint) map[string]any {
	vals := map[string]any{"ClientID": clientID}
	db, err := getters.DBFromContext(ctx)
	if err != nil {
		slog.Error("appointmentCreateDefaultsFromClientID: db from context", "error", err)
		return vals
	}
	client, err := gorm.G[p_totschool_clients.Client](db).Where("id = ?", clientID).First(ctx)
	if err != nil {
		slog.Error("appointmentCreateDefaultsFromClientID: load client", "error", err, "clientID", clientID)
		return vals
	}
	if client.Address != nil {
		if loc := strings.TrimSpace(*client.Address); loc != "" {
			vals["Location"] = loc
		}
	}
	return vals
}

func (appointmentCreateQueryDefaultsLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}
		var vals map[string]any
		if cid := r.URL.Query().Get("ClientID"); cid != "" {
			if id64, err := strconv.ParseUint(cid, 10, 32); err == nil && id64 != 0 {
				vals = appointmentCreateDefaultsFromClientID(r.Context(), uint(id64))
			}
		}
		if len(vals) == 0 {
			next.ServeHTTP(w, r)
			return
		}
		ctx := views.ContextWithErrorsAndValues(r.Context(), vals, nil)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func requestReturnClient(r *http.Request) bool {
	return r.URL.Query().Get("return") == "client"
}
