package services

import (
	"PR/internal/database"
	"PR/internal/models"
	"fmt"
)

type TeamService struct {
}

func (tm *TeamService) Create(team *models.Team) error {
	var teamID int
	query := "SELECT id FROM teams WHERE name=$1"
	err := database.DB.QueryRow(query, team.Name).Scan(&teamID)

	if err == nil {
		return fmt.Errorf("team exists")
	}

	query = "INSERT INTO teams (name) VALUES ($1) RETURNING id"
	err = database.DB.QueryRow(query).Scan(&teamID)
	if err != nil {
		return err
	}

	for _, user := range team.Members {
		query = "UPDATE users SET team_id = $1 WHERE id = $2"
		_, err := database.DB.Exec(user.UserID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (tm *TeamService) GetTeamMembers(teamName string) (*models.Team, error) {
	var teamID int
	query := "SELECT id FROM teams WHERE name = $1"
	if err := database.DB.QueryRow(query, teamName).Scan(&teamID); err != nil {
		return nil, fmt.Errorf("author not found: %w", err)
	}

	query = "SELECT id, name, is_active FROM users WHERE team_id = $1"
	rows, err := database.DB.Query(query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	team := &models.Team{
		Name:    teamName,
		Members: []models.TeamMember{},
	}
	for rows.Next() {
		var TeamMember models.TeamMember
		if err := rows.Scan(&TeamMember.UserID); err != nil {
			return nil, err
		}
		team.Members = append(team.Members, TeamMember)
	}
	return team, nil
}
