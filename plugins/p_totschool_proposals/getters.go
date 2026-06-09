package p_totschool_proposals

import (
	"context"
	"log/slog"

	"github.com/UniquityVentures/lamu/getters"
)

func getterGenerated() getters.Getter[bool] {
	return func(ctx context.Context) (bool, error) {
		id, err := getters.Key[*int]("$in.GenerationID")(ctx)
		if err != nil {
			slog.Error("Error while getting id for checking if appointment is idle", "error", err)
			return false, err
		}
		content, err := getters.Key[string]("$in.GeneratedContent")(ctx)
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
		content, err := getters.Key[string]("$in.GeneratedContent")(ctx)
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
		content, err := getters.Key[string]("$in.GeneratedContent")(ctx)
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

func proposalFormClientID() getters.Getter[uint] {
	return func(ctx context.Context) (uint, error) {
		if id, err := getters.Key[uint]("$in.ClientID")(ctx); err == nil {
			return id, nil
		}
		ptr, err := getters.Key[*uint]("$in.ClientID")(ctx)
		if err != nil {
			return 0, err
		}
		if ptr == nil {
			return 0, nil
		}
		return *ptr, nil
	}
}

func getterProposalUnassigned() getters.Getter[bool] {
	return func(ctx context.Context) (bool, error) {
		clientID, err := proposalFormClientID()(ctx)
		if err != nil {
			return false, err
		}
		return clientID == 0, nil
	}
}

func getterProposalAssigned() getters.Getter[bool] {
	return func(ctx context.Context) (bool, error) {
		unassigned, err := getterProposalUnassigned()(ctx)
		if err != nil {
			return false, err
		}
		return !unassigned, nil
	}
}
