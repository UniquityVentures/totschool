package p_totschool_clients

import (
	"context"

	"github.com/UniquityVentures/lamu/getters"
)

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
