package domain

type Team struct {
	ID   string
	Name string
}

type CreateTeamRequest struct {
	Name    string              `json:"team_name" validate:"required,min=2,max=50"`
	Members []CreateUserRequest `json:"members"   validate:"required,min=2,max=50"`
}

type DeactivateUsersRequest struct {
	TeamName string   `json:"team_name" validate:"required,min=2,max=50"`
	UserIDs  []string `json:"user_ids"  validate:"required,min=2,max=50,dive,min=1,max=36"`
}
