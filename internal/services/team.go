package services

import (
	"PR/internal/database"
	"PR/internal/models"
	"fmt"
)

type TeamService struct {
}

func (tm *TeamService) Create(Team *models.Team) error {
	var teamID int
	query := "INSERT INTO teams (name) VALUES ($1) RETURNING id"
	err := database.DB.QueryRow(query).Scan(&teamID)
	if err != nil {
		return err
	}

	for _, user := range Team.Members {
		query = "UPDATE users SET team_id = $1 WHERE id = $2"
		_, err := database.DB.Exec(query, teamID, user.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (tm *TeamService) GetTeamMembers(Team *models.Team) error {
	var teamID int
	query := "SELECT id FROM teams WHERE name = $1"
	if err := database.DB.QueryRow(query, Team.Name).Scan(&teamID); err != nil {
		return fmt.Errorf("author not found: %w", err)
	}

	query = "SELECT id, name, is_active FROM users WHERE team_id = $1"
	rows, err := database.DB.Query(query, teamID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var TeamMember models.User
		if err := rows.Scan(&TeamMember.Id); err != nil {
			return err
		}
		Team.Members = append(Team.Members, TeamMember)
	}
	return nil
}
