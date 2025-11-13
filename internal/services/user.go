package services

import (
	"PR/internal/database"
	"PR/internal/models"
)

type UserService struct {
	User models.User
}

func (us *UserService) SetIsActive() {
	us.User.IsActive = true
	//TODO написать sql

}

func (us *UserService) GetReviews() ([]models.PullRequest, error) {
	//TODO написать sql
	query := ""
	rows, err := database.DB.Query(query)
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
