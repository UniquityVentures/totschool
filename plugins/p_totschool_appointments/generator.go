package p_totschool_appointments

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/registry"
	"github.com/UniquityVentures/totschool/plugins/p_totschool_clients"
	"gorm.io/gorm"
)

var AppointmentNames = []string{
	"Initial Consultation",
	"Follow-up Meeting",
	"Project Kickoff",
	"Strategy Session",
	"Quarterly Review",
	"Annual Check-in",
	"Onboarding Session",
	"Technical Interview",
	"Design Review",
	"Stakeholder Meeting",
	"Brainstorming Session",
	"Sprint Planning",
	"Retrospective",
	"Client Presentation",
	"Vendor Negotiation",
}

var Locations = []string{
	"Conference Room A (Headquarters)",
	"Conference Room B (Headquarters)",
	"Virtual (Zoom Link: https://zoom.us/j/123456789)",
	"Virtual (Google Meet: https://meet.google.com/abc-defg-hij)",
	"Client Office (Downtown)",
	"Client Office (Northside)",
	"Coffee Shop (Main St.)",
	"Co-working Space (Desk 12)",
	"Branch Office (West)",
	"Branch Office (East)",
}

func appointmentStatusForDatetime(dt time.Time) AppointmentStatus {
	if dt.Before(time.Now()) {
		return AppointmentStatusPending
	}
	if dt.After(time.Now()) {
		return AppointmentStatusDone
	}
	return AppointmentStatusPending
}

func GenerateAppointmentsForUser(db *gorm.DB, user p_users.User, count int) {
	now := time.Now()
	for range count {
		daysOffset := rand.Intn(30) + 1
		hoursOffset := rand.Intn(8) + 9
		minutesOffset := (rand.Intn(4)) * 15

		apptDate := time.Date(now.Year(), now.Month(), now.Day()+daysOffset, hoursOffset, minutesOffset, 0, 0, now.Location())

		for {
			overlappingCount, err := gorm.G[Appointment](db).Where(
				"created_by_id = ? AND datetime = ?",
				user.ID, apptDate,
			).Count(context.Background(), "*")
			if err != nil {
				overlappingCount = 0
			}

			if overlappingCount == 0 {
				break
			}
			apptDate = apptDate.Add(30 * time.Minute)
		}

		name := AppointmentNames[rand.Intn(len(AppointmentNames))]
		location := Locations[rand.Intn(len(Locations))]
		phone := fmt.Sprintf("(%03d) %03d-%04d", rand.Intn(800)+200, rand.Intn(900)+100, rand.Intn(10000))

		remarks := ""
		if rand.Float64() > 0.5 {
			remarks = "Please review the attached documents before the meeting."
		}

		extraInfo := ""
		if rand.Float64() > 0.7 {
			extraInfo = "Client prefers formal tone. Mention their recent project."
		}

		addr := location
		ph := phone
		clientRemarks := remarks
		client := p_totschool_clients.Client{
			CreatedByID: user.ID,
			Name:        name,
			Address:     &addr,
			Phone:       &ph,
			Remarks:     &clientRemarks,
		}
		if err := gorm.G[p_totschool_clients.Client](db).Create(context.Background(), &client); err != nil {
			continue
		}

		appointment := Appointment{
			CreatedByID: user.ID,
			ClientID:    client.ID,
			Datetime:    apptDate,
			Status:      appointmentStatusForDatetime(apptDate),
			Remarks:     remarks,
			ExtraInfo:   extraInfo,
		}
		_ = gorm.G[Appointment](db).Create(context.Background(), &appointment)
	}
}

func pluginGenerators() lamu.PluginFeatures[lamu.Generator] {
	return lamu.PluginFeatures[lamu.Generator]{
		Entries: []registry.Pair[string, lamu.Generator]{
			{
				Key: "appointments.Generator",
				Value: lamu.Generator{
					Create: func(db *gorm.DB) error {
						users, err := gorm.G[p_users.User](db).Find(context.Background())
						if err != nil {
							return err
						}

						for _, user := range users {
							count := 10 + rand.Intn(15)
							GenerateAppointmentsForUser(db, user, count)
						}
						return nil
					},
					Remove: func(db *gorm.DB) error {
						return db.Unscoped().Where("1=1").Delete(&Appointment{}).Error
					},
				},
			},
		},
	}
}
