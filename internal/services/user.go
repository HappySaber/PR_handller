package services

import (
	"PR/internal/database"
	"PR/internal/models"
)

type UserService struct {
	User models.User
}

func (us *UserService) SetIsActive(user *models.User) error {
	us.User.IsActive = true
	query := "UPDATE users SET is_active = TRUE WHERE id = $1"

	_, err := database.DB.Exec(query, user.Id)
	if err != nil {
		return err
	}
	user.IsActive = true
	return nil
}

func (us *UserService) GetReviews(user models.User) ([]models.PullRequest, error) {
	query := "SELECT * FROM pull_requests WHERE $1 = ANY(reviewer_ids)"
	rows, err := database.DB.Query(query, user.Id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var reviews []models.PullRequest
	for rows.Next() {
		var review models.PullRequest
		if err := rows.Scan(&review); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}
