package p_totschool_appointments

import (
	"database/sql/driver"
	"fmt"

	"github.com/UniquityVentures/lamu/registry"
)

// AppointmentStatus mirrors the PostgreSQL enum appointment_status.
type AppointmentStatus string

const (
	AppointmentStatusPending    AppointmentStatus = "pending"
	AppointmentStatusDone       AppointmentStatus = "done"
	AppointmentStatusCancelled  AppointmentStatus = "cancelled"
	AppointmentStatusPostponed  AppointmentStatus = "postponed"
)

var AppointmentStatusChoices = []registry.Pair[AppointmentStatus, string]{
	{Key: AppointmentStatusPending, Value: "Pending"},
	{Key: AppointmentStatusDone, Value: "Done"},
	{Key: AppointmentStatusCancelled, Value: "Cancelled"},
	{Key: AppointmentStatusPostponed, Value: "Postponed"},
}

func (s AppointmentStatus) Value() (driver.Value, error) {
	switch s {
	case AppointmentStatusPending, AppointmentStatusDone, AppointmentStatusCancelled, AppointmentStatusPostponed:
		return string(s), nil
	default:
		return nil, fmt.Errorf("invalid AppointmentStatus: %q", s)
	}
}

func (s *AppointmentStatus) Scan(src any) error {
	if src == nil {
		return fmt.Errorf("AppointmentStatus: NULL")
	}
	var str string
	switch v := src.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return fmt.Errorf("AppointmentStatus: cannot scan %T", src)
	}
	switch AppointmentStatus(str) {
	case AppointmentStatusPending, AppointmentStatusDone, AppointmentStatusCancelled, AppointmentStatusPostponed:
		*s = AppointmentStatus(str)
		return nil
	default:
		return fmt.Errorf("AppointmentStatus: unknown value %q", str)
	}
}
