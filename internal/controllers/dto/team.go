package dto

type CreateTeamRequest struct {
	TeamName string          `json:"team_name" binding:"required"`
	Members  []TeamMemberDTO `json:"members" binding:"required"`
}

type TeamMemberDTO struct {
	UserID   string `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	IsActive bool   `json:"is_active" binding:"required"`
}

type TeamResponse struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDTO `json:"members"`
}
