package services

import (
	"PR/internal/database"
	"PR/internal/models"
)

type TeamService struct {
}

func (tm *TeamService) Create(Team *models.Team) error {
	//TODO написать sql query
	query := "INSERT INTO"
	_, err := database.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (tm *TeamService) GetTeamMembers(Team *models.Team) error {
	//TODO написать sql query
	query := "SELECT"
	rows, err := database.DB.Query(query, Team.Name)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var TeamMember models.User
		if err := rows.Scan(&TeamMember); err != nil {
			return err
		}
		Team.Members = append(Team.Members, TeamMember)
	}
	return nil
}
