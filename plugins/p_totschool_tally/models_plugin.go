package p_totschool_tally

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginModels() lamu.PluginFeatures[any] {
	return lamu.PluginFeatures[any]{
		Entries: []registry.Pair[string, any]{
			{Key: "p_totschool_tally.TotSchoolSession", Value: TotSchoolSession{}},
			{Key: "p_totschool_tally.Tally", Value: Tally{}},
		},
	}
}
