package domain

type PullRequestStats struct {
	PullRequestID    string
	AssignmentsCount int64
}

type UserStats struct {
	UserID             string
	AssignmentsCount   int64
	StatusChangesCount int64
}
