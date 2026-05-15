package p_totschool_proposals

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
	"gorm.io/gorm"
)

func pluginDBInitHooks() lamu.PluginFeatures[lamu.DBInitHook] {
	return lamu.PluginFeatures[lamu.DBInitHook]{
		Entries: []registry.Pair[string, lamu.DBInitHook]{
			{
				Key: "p_totschool_proposals.bootstrap",
				Value: func(d *gorm.DB) *gorm.DB {
					d.Model(&Proposal{}).Where("generation_id IS NOT NULL").Update("generation_id", nil)
					go runWorker(d)
					return d
				},
			},
		},
	}
}
