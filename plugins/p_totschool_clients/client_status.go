package p_totschool_clients

import (
	"database/sql/driver"
	"fmt"

	"github.com/UniquityVentures/lamu/registry"
)

// ClientStatus mirrors the PostgreSQL enum client_status.
type ClientStatus string

const (
	ClientStatusActive   ClientStatus = "active"
	ClientStatusArchived ClientStatus = "archived"
)

var ClientStatusChoices = []registry.Pair[ClientStatus, string]{
	{Key: ClientStatusActive, Value: "Active"},
	{Key: ClientStatusArchived, Value: "Archived"},
}

func (s ClientStatus) Value() (driver.Value, error) {
	switch s {
	case ClientStatusActive, ClientStatusArchived:
		return string(s), nil
	default:
		return nil, fmt.Errorf("invalid ClientStatus: %q", s)
	}
}

func (s *ClientStatus) Scan(src any) error {
	if src == nil {
		return fmt.Errorf("ClientStatus: NULL")
	}
	var str string
	switch v := src.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return fmt.Errorf("ClientStatus: cannot scan %T", src)
	}
	switch ClientStatus(str) {
	case ClientStatusActive, ClientStatusArchived:
		*s = ClientStatus(str)
		return nil
	default:
		return fmt.Errorf("ClientStatus: unknown value %q", str)
	}
}
