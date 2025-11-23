package domain

type User struct {
	ID       string
	Name     string
	IsActive bool
	TeamID   string
}

type CreateUserRequest struct {
	ID       string `json:"user_id"   validate:"required,min=1,max=36"`
	Name     string `json:"username"  validate:"required,min=2,max=50"`
	IsActive bool   `json:"is_active" validate:"required"`
}

type UpdateUserStatusRequest struct {
	ID       string `json:"username"  validate:"required,min=1,max=36"`
	IsActive bool   `json:"is_active"`
}
