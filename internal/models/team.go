package models

type Team struct {
	Name    string
	Members []TeamMembers
}

type TeamMembers struct {
	UserID   string
	Username string
	IsActive bool
}
