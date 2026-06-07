package p_totschool_appointments

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"maragu.dev/gomponents"
)

func appointmentCreateSuccessURL(ctx context.Context) (string, error) {
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
	return lamu.RoutePath("appointments.DetailRoute", map[string]getters.Getter[any]{
		"id": getters.Any(getters.Key[uint]("$id")),
	})(ctx)
}

func appointmentUpdateSuccessURL(ctx context.Context) (string, error) {
	if r, ok := ctx.Value("$request").(*http.Request); ok && requestReturnClient(r) {
		return lamu.RoutePath("clients.DetailRoute", map[string]getters.Getter[any]{
			"id": getters.Any(getters.Key[uint]("appointment.ClientID")),
		})(ctx)
	}
	return lamu.RoutePath("appointments.DetailRoute", map[string]getters.Getter[any]{
		"id": getters.Any(getters.Key[uint]("appointment.ID")),
	})(ctx)
}

func appointmentTimelineDateGetter() getters.Getter[time.Time] {
	return getters.IfOrElse(getters.Key[time.Time]("$get.Date"), func(ctx context.Context) (time.Time, error) {
		return time.Now(), nil
	})
}

// appointmentTimelineDateFilterAttr is a boosted GET form that reloads on Date change (no filter dropdown).
func appointmentTimelineDateFilterAttr() getters.Getter[gomponents.Node] {
	route := lamu.RoutePath("appointments.CardTimelineRoute", nil)
	return func(ctx context.Context) (gomponents.Node, error) {
		boosted, err := getters.FormBoostedGet(route)(ctx)
		if err != nil {
			return nil, err
		}
		url, err := route(ctx)
		if err != nil {
			return nil, err
		}
		urlLit, err := json.Marshal(url)
		if err != nil {
			return nil, err
		}
		changeScript := fmt.Sprintf(
			`(function(evt){if(evt.target.name!=='Date')return;var f=evt.target.closest('form');if(!f)return;var m=f.closest('dialog.modal');var o={source:f,swap:'outerHTML',values:htmx.values(f),headers:{'HX-Boosted':'true'}};o.target=m||'body';htmx.ajax('GET',%s,o)})($event)`,
			urlLit,
		)
		return gomponents.Group{boosted, gomponents.Attr("@change", changeScript)}, nil
	}
}

func appointmentDeleteSuccessURL(ctx context.Context) (string, error) {
	if r, ok := ctx.Value("$request").(*http.Request); ok && requestReturnClient(r) {
		return lamu.RoutePath("clients.DetailRoute", map[string]getters.Getter[any]{
			"id": getters.Any(getters.Key[uint]("appointment.ClientID")),
		})(ctx)
	}
	return lamu.RoutePath("appointments.ListRoute", nil)(ctx)
}

func getterGenerated() getters.Getter[bool] {
	return func(ctx context.Context) (bool, error) {
		id, err := getters.Key[*int]("$in.GenerationID")(ctx)
		if err != nil {
			slog.Error("Error while getting id for checking if appointment is idle", "error", err)
			return false, err
		}
		content, err := getters.Key[string]("$in.GeneratedLetter")(ctx)
		if err != nil {
			slog.Error("Error while getting content for checking if appointment is idle", "error", err)
			return false, err
		}
		if id == nil && content != "" {
			return true, nil
		}
		return false, nil
	}
}

func getterGenerationPending() getters.Getter[bool] {
	return func(ctx context.Context) (bool, error) {
		id, err := getters.Key[*int]("$in.GenerationID")(ctx)
		if err != nil {
			slog.Error("Error while getting id for checking if appointment is idle", "error", err)
			return false, err
		}
		content, err := getters.Key[string]("$in.GeneratedLetter")(ctx)
		if err != nil {
			slog.Error("Error while getting content for checking if appointment is idle", "error", err)
			return false, err
		}
		if id != nil && content == "" {
			return true, nil
		}
		return false, nil
	}
}

func getterIdleGeneration() getters.Getter[bool] {
	return func(ctx context.Context) (bool, error) {
		id, err := getters.Key[*int]("$in.GenerationID")(ctx)
		if err != nil {
			slog.Error("Error while getting id for checking if appointment is idle", "error", err)
			return false, err
		}
		content, err := getters.Key[string]("$in.GeneratedLetter")(ctx)
		if err != nil {
			slog.Error("Error while getting content for checking if appointment is idle", "error", err)
			return false, err
		}
		if id == nil && content == "" {
			return true, nil
		}
		return false, nil
	}
}

func overlapAppointmentLinkLabel() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		name, err := getters.Key[string]("$row.Name")(ctx)
		if err != nil {
			return "", err
		}
		t, err := getters.Key[time.Time]("$row.Date")(ctx)
		if err != nil {
			return "", err
		}
		timezone, _ := ctx.Value("$tz").(*time.Location)
		if timezone == nil {
			timezone = components.DefaultTimeZone
		}
		dateStr := ""
		if !t.IsZero() {
			dateStr = t.In(timezone).Format("Mon, 02 Jan 2006 15:04:05")
		}
		return name + " — " + dateStr, nil
	}
}
