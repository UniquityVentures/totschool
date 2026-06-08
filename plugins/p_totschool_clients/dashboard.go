package p_totschool_clients

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/views"
	"gorm.io/gorm"
)

const dashboardTodayScheduleContextKey = "dashboardTodaySchedule"

// DashboardAppointment maps the appointments table for the clients dashboard
// without importing the appointments plugin (avoids a circular dependency).
type DashboardAppointment struct {
	gorm.Model
	CreatedByID uint   `gorm:"notnull"`
	ClientID    uint   `gorm:"notnull"`
	Client      Client `gorm:"foreignKey:ClientID"`
	Datetime    time.Time
	Location    string `gorm:"type:text"`
}

func (DashboardAppointment) TableName() string { return "appointments" }

type dashboardTodayScheduleLayer struct{}

func (dashboardTodayScheduleLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db, dberr := getters.DBFromContext(r.Context())
		if dberr != nil {
			slog.Error("dashboardTodayScheduleLayer: db from context", "error", dberr)
			next.ServeHTTP(w, r)
			return
		}

		now := time.Now()
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end := start.Add(24 * time.Hour)

		query := gorm.G[DashboardAppointment](db).
			Preload("Client", nil).
			Where("datetime >= ? AND datetime < ?", start, end).
			Order("datetime ASC").
			Order("id ASC")
		query = scopeDashboardAppointmentsToCurrentUser(r, query)

		rows, err := query.Find(r.Context())
		if err != nil {
			slog.Error("dashboardTodayScheduleLayer: query failed", "error", err)
			next.ServeHTTP(w, r)
			return
		}

		ol := components.ObjectList[DashboardAppointment]{
			Items:    rows,
			Number:   1,
			NumPages: 1,
			Total:    uint64(len(rows)),
		}
		ctx := context.WithValue(r.Context(), dashboardTodayScheduleContextKey, ol)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func scopeDashboardAppointmentsToCurrentUser(r *http.Request, query gorm.ChainInterface[DashboardAppointment]) gorm.ChainInterface[DashboardAppointment] {
	user, role := p_users.UserAndRoleFromContext(r.Context(), "scopeDashboardAppointmentsToCurrentUser")
	if user.IsSuperuser || role == "totschool_admin" {
		return query
	}
	return query.Where("created_by_id = ?", user.ID)
}
