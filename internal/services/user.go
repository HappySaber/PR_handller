package services

import (
	"PR/internal/controllers/dto"
	"PR/internal/database"
	"PR/internal/models"
	"database/sql"

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
    FROM teams t
    WHERE u.id = $2 AND u.team_id = t.id
    RETURNING u.id, u.name, u.is_active, t.name
`
	var resp dto.UserResponse
	err := database.DB.QueryRow(query, isActive, userID).Scan(&resp.UserID, &resp.TeamName, &resp.IsActive, &resp.TeamName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &resp, nil
}

func (us *UserService) GetReviews(userID string) ([]models.PullRequest, error) {
	query := "SELECT id, title, author_id, status, reviewer_ids, created_at, merged_at FROM pull_requests WHERE $1 = ANY(reviewer_ids)"
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

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
