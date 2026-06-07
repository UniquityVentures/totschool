package p_totschool_clients

import (
	"context"

	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/registry"
)

func clientStatusSelectGetter(ctxKey string) getters.Getter[registry.Pair[ClientStatus, string]] {
	return func(ctx context.Context) (registry.Pair[ClientStatus, string], error) {
		status, err := getters.Key[ClientStatus](ctxKey)(ctx)
		if err != nil {
			return registry.Pair[ClientStatus, string]{}, err
		}
		if p, ok := registry.PairFromPairs(status, ClientStatusChoices); ok {
			return p, nil
		}
		return registry.Pair[ClientStatus, string]{Key: status, Value: string(status)}, nil
	}
}

func clientStatusLabelFromRow() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		status, err := getters.Key[ClientStatus]("$row.Status")(ctx)
		if err != nil {
			return "", err
		}
		if p, ok := registry.PairFromPairs(status, ClientStatusChoices); ok {
			return p.Value, nil
		}
		return string(status), nil
	}
}

func clientStatusLabelFromIn() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		status, err := getters.Key[ClientStatus]("$in.Status")(ctx)
		if err != nil {
			return "", err
		}
		if p, ok := registry.PairFromPairs(status, ClientStatusChoices); ok {
			return p.Value, nil
		}
		return string(status), nil
	}
}
