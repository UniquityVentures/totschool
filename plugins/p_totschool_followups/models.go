package p_totschool_followups

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/registry"
	"github.com/UniquityVentures/totschool/plugins/p_totschool_clients"
	"gorm.io/gorm"
)

type Followup struct {
	gorm.Model
	CreatedByID     uint                       `gorm:"notnull"`
	CreatedBy       p_users.User               `gorm:"foreignKey:CreatedByID"`
	ClientID        uint                       `gorm:"notnull"`
	Client          p_totschool_clients.Client `gorm:"foreignKey:ClientID"`
	Title           string                     `gorm:"size:250;notnull"`
	ExtraInfo       string                     `gorm:"type:text"`
	GeneratedLetter string                     `gorm:"type:text"`
	GenerationID    *int                       // non-nil while AI generation is in progress
}

func pluginDBInitHooks() lamu.PluginFeatures[lamu.DBInitHook] {
	return lamu.PluginFeatures[lamu.DBInitHook]{
		Entries: []registry.Pair[string, lamu.DBInitHook]{
			{
				Key: "p_totschool_followups.bootstrap",
				Value: func(d *gorm.DB) *gorm.DB {
					d.Model(&Followup{}).Where("generation_id IS NOT NULL").Update("generation_id", nil)
					go runWorker(d)
					return d
				},
			},
		},
	}
}

func init() {
	lamu.RegistryAdmin.Register("p_totschool_followups.Followup", lamu.AdminPanel[Followup]{
		SearchField: "Title",
		Preload:     []string{"Client"},
	})
}
