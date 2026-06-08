package p_totschool_followups

import (
	"context"
	"log/slog"

	"github.com/UniquityVentures/lamu/getters"
)

func getterGenerated() getters.Getter[bool] {
	return func(ctx context.Context) (bool, error) {
		id, err := getters.Key[*int]("$in.GenerationID")(ctx)
		if err != nil {
			slog.Error("Error while getting id for checking if followup is generated", "error", err)
			return false, err
		}
		content, err := getters.Key[string]("$in.GeneratedLetter")(ctx)
		if err != nil {
			slog.Error("Error while getting content for checking if followup is generated", "error", err)
			return false, err
		}
		return id == nil && content != "", nil
	}
}

func getterGenerationPending() getters.Getter[bool] {
	return func(ctx context.Context) (bool, error) {
		id, err := getters.Key[*int]("$in.GenerationID")(ctx)
		if err != nil {
			slog.Error("Error while getting id for checking if followup generation is pending", "error", err)
			return false, err
		}
		content, err := getters.Key[string]("$in.GeneratedLetter")(ctx)
		if err != nil {
			slog.Error("Error while getting content for checking if followup generation is pending", "error", err)
			return false, err
		}
		return id != nil && content == "", nil
	}
}

func getterIdleGeneration() getters.Getter[bool] {
	return func(ctx context.Context) (bool, error) {
		id, err := getters.Key[*int]("$in.GenerationID")(ctx)
		if err != nil {
			slog.Error("Error while getting id for checking if followup generation is idle", "error", err)
			return false, err
		}
		content, err := getters.Key[string]("$in.GeneratedLetter")(ctx)
		if err != nil {
			slog.Error("Error while getting content for checking if followup generation is idle", "error", err)
			return false, err
		}
		return id == nil && content == "", nil
	}
}
