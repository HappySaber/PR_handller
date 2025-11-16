package services

import (
	"PR/internal/controllers/dto"
	"PR/internal/database"
	"PR/internal/models"
	"database/sql"
	"log"

	"github.com/lib/pq"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (us *UserService) SetIsActive(userID string, isActive bool) (*dto.UserResponse, error) {
	query := `
			UPDATE users u
			SET is_active = $1
			FROM (
				SELECT id, name
				FROM teams
			) t
			WHERE u.id = $2
			RETURNING u.id,
					u.name,
					u.is_active,
					(SELECT name FROM teams WHERE id = u.team_id) AS team_name;

`
	var resp dto.UserResponse
	var teamName sql.NullString
	err := database.DB.QueryRow(query, isActive, userID).Scan(&resp.UserID, &resp.Username, &resp.IsActive, &teamName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	resp.TeamName = ""
	if teamName.Valid {
		resp.TeamName = teamName.String
	}
	return &resp, nil
}

func (us *UserService) GetReviews(userID string) ([]models.PullRequest, error) {
	query := "SELECT id, title, author_id, status, reviewer_ids, created_at, merged_at FROM pull_requests WHERE $1 = ANY(reviewer_ids)"
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("failed to close rows:", err)
		}
	}()

	var reviews []models.PullRequest
	var reviewers pq.StringArray
	for rows.Next() {
		var pr models.PullRequest
		if err := rows.Scan(
			&pr.PullRequestID,
			&pr.PullRequestName,
			&pr.AuthorID,
			&pr.Status,
			&reviewers,
			&pr.CreatedAt,
			&pr.MergedAt,
		); err != nil {
			return nil, err
		}
		pr.AssignedReviewers = []string(reviewers)
		reviews = append(reviews, pr)
	}

	return reviews, nil
}

type ReviewStat struct {
	ReviewerID  string `json:"reviewer_id"`
	Assignments int    `json:"assignments"`
}

func (us *UserService) GetReviewStats() ([]ReviewStat, error) {
	query := `
        SELECT reviewer_id, COUNT(*) AS assignments
        FROM pull_requests, unnest(reviewer_ids) AS reviewer_id
        GROUP BY reviewer_id
        ORDER BY assignments DESC
    `

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("failed to close rows:", err)
		}
	}()

	var stats []ReviewStat
	for rows.Next() {
		var s ReviewStat
		if err := rows.Scan(&s.ReviewerID, &s.Assignments); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, nil
}
