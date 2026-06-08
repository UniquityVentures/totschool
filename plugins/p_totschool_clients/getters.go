package p_totschool_clients

import (
	"context"
	"fmt"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
)

func dashboardTodayAppointmentsSummary() getters.Getter[string] {
	return getters.Map(
		getters.Key[components.ObjectList[DashboardAppointment]](dashboardTodayScheduleContextKey),
		func(_ context.Context, list components.ObjectList[DashboardAppointment]) (string, error) {
			switch list.Total {
			case 0:
				return "You have no Appointments Today", nil
			case 1:
				return "You have 1 Appointment Today", nil
			default:
				return fmt.Sprintf("You have %d Appointments Today", list.Total), nil
			}
		},
	)
}

func clientStatusRowClass() getters.Getter[string] {
	return getters.Map(getters.Key[ClientStatus]("$row.Status"), func(_ context.Context, status ClientStatus) (string, error) {
		switch status {
		case ClientStatusActive:
			return "'bg-success/10 hover:bg-success/20'", nil
		case ClientStatusArchived:
			return "'bg-base-200 opacity-60 hover:bg-base-300'", nil
		default:
			return "'hover:bg-base-200'", nil
		}
	})
}
