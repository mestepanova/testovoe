package usecases

import (
	"context"
	"testing"
	"time"

	"pr-manager-service/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUsecases_CreatePullRequest(t *testing.T) {
	mockUnitOfWork := func(ms *MockStorage) {
		ms.EXPECT().UnitOfWork(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, fn func(s Storage) error) error {
				return fn(ms)
			})
	}

	timeNow := time.Now()

	const (
		prID       = "100"
		prName     = "prname1"
		prAuthorID = "200"

		teamID = "300"

		userID1   = "101"
		userName1 = "user1"

		userID2   = "102"
		userName2 = "user2"
	)

	testCases := []struct {
		name         string
		mock         func(*MockStorage)
		in           domain.CreatePullRequestRequest
		expect       domain.PullRequest
		expectErrMsg string
	}{
		{
			name: "happy_path",
			in: domain.CreatePullRequestRequest{
				ID:           prID,
				Name:         prName,
				AuthorUserID: prAuthorID,
			},
			expect: domain.PullRequest{
				ID:                prID,
				Name:              prName,
				AuthorUserID:      prAuthorID,
				ReviewersUsersIDs: []string{userID1, userID2},
				CreatedAt:         &timeNow,
				Status:            domain.StatusOpen,
			},
			mock: func(ms *MockStorage) {
				ms.EXPECT().
					GetUserShort(gomock.Any(), prAuthorID).
					Return(domain.User{
						ID:       prAuthorID,
						IsActive: true,
					}, nil)

				mockUnitOfWork(ms)

				ms.EXPECT().
					GetActiveColleagues(gomock.Any(), prAuthorID).
					Return(
						[]domain.User{
							{
								ID:       userID1,
								Name:     userName1,
								IsActive: true,
								TeamID:   teamID,
							},
							{
								ID:       userID2,
								Name:     userName2,
								IsActive: true,
								TeamID:   teamID,
							},
						},
						nil,
					)

				ms.EXPECT().
					CreatePullRequest(
						gomock.Any(),
						domain.CreatePullRequestRequest{
							ID:           prID,
							Name:         prName,
							AuthorUserID: prAuthorID,
						},
						gomock.InAnyOrder([]string{userID1, userID2}),
					).
					Return(
						domain.PullRequest{
							ID:                prID,
							Name:              prName,
							AuthorUserID:      prAuthorID,
							ReviewersUsersIDs: []string{userID1, userID2},
							CreatedAt:         &timeNow,
							Status:            domain.StatusOpen,
						},
						nil,
					)

				ms.EXPECT().
					UserAssignmentsIncrementBatch(
						gomock.Any(),
						gomock.InAnyOrder([]string{userID1, userID2}),
					).
					Return(nil)

				ms.EXPECT().
					PullRequestStatsCreate(gomock.Any(), prID, 2).
					Return(nil)
			},
		},
		{
			name: "one_active_collegue",
			in: domain.CreatePullRequestRequest{
				ID:           prID,
				Name:         prName,
				AuthorUserID: prAuthorID,
			},
			expect: domain.PullRequest{
				ID:                prID,
				Name:              prName,
				AuthorUserID:      prAuthorID,
				ReviewersUsersIDs: []string{userID1},
				CreatedAt:         &timeNow,
				Status:            domain.StatusOpen,
			},
			mock: func(ms *MockStorage) {
				ms.EXPECT().
					GetUserShort(gomock.Any(), prAuthorID).
					Return(domain.User{
						ID:       prAuthorID,
						IsActive: true,
					}, nil)

				mockUnitOfWork(ms)

				ms.EXPECT().
					GetActiveColleagues(gomock.Any(), prAuthorID).
					Return(
						[]domain.User{
							{
								ID:       userID1,
								Name:     userName1,
								IsActive: true,
								TeamID:   teamID,
							},
						},
						nil,
					)

				ms.EXPECT().
					CreatePullRequest(
						gomock.Any(),
						domain.CreatePullRequestRequest{
							ID:           prID,
							Name:         prName,
							AuthorUserID: prAuthorID,
						},
						[]string{userID1},
					).
					Return(
						domain.PullRequest{
							ID:                prID,
							Name:              prName,
							AuthorUserID:      prAuthorID,
							ReviewersUsersIDs: []string{userID1},
							CreatedAt:         &timeNow,
							Status:            domain.StatusOpen,
						},
						nil,
					)

				ms.EXPECT().
					UserAssignmentsIncrementBatch(
						gomock.Any(),
						[]string{userID1},
					).
					Return(nil)

				ms.EXPECT().
					PullRequestStatsCreate(gomock.Any(), prID, 1).
					Return(nil)
			},
		},
		{
			name: "zero_active_collegues",
			in: domain.CreatePullRequestRequest{
				ID:           prID,
				Name:         prName,
				AuthorUserID: prAuthorID,
			},
			expect: domain.PullRequest{
				ID:                prID,
				Name:              prName,
				AuthorUserID:      prAuthorID,
				ReviewersUsersIDs: []string{},
				CreatedAt:         &timeNow,
				Status:            domain.StatusOpen,
			},
			mock: func(ms *MockStorage) {
				ms.EXPECT().
					GetUserShort(gomock.Any(), prAuthorID).
					Return(domain.User{
						ID:       prAuthorID,
						IsActive: true,
					}, nil)

				mockUnitOfWork(ms)

				ms.EXPECT().
					GetActiveColleagues(gomock.Any(), prAuthorID).
					Return(
						[]domain.User{},
						nil,
					)

				ms.EXPECT().
					CreatePullRequest(
						gomock.Any(),
						domain.CreatePullRequestRequest{
							ID:           prID,
							Name:         prName,
							AuthorUserID: prAuthorID,
						},
						[]string{},
					).
					Return(
						domain.PullRequest{
							ID:                prID,
							Name:              prName,
							AuthorUserID:      prAuthorID,
							ReviewersUsersIDs: []string{},
							CreatedAt:         &timeNow,
							Status:            domain.StatusOpen,
						},
						nil,
					)

				ms.EXPECT().
					UserAssignmentsIncrementBatch(
						gomock.Any(),
						[]string{},
					).
					Return(nil)

				ms.EXPECT().
					PullRequestStatsCreate(gomock.Any(), prID, 0).
					Return(nil)
			},
		},
		{
			name: "author_not_found",
			in: domain.CreatePullRequestRequest{
				ID:           prID,
				Name:         prName,
				AuthorUserID: prAuthorID,
			},
			mock: func(ms *MockStorage) {
				ms.EXPECT().
					GetUserShort(gomock.Any(), prAuthorID).
					Return(domain.User{}, domain.ErrUserNotFound)
			},
			expectErrMsg: "storage.GetUserShort: user not found",
		},
		{
			name: "author_not_active",
			in: domain.CreatePullRequestRequest{
				ID:           prID,
				Name:         prName,
				AuthorUserID: prAuthorID,
			},
			mock: func(ms *MockStorage) {
				ms.EXPECT().
					GetUserShort(gomock.Any(), prAuthorID).
					Return(domain.User{
						ID:       prAuthorID,
						IsActive: false,
					}, nil)
			},
			expectErrMsg: domain.ErrUserInactive.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			storageMock := NewMockStorage(ctrl)
			tc.mock(storageMock)

			u := NewUsecases(storageMock)
			gotPullRequest, err := u.CreatePullRequest(context.Background(), tc.in)
			if tc.expectErrMsg != "" {
				require.Error(t, err)
				assert.Equal(t, tc.expectErrMsg, err.Error())
				assert.Equal(t, domain.PullRequest{}, gotPullRequest)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expect, gotPullRequest)
			}
		})
	}
}
