package p_totschool_followups

import (
	"context"
	"net/http"
	"strconv"

	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
)

func followupCreateSuccessURL(ctx context.Context) (string, error) {
	if r, ok := ctx.Value("$request").(*http.Request); ok && requestReturnClient(r) {
		if cid := r.URL.Query().Get("ClientID"); cid != "" {
			if id64, err := strconv.ParseUint(cid, 10, 32); err == nil && id64 != 0 {
				return lamu.RoutePath("clients.DetailRoute", map[string]getters.Getter[any]{
					"id": getters.Any(getters.Static(uint(id64))),
				})(ctx)
			}
		}
		if clientID, err := getters.Key[uint]("$in.ClientID")(ctx); err == nil && clientID != 0 {
			return lamu.RoutePath("clients.DetailRoute", map[string]getters.Getter[any]{
				"id": getters.Any(getters.Static(clientID)),
			})(ctx)
		}
	}
	return lamu.RoutePath("followups.DetailRoute", map[string]getters.Getter[any]{
		"id": getters.Any(getters.Key[uint]("$id")),
	})(ctx)
}

func followupUpdateSuccessURL(ctx context.Context) (string, error) {
	if r, ok := ctx.Value("$request").(*http.Request); ok && requestReturnClient(r) {
		return lamu.RoutePath("clients.DetailRoute", map[string]getters.Getter[any]{
			"id": getters.Any(getters.Key[uint]("followup.ClientID")),
		})(ctx)
	}
	return lamu.RoutePath("followups.DetailRoute", map[string]getters.Getter[any]{
		"id": getters.Any(getters.Key[uint]("followup.ID")),
	})(ctx)
}

func followupDeleteSuccessURL(ctx context.Context) (string, error) {
	if r, ok := ctx.Value("$request").(*http.Request); ok && requestReturnClient(r) {
		return lamu.RoutePath("clients.DetailRoute", map[string]getters.Getter[any]{
			"id": getters.Any(getters.Key[uint]("followup.ClientID")),
		})(ctx)
	}
	return lamu.RoutePath("followups.ListRoute", nil)(ctx)
}
