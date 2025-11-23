package usecases

import (
	context "context"
	domain "pr-manager-service/internal/domain"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
	isgomock struct{}
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// CreatePullRequest mocks base method.
func (m *MockStorage) CreatePullRequest(ctx context.Context, request domain.CreatePullRequestRequest, reviewersIDs []string) (domain.PullRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePullRequest", ctx, request, reviewersIDs)
	ret0, _ := ret[0].(domain.PullRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePullRequest indicates an expected call of CreatePullRequest.
func (mr *MockStorageMockRecorder) CreatePullRequest(ctx, request, reviewersIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePullRequest", reflect.TypeOf((*MockStorage)(nil).CreatePullRequest), ctx, request, reviewersIDs)
}

// CreateTeam mocks base method.
func (m *MockStorage) CreateTeam(ctx context.Context, request domain.CreateTeamRequest, teamID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTeam", ctx, request, teamID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTeam indicates an expected call of CreateTeam.
func (mr *MockStorageMockRecorder) CreateTeam(ctx, request, teamID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTeam", reflect.TypeOf((*MockStorage)(nil).CreateTeam), ctx, request, teamID)
}

// CreateUsers mocks base method.
func (m *MockStorage) CreateUsers(ctx context.Context, requests []domain.CreateUserRequest, teamID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUsers", ctx, requests, teamID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUsers indicates an expected call of CreateUsers.
func (mr *MockStorageMockRecorder) CreateUsers(ctx, requests, teamID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUsers", reflect.TypeOf((*MockStorage)(nil).CreateUsers), ctx, requests, teamID)
}

// GetActiveColleagues mocks base method.
func (m *MockStorage) GetActiveColleagues(ctx context.Context, userID string) ([]domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActiveColleagues", ctx, userID)
	ret0, _ := ret[0].([]domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActiveColleagues indicates an expected call of GetActiveColleagues.
func (mr *MockStorageMockRecorder) GetActiveColleagues(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActiveColleagues", reflect.TypeOf((*MockStorage)(nil).GetActiveColleagues), ctx, userID)
}

// GetPullRequestByID mocks base method.
func (m *MockStorage) GetPullRequestByID(ctx context.Context, prID string) (domain.PullRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPullRequestByID", ctx, prID)
	ret0, _ := ret[0].(domain.PullRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPullRequestByID indicates an expected call of GetPullRequestByID.
func (mr *MockStorageMockRecorder) GetPullRequestByID(ctx, prID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPullRequestByID", reflect.TypeOf((*MockStorage)(nil).GetPullRequestByID), ctx, prID)
}

// GetPullRequestsByReviewer mocks base method.
func (m *MockStorage) GetPullRequestsByReviewer(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPullRequestsByReviewer", ctx, userID)
	ret0, _ := ret[0].([]domain.PullRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPullRequestsByReviewer indicates an expected call of GetPullRequestsByReviewer.
func (mr *MockStorageMockRecorder) GetPullRequestsByReviewer(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPullRequestsByReviewer", reflect.TypeOf((*MockStorage)(nil).GetPullRequestsByReviewer), ctx, userID)
}

// GetPullRequestsStats mocks base method.
func (m *MockStorage) GetPullRequestsStats(ctx context.Context) ([]domain.PullRequestStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPullRequestsStats", ctx)
	ret0, _ := ret[0].([]domain.PullRequestStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPullRequestsStats indicates an expected call of GetPullRequestsStats.
func (mr *MockStorageMockRecorder) GetPullRequestsStats(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPullRequestsStats", reflect.TypeOf((*MockStorage)(nil).GetPullRequestsStats), ctx)
}

// GetTeamByName mocks base method.
func (m *MockStorage) GetTeamByName(ctx context.Context, teamName string) (domain.Team, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTeamByName", ctx, teamName)
	ret0, _ := ret[0].(domain.Team)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTeamByName indicates an expected call of GetTeamByName.
func (mr *MockStorageMockRecorder) GetTeamByName(ctx, teamName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTeamByName", reflect.TypeOf((*MockStorage)(nil).GetTeamByName), ctx, teamName)
}

// GetTeamFullByName mocks base method.
func (m *MockStorage) GetTeamFullByName(ctx context.Context, teamName string) (domain.Team, []domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTeamFullByName", ctx, teamName)
	ret0, _ := ret[0].(domain.Team)
	ret1, _ := ret[1].([]domain.User)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetTeamFullByName indicates an expected call of GetTeamFullByName.
func (mr *MockStorageMockRecorder) GetTeamFullByName(ctx, teamName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTeamFullByName", reflect.TypeOf((*MockStorage)(nil).GetTeamFullByName), ctx, teamName)
}

// GetUserFull mocks base method.
func (m *MockStorage) GetUserFull(ctx context.Context, userID string) (domain.User, domain.Team, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserFull", ctx, userID)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(domain.Team)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUserFull indicates an expected call of GetUserFull.
func (mr *MockStorageMockRecorder) GetUserFull(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserFull", reflect.TypeOf((*MockStorage)(nil).GetUserFull), ctx, userID)
}

// GetUserShort mocks base method.
func (m *MockStorage) GetUserShort(ctx context.Context, userID string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserShort", ctx, userID)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserShort indicates an expected call of GetUserShort.
func (mr *MockStorageMockRecorder) GetUserShort(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserShort", reflect.TypeOf((*MockStorage)(nil).GetUserShort), ctx, userID)
}

// GetUsersStats mocks base method.
func (m *MockStorage) GetUsersStats(ctx context.Context) ([]domain.UserStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersStats", ctx)
	ret0, _ := ret[0].([]domain.UserStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersStats indicates an expected call of GetUsersStats.
func (mr *MockStorageMockRecorder) GetUsersStats(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersStats", reflect.TypeOf((*MockStorage)(nil).GetUsersStats), ctx)
}

// PullRequestAssignmentsIncrement mocks base method.
func (m *MockStorage) PullRequestAssignmentsIncrement(ctx context.Context, pullRequestID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PullRequestAssignmentsIncrement", ctx, pullRequestID)
	ret0, _ := ret[0].(error)
	return ret0
}

// PullRequestAssignmentsIncrement indicates an expected call of PullRequestAssignmentsIncrement.
func (mr *MockStorageMockRecorder) PullRequestAssignmentsIncrement(ctx, pullRequestID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PullRequestAssignmentsIncrement", reflect.TypeOf((*MockStorage)(nil).PullRequestAssignmentsIncrement), ctx, pullRequestID)
}

// PullRequestStatsCreate mocks base method.
func (m *MockStorage) PullRequestStatsCreate(ctx context.Context, pullRequestID string, assignmentsCount int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PullRequestStatsCreate", ctx, pullRequestID, assignmentsCount)
	ret0, _ := ret[0].(error)
	return ret0
}

// PullRequestStatsCreate indicates an expected call of PullRequestStatsCreate.
func (mr *MockStorageMockRecorder) PullRequestStatsCreate(ctx, pullRequestID, assignmentsCount any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PullRequestStatsCreate", reflect.TypeOf((*MockStorage)(nil).PullRequestStatsCreate), ctx, pullRequestID, assignmentsCount)
}

// UnitOfWork mocks base method.
func (m *MockStorage) UnitOfWork(ctx context.Context, do func(Storage) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnitOfWork", ctx, do)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnitOfWork indicates an expected call of UnitOfWork.
func (mr *MockStorageMockRecorder) UnitOfWork(ctx, do any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnitOfWork", reflect.TypeOf((*MockStorage)(nil).UnitOfWork), ctx, do)
}

// UpdatePullRequestReviewersIDs mocks base method.
func (m *MockStorage) UpdatePullRequestReviewersIDs(ctx context.Context, prID string, reviewersIDs []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePullRequestReviewersIDs", ctx, prID, reviewersIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePullRequestReviewersIDs indicates an expected call of UpdatePullRequestReviewersIDs.
func (mr *MockStorageMockRecorder) UpdatePullRequestReviewersIDs(ctx, prID, reviewersIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePullRequestReviewersIDs", reflect.TypeOf((*MockStorage)(nil).UpdatePullRequestReviewersIDs), ctx, prID, reviewersIDs)
}

// UpdatePullRequestStatus mocks base method.
func (m *MockStorage) UpdatePullRequestStatus(ctx context.Context, prID string, newStatus domain.PullRequestStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePullRequestStatus", ctx, prID, newStatus)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePullRequestStatus indicates an expected call of UpdatePullRequestStatus.
func (mr *MockStorageMockRecorder) UpdatePullRequestStatus(ctx, prID, newStatus any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePullRequestStatus", reflect.TypeOf((*MockStorage)(nil).UpdatePullRequestStatus), ctx, prID, newStatus)
}

// UpdateUserStatus mocks base method.
func (m *MockStorage) UpdateUserStatus(ctx context.Context, userID string, isActive bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserStatus", ctx, userID, isActive)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserStatus indicates an expected call of UpdateUserStatus.
func (mr *MockStorageMockRecorder) UpdateUserStatus(ctx, userID, isActive any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserStatus", reflect.TypeOf((*MockStorage)(nil).UpdateUserStatus), ctx, userID, isActive)
}

// UserAssignmentsIncrementBatch mocks base method.
func (m *MockStorage) UserAssignmentsIncrementBatch(ctx context.Context, userIDs []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserAssignmentsIncrementBatch", ctx, userIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// UserAssignmentsIncrementBatch indicates an expected call of UserAssignmentsIncrementBatch.
func (mr *MockStorageMockRecorder) UserAssignmentsIncrementBatch(ctx, userIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserAssignmentsIncrementBatch", reflect.TypeOf((*MockStorage)(nil).UserAssignmentsIncrementBatch), ctx, userIDs)
}

// UserStatsCreateBatch mocks base method.
func (m *MockStorage) UserStatsCreateBatch(ctx context.Context, userIDs []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserStatsCreateBatch", ctx, userIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// UserStatsCreateBatch indicates an expected call of UserStatsCreateBatch.
func (mr *MockStorageMockRecorder) UserStatsCreateBatch(ctx, userIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserStatsCreateBatch", reflect.TypeOf((*MockStorage)(nil).UserStatsCreateBatch), ctx, userIDs)
}

// UserStatusChangesIncrementBatch mocks base method.
func (m *MockStorage) UserStatusChangesIncrementBatch(ctx context.Context, userIDs []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserStatusChangesIncrementBatch", ctx, userIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// UserStatusChangesIncrementBatch indicates an expected call of UserStatusChangesIncrementBatch.
func (mr *MockStorageMockRecorder) UserStatusChangesIncrementBatch(ctx, userIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserStatusChangesIncrementBatch", reflect.TypeOf((*MockStorage)(nil).UserStatusChangesIncrementBatch), ctx, userIDs)
}
