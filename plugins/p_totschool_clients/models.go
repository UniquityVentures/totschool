package p_totschool_clients

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	CreatedByID uint         `gorm:"notnull"`
	CreatedBy   p_users.User `gorm:"foreignKey:CreatedByID"`
	Name        string       `gorm:"size:250;notnull"`
	Status      ClientStatus `gorm:"type:client_status;notnull;default:active"`
	Address     *string      `gorm:"type:text"`
	Phone       *string      `gorm:"size:20"`
	Remarks     *string      `gorm:"type:text"`
}

func init() {
	lamu.RegistryAdmin.Register("p_totschool_clients", lamu.AdminPanel[Client]{SearchField: "Name"})
}
