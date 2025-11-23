//go:build integration

package tests

import (
	"context"
	"sort"
	"testing"

	"pr-manager-service/internal/domain"
	"pr-manager-service/internal/generated/api"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPullRequestCreate(t *testing.T) {
	ctx := context.Background()

	t.Run("two_reviewers", func(t *testing.T) {
		cleanupDB(ctx, t)
		defer cleanupDB(ctx, t)

		const (
			teamName1 = "test name 1"

			userID1 = "100"
			userID2 = "101"
			userID3 = "102"
			userID4 = "103"

			username1 = "user1"
			username2 = "user2"
			username3 = "user3"
			username4 = "user4"

			prID   = "100"
			prName = "prname 1"
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
					IsActive: true,
				},
				{
					UserId:   userID4,
					Username: username4,
					IsActive: false,
				},
			},
		}
		activeReviewersIDs := lo.Map(apiTeam1.Members, func(m api.TeamMember, _ int) string {
			return m.UserId
		})

		// NOTE: добавление команды
		teamAddResp, err := client.PostTeamAdd(ctx, apiTeam1)
		require.NoError(t, err)
		require.Equal(t, 201, teamAddResp.StatusCode)

		// NOTE: создание пулреквеста
		pullRequestCreateResp, err := client.PostPullRequestCreateWithResponse(ctx, api.PostPullRequestCreateJSONRequestBody{
			AuthorId:        userID1,
			PullRequestId:   prID,
			PullRequestName: prName,
		})
		require.NoError(t, err)
		require.Equal(t, 201, pullRequestCreateResp.StatusCode())

		createdPrApi := pullRequestCreateResp.JSON201.Pr

		assert.Equal(t, userID1, createdPrApi.AuthorId)
		assert.Nil(t, createdPrApi.MergedAt)
		assert.NotNil(t, createdPrApi.CreatedAt)
		assert.Equal(t, prID, createdPrApi.PullRequestId)
		assert.Equal(t, prName, createdPrApi.PullRequestName)
		assert.Equal(t, string(api.OPEN), createdPrApi.Status)
		require.Equal(t, 2, len(createdPrApi.AssignedReviewers))
		assert.Subset(t, activeReviewersIDs, createdPrApi.AssignedReviewers)

		// NOTE: проверка, что в БД лежит правильный ПР
		domainPullRequest, err := testStorage.GetPullRequestByID(ctx, prID)
		require.NoError(t, err)

		assert.Equal(t, userID1, domainPullRequest.AuthorUserID)
		assert.Nil(t, domainPullRequest.MergedAt)
		assert.NotNil(t, domainPullRequest.CreatedAt)
		assert.Equal(t, prID, domainPullRequest.ID)
		assert.Equal(t, prName, domainPullRequest.Name)
		assert.Equal(t, domain.StatusOpen, domainPullRequest.Status)
		require.Equal(t, 2, len(domainPullRequest.ReviewersUsersIDs))
		assert.Subset(t, activeReviewersIDs, domainPullRequest.ReviewersUsersIDs)

		// NOTE: проверка статистики
		statsResp, err := client.GetStatsGetWithResponse(ctx)
		require.NoError(t, err)
		require.Equal(t, 200, statsResp.StatusCode())

		prStats := statsResp.JSON200.PullRequestsStats
		usersStats := statsResp.JSON200.UserStats

		require.Equal(t, 1, len(prStats))
		require.Equal(t, 4, len(usersStats))

		sort.Slice(usersStats, func(i, j int) bool {
			return usersStats[i].UserId < usersStats[j].UserId
		})

		assert.Equal(t, prID, prStats[0].PullRequestId)
		assert.Equal(t, 2, prStats[0].AssignmentsCount)

		usersWithOneAssignment := 0

		assert.Equal(t, userID4, usersStats[3].UserId)
		assert.Equal(t, 0, usersStats[3].StatusChangesCount)
		assert.True(t, usersStats[3].AssignmentsCount <= 1)
		if usersStats[3].AssignmentsCount == 1 {
			usersWithOneAssignment++
		}

		assert.Equal(t, userID3, usersStats[2].UserId)
		assert.Equal(t, 0, usersStats[2].StatusChangesCount)
		assert.True(t, usersStats[2].AssignmentsCount <= 1)
		if usersStats[2].AssignmentsCount == 1 {
			usersWithOneAssignment++
		}

		assert.Equal(t, userID2, usersStats[1].UserId)
		assert.Equal(t, 0, usersStats[1].StatusChangesCount)
		assert.True(t, usersStats[1].AssignmentsCount <= 1)
		if usersStats[1].AssignmentsCount == 1 {
			usersWithOneAssignment++
		}

		assert.Equal(t, userID1, usersStats[0].UserId)
		assert.Equal(t, 0, usersStats[0].StatusChangesCount)
		assert.Equal(t, 0, usersStats[0].AssignmentsCount)

		assert.Equal(t, 2, usersWithOneAssignment)
	})
}
