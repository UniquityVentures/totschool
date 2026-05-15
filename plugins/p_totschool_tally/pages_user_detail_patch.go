package p_totschool_tally

import (
	"context"
	"fmt"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/plugins/p_users"
)

// patchUserDetailForTally extends the users detail page with session environment and tally widgets.
func patchUserDetailForTally(page components.PageInterface) components.PageInterface {
	scaffold, ok := page.(*components.ShellScaffold)
	if !ok {
		panic("Base page for p_users.UserDetail was not ShellScaffold")
	}
	components.InsertChildAfter(scaffold,
		"p_users.UserDetailContent",
		func(*components.Detail[p_users.User]) components.ContainerColumn {
			return components.ContainerColumn{
				Children: []components.PageInterface{
					&components.Environment[uint]{
						Label:   "Session",
						Key:     getters.Static("session"),
						Options: SessionsListGetter,
						Default: tallySessionEnvironmentDefault,
					},
					TallySessionEntries{
						Page: components.Page{
							Key: "tally.UserSessionTallies",
						},
						UserGetter:    getters.Key[p_users.User]("user"),
						SessionGetter: CurrentEnvironmentSessionGetter,
					},
					StatLineChart{
						Page: components.Page{
							Key: "tally.UserSessionTalliesChart",
						},
						TalliesGetter: func(ctx context.Context) ([]Tally, error) {
							db, err := getters.DBFromContext(ctx)
							if err != nil {
								return nil, fmt.Errorf("StatLineChart: %w", err)
							}
							user, ok := ctx.Value("user").(p_users.User)
							if !ok {
								return nil, fmt.Errorf("StatLineChart: missing user in context")
							}
							session, err := CurrentEnvironmentSessionGetter(ctx)
							if err != nil {
								return nil, err
							}
							var tallies []Tally
							if err := db.
								Where("user_id = ? AND date >= ? AND date <= ?", user.ID, session.Start, session.End).
								Order("date ASC").
								Find(&tallies).Error; err != nil {
								return nil, err
							}
							return tallies, nil
						},
						Keys: []string{
							"Visits",
							"Appointments",
							"Leads",
							"Presentations",
							"Demos",
							"Letters",
							"FollowUps",
							"Proposals",
							"Policies",
							"Premium",
						},
					},
				},
			}
		},
	)
	return scaffold
}
