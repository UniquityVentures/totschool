package p_totschool_followups

import (
	"net/http"
	"strconv"

	"github.com/UniquityVentures/lamu/views"
)

type followupCreateQueryDefaultsLayer struct{}

func (followupCreateQueryDefaultsLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}
		vals := map[string]any{}
		if cid := r.URL.Query().Get("ClientID"); cid != "" {
			if id64, err := strconv.ParseUint(cid, 10, 32); err == nil && id64 != 0 {
				vals["ClientID"] = uint(id64)
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
