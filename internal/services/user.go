package services

import (
	"PR/internal/database"
	"PR/internal/models"
	"errors"
)

type UserService struct {
	User models.User
}

func (us *UserService) SetIsActive(userID string, isActive bool) error {
	query := "UPDATE users SET is_active = TRUE WHERE id = $1"

	res, err := database.DB.Exec(query, isActive, userID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (us *UserService) GetReviews(userID string) ([]models.PullRequest, error) {
	query := "SELECT id, title, author_id, status, reviewer_ids, created_at, merged_at FROM pull_requests WHERE $1 = ANY(reviewer_ids)"
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var reviews []models.PullRequest
	for rows.Next() {
		var pr models.PullRequest
		if err := rows.Scan(
			&pr.PullRequestID,
			&pr.PullRequestName,
			&pr.AuthorID,
			&pr.Status,
			&pr.AssignedReviewers,
			&pr.CreatedAt,
			&pr.MergedAt,
		); err != nil {
			return nil, err
		}
		reviews = append(reviews, pr)
	}

	return reviews, nil
}
