package services

import (
	"PR/internal/database"
	"PR/internal/models"
	"fmt"
)

type TestUserService struct{}

func NewTestUserService() *TestUserService {
	return &TestUserService{}

}

func (tus *TestUserService) CreateTestUsers(count int) ([]models.User, error) {
	users := make([]models.User, 0, 8)

	for i := 1; i <= 8; i++ {
		user_id := fmt.Sprintf("test_u%d", i)
		username := fmt.Sprintf("test%d", i)

		query := `INSERT INTO users (id, name) VALUES ($1, $2) RETURNING id`
		var id string
		err := database.DB.QueryRow(query, user_id, username).Scan(&id)
		if err != nil {
			return nil, err
		}

		user := models.User{
			UserID:   id,
			Username: username,
		}
		users = append(users, user)
	}

	return users, nil

}

func (tus *TestUserService) DeleteTestUsers() error {
	for i := 1; i <= 8; i++ {
		id := fmt.Sprintf("test_u%d", i)
		_, err := database.DB.Exec(`DELETE FROM pull_requests WHERE author_id = $1`, id)
		if err != nil {
			return fmt.Errorf("failed to delete PRs for %s: %w", id, err)
		}

		_, err = database.DB.Exec(`DELETE FROM users WHERE id = $1`, id)
		if err != nil {
			return fmt.Errorf("failed to delete user %s: %w", id, err)
		}
	}
	return nil
}
