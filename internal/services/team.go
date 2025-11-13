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

func (tm *TeamService) GetTeamMembers(Team *models.Team) ([]models.User, error) {
	//TODO написать sql query
	query := "SELECT"
	rows, err := database.DB.Query(query, Team.Name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	TeamMembers := make([]models.User, 0)
	for rows.Next() {
		var TeamMember models.User
		if err := rows.Scan(&TeamMember); err != nil {
			return nil, err
		}
		TeamMembers = append(TeamMembers, TeamMember)
	}
	return TeamMembers, nil
}
