//go:build integration

package tests

import (
	"context"
	"slices"
	"sort"
	"testing"

	"pr-manager-service/internal/domain"
	"pr-manager-service/internal/generated/api"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTeam(t *testing.T) {
	ctx := context.Background()

	t.Run("found", func(t *testing.T) {
		cleanupDB(ctx, t)
		defer cleanupDB(ctx, t)

		const (
			teamName  = "test name"
			userID1   = "100"
			userID2   = "101"
			userID3   = "102"
			username1 = "user1"
			username2 = "user2"
			username3 = "user3"
		)

		apiTeam := api.Team{
			TeamName: teamName,
			Members: []api.TeamMember{
				{
					UserId:   userID1,
					Username: username1,
					IsActive: true,
				},
				{
					UserId:   userID2,
					Username: username2,
					IsActive: true,
				},
				{
					UserId:   userID3,
					Username: username3,
					IsActive: false,
				},
			},
		}

		teamAddResp, err := client.PostTeamAddWithResponse(ctx, apiTeam)
		require.NoError(t, err)

		require.Equal(t, 201, teamAddResp.StatusCode())
		require.NotNil(t, teamAddResp.JSON201)
		require.NotNil(t, teamAddResp.JSON201.Team)
		require.Nil(t, teamAddResp.JSON400)
		assert.Equal(t, teamAddResp.JSON201.Team, apiTeam)

		getTeamResp, err := client.GetTeamGetWithResponse(ctx, &api.GetTeamGetParams{
			TeamName: teamName,
		})
		require.NoError(t, err)
		require.Equal(t, 200, getTeamResp.StatusCode())
		require.NotNil(t, getTeamResp.JSON200)
		require.Nil(t, getTeamResp.JSON404)
		assert.Equal(t, *getTeamResp.JSON200, apiTeam)
	})

	t.Run("not_found", func(t *testing.T) {
		cleanupDB(ctx, t)
		defer cleanupDB(ctx, t)

		const (
			teamName  = "test name"
			userID1   = "100"
			userID2   = "101"
			userID3   = "102"
			username1 = "user1"
			username2 = "user2"
			username3 = "user3"
		)

		apiTeam := api.Team{
			TeamName: teamName,
			Members: []api.TeamMember{
				{
					UserId:   userID1,
					Username: username1,
					IsActive: true,
				},
				{
					UserId:   userID2,
					Username: username2,
					IsActive: true,
				},
				{
					UserId:   userID3,
					Username: username3,
					IsActive: false,
				},
			},
		}

		teamAddResp, err := client.PostTeamAddWithResponse(ctx, apiTeam)
		require.NoError(t, err)

		require.Equal(t, 201, teamAddResp.StatusCode())
		require.NotNil(t, teamAddResp.JSON201)
		require.NotNil(t, teamAddResp.JSON201.Team)
		require.Nil(t, teamAddResp.JSON400)
		assert.Equal(t, teamAddResp.JSON201.Team, apiTeam)

		expectErr := api.ErrorResponse{
			Error: struct {
				Code    api.ErrorCode `json:"code"`
				Message string        `json:"message"`
			}{
				Code:    api.NOTFOUND,
				Message: "team not found",
			},
		}

		getTeamResp, err := client.GetTeamGetWithResponse(ctx, &api.GetTeamGetParams{
			TeamName: teamName + "asdfasdf",
		})
		require.NoError(t, err)
		require.Equal(t, 404, getTeamResp.StatusCode())
		require.Nil(t, getTeamResp.JSON200)
		require.NotNil(t, getTeamResp.JSON404)
		assert.Equal(t, *getTeamResp.JSON404, expectErr)
	})
}

func TestCreateTeam(t *testing.T) {
	ctx := context.Background()

	t.Run("create_one", func(t *testing.T) {
		cleanupDB(ctx, t)
		defer cleanupDB(ctx, t)

		const (
			teamName  = "test name"
			userID1   = "100"
			userID2   = "101"
			userID3   = "102"
			username1 = "user1"
			username2 = "user2"
			username3 = "user3"
		)

		apiTeam := api.Team{
			TeamName: teamName,
			Members: []api.TeamMember{
				{
					UserId:   userID1,
					Username: username1,
					IsActive: true,
				},
				{
					UserId:   userID2,
					Username: username2,
					IsActive: true,
				},
				{
					UserId:   userID3,
					Username: username3,
					IsActive: false,
				},
			},
		}

		teamAddResp, err := client.PostTeamAddWithResponse(ctx, apiTeam)
		require.NoError(t, err)

		require.Equal(t, 201, teamAddResp.StatusCode())
		require.NotNil(t, teamAddResp.JSON201)
		require.NotNil(t, teamAddResp.JSON201.Team)
		require.Nil(t, teamAddResp.JSON400)
		assert.Equal(t, teamAddResp.JSON201.Team, apiTeam)

		// NOTE: проверка, что в БД сохранилась правильная команда
		domainTeamFromDB, domainUsersFromDB, err := testStorage.GetTeamFullByName(ctx, teamName)
		require.NoError(t, err)

		assert.Equal(t, apiTeam.TeamName, domainTeamFromDB.Name)

		expectDomainUsers := []domain.User{
			{
				ID:       userID1,
				Name:     username1,
				IsActive: true,
				TeamID:   domainTeamFromDB.ID,
			},
			{
				ID:       userID2,
				Name:     username2,
				IsActive: true,
				TeamID:   domainTeamFromDB.ID,
			},
			{
				ID:       userID3,
				Name:     username3,
				IsActive: false,
				TeamID:   domainTeamFromDB.ID,
			},
		}

		sortDomainUsers(expectDomainUsers)
		sortDomainUsers(domainUsersFromDB)

		assert.Equal(t, expectDomainUsers, domainUsersFromDB)

		// NOTE: проверка, что статистика посчиталась правильно
		statsResp, err := client.GetStatsGetWithResponse(ctx)
		require.NoError(t, err)
		require.NotNil(t, statsResp.JSON200)

		require.Equal(t, 3, len(statsResp.JSON200.UserStats))
		assert.Equal(t, 0, len(statsResp.JSON200.PullRequestsStats))

		userStats := statsResp.JSON200.UserStats
		sort.Slice(userStats, func(i, j int) bool {
			return userStats[i].UserId > userStats[j].UserId
		})

		for i, userStat := range userStats {
			assert.Equal(t, 0, userStat.AssignmentsCount)
			assert.Equal(t, 0, userStat.StatusChangesCount)
			assert.Equal(t, expectDomainUsers[i].ID, userStat.UserId)
		}
	})

	t.Run("create_second,_update_user", func(t *testing.T) {
		cleanupDB(ctx, t)
		defer cleanupDB(ctx, t)

		const (
			teamName1 = "test name 1"
			teamName2 = "test name 2"
			userID1   = "100"
			userID2   = "101"
			userID3   = "102"
			username1 = "user1"
			username2 = "user2"
			username3 = "user3"
			username4 = "user4"
		)

		apiTeam1 := api.Team{
			TeamName: teamName1,
			Members: []api.TeamMember{
				{
					UserId:   userID1,
					Username: username1,
					IsActive: true,
				},
				{
					UserId:   userID2,
					Username: username2,
					IsActive: true,
				},
				{
					UserId:   userID3,
					Username: username3,
					IsActive: false,
				},
			},
		}

		teamAddResp1, err := client.PostTeamAddWithResponse(ctx, apiTeam1)
		require.NoError(t, err)
		require.Equal(t, 201, teamAddResp1.StatusCode())

		apiTeam2 := api.Team{
			TeamName: teamName2,
			Members: []api.TeamMember{
				{
					UserId:   userID1,
					Username: username1,
					IsActive: false,
				},
				{
					UserId:   userID3,
					Username: username4,
					IsActive: true,
				},
			},
		}

		teamAddResp2, err := client.PostTeamAddWithResponse(ctx, apiTeam2)
		require.NoError(t, err)
		require.Equal(t, 201, teamAddResp2.StatusCode())
		require.NotNil(t, teamAddResp2.JSON201)
		assert.Equal(t, apiTeam2, teamAddResp2.JSON201.Team)

		// NOTE: проверка состояния в БД
		teamDomainFromDB1, usersDomainFromDB1, err := testStorage.GetTeamFullByName(ctx, teamName1)
		require.NoError(t, err)
		assert.Equal(t, teamDomainFromDB1.Name, teamName1)

		teamDomainFromDB2, usersDomainFromDB2, err := testStorage.GetTeamFullByName(ctx, teamName2)
		require.NoError(t, err)
		assert.Equal(t, teamDomainFromDB2.Name, teamName2)

		expectDomainUsers1 := []domain.User{
			{
				ID:       userID2,
				Name:     username2,
				IsActive: true,
				TeamID:   teamDomainFromDB1.ID,
			},
		}

		expectDomainUsers2 := []domain.User{
			{
				ID:       userID1,
				Name:     username1,
				IsActive: false,
				TeamID:   teamDomainFromDB2.ID,
			},
			{
				ID:       userID3,
				Name:     username4,
				IsActive: true,
				TeamID:   teamDomainFromDB2.ID,
			},
		}

		sortDomainUsers(usersDomainFromDB1)
		sortDomainUsers(usersDomainFromDB2)
		sortDomainUsers(expectDomainUsers1)
		sortDomainUsers(expectDomainUsers2)

		assert.Equal(t, expectDomainUsers1, usersDomainFromDB1)
		assert.Equal(t, expectDomainUsers2, usersDomainFromDB2)

		// NOTE: проверка статистики
		statsResp, err := client.GetStatsGetWithResponse(ctx)
		require.NoError(t, err)
		require.NotNil(t, statsResp.JSON200)

		require.Equal(t, 3, len(statsResp.JSON200.UserStats))
		assert.Equal(t, 0, len(statsResp.JSON200.PullRequestsStats))

		userStats := statsResp.JSON200.UserStats
		sort.Slice(userStats, func(i, j int) bool {
			return userStats[i].UserId > userStats[j].UserId
		})

		for _, userStat := range userStats {
			assert.Equal(t, 0, userStat.AssignmentsCount)
			assert.Equal(t, 0, userStat.StatusChangesCount)

			isInTeam1 := slices.ContainsFunc(expectDomainUsers1, func(u domain.User) bool {
				return u.ID == userStat.UserId
			})

			isInTeam2 := slices.ContainsFunc(expectDomainUsers2, func(u domain.User) bool {
				return u.ID == userStat.UserId
			})

			assert.True(t, isInTeam1 || isInTeam2)
		}
	})

	t.Run("create_second,_already_exists", func(t *testing.T) {
		cleanupDB(ctx, t)
		defer cleanupDB(ctx, t)

		const (
			teamName1 = "test name 1"
			userID1   = "100"
			userID2   = "101"
			userID3   = "102"
			username1 = "user1"
			username2 = "user2"
			username3 = "user3"
			username4 = "user4"
		)

		apiTeam1 := api.Team{
			TeamName: teamName1,
			Members: []api.TeamMember{
				{
					UserId:   userID1,
					Username: username1,
					IsActive: true,
				},
				{
					UserId:   userID2,
					Username: username2,
					IsActive: true,
				},
				{
					UserId:   userID3,
					Username: username3,
					IsActive: false,
				},
			},
		}

		teamAddResp1, err := client.PostTeamAddWithResponse(ctx, apiTeam1)
		require.NoError(t, err)
		require.Equal(t, 201, teamAddResp1.StatusCode())

		apiTeam2 := api.Team{
			TeamName: teamName1,
			Members: []api.TeamMember{
				{
					UserId:   userID1,
					Username: username1,
					IsActive: false,
				},
				{
					UserId:   userID3,
					Username: username4,
					IsActive: true,
				},
			},
		}

		expectErr := api.ErrorResponse{
			Error: struct {
				Code    api.ErrorCode `json:"code"`
				Message string        `json:"message"`
			}{
				Code:    api.TEAMEXISTS,
				Message: "team_name already exists",
			},
		}

		teamAddResp2, err := client.PostTeamAddWithResponse(ctx, apiTeam2)
		require.Equal(t, 400, teamAddResp2.StatusCode())
		require.Nil(t, teamAddResp2.JSON201)
		require.NotNil(t, teamAddResp2.JSON400)
		assert.Equal(t, expectErr, *teamAddResp2.JSON400)

		// NOTE: проверка состояния в БД
		teamDomainFromDB1, usersDomainFromDB1, err := testStorage.GetTeamFullByName(ctx, teamName1)
		require.NoError(t, err)
		assert.Equal(t, teamDomainFromDB1.Name, teamName1)

		expectDomainUsers1 := []domain.User{
			{
				ID:       userID1,
				Name:     username1,
				IsActive: true,
				TeamID:   teamDomainFromDB1.ID,
			},
			{
				ID:       userID2,
				Name:     username2,
				IsActive: true,
				TeamID:   teamDomainFromDB1.ID,
			},
			{
				ID:       userID3,
				Name:     username3,
				IsActive: false,
				TeamID:   teamDomainFromDB1.ID,
			},
		}

		sortDomainUsers(usersDomainFromDB1)
		sortDomainUsers(expectDomainUsers1)
		assert.Equal(t, expectDomainUsers1, usersDomainFromDB1)

		// NOTE: проверка статистики
		statsResp, err := client.GetStatsGetWithResponse(ctx)
		require.NoError(t, err)
		require.NotNil(t, statsResp.JSON200)

		require.Equal(t, 3, len(statsResp.JSON200.UserStats))
		assert.Equal(t, 0, len(statsResp.JSON200.PullRequestsStats))

		userStats := statsResp.JSON200.UserStats
		sort.Slice(userStats, func(i, j int) bool {
			return userStats[i].UserId > userStats[j].UserId
		})

		for i, userStat := range userStats {
			assert.Equal(t, 0, userStat.AssignmentsCount)
			assert.Equal(t, 0, userStat.StatusChangesCount)
			assert.Equal(t, expectDomainUsers1[i].ID, userStat.UserId)
		}
	})
}

func sortDomainUsers(u []domain.User) {
	sort.Slice(u, func(i, j int) bool {
		return u[i].ID > u[j].ID
	})
}
