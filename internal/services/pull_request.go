package services

import (
	"PR/internal/database"
	"PR/internal/models"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/lib/pq"
)

type PullRequestService struct {
	PullRequest models.PullRequest
}

func (prs *PullRequestService) Create(PR *models.PullRequest) error {
	var teamName string
	query := "SELECT t.name FROM teams t LEFT JOIN users u ON  t.id = u.team_id WHERE u.id = $1"
	if err := database.DB.QueryRow(query, PR.AuthorId).Scan(&teamName); err != nil {
		return fmt.Errorf("author not found: %w", err)
	}

	var TeamMembers []models.User
	query = "SELECT u.id, u.name, u.is_active FROM users u LEFT JOIN teams t ON u.team_id = t.id WHERE t.name = $1"
	rows, err := database.DB.Query(query, teamName)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.Id, &u.Name, &u.IsActive); err != nil {
			return err
		}
		TeamMembers = append(TeamMembers, u)
	}

	reviewers, err := prs.сhooseReviewers(TeamMembers)
	if err != nil {
		return err
	}

	reviewerIDs := []int{}
	for _, r := range reviewers {
		reviewerIDs = append(reviewerIDs, r.Id)
	}

	query = "INSERT INTO pull_requests (title, author_id, status, reviewer_ids) VALUES ($1, $2, $3, $4) RETURNING id, created_at"

	err = database.DB.QueryRow(
		query,
		PR.Title,
		PR.AuthorId,
		"OPEN",
		pq.Array(reviewerIDs),
	).Scan(&PR.Id, &PR.CreatedAt)
	if err != nil {
		return err
	}

	PR.Status = "OPEN"
	PR.Reviewers = reviewers
	return nil
}

func (prs *PullRequestService) Merge(PR *models.PullRequest) error {
	query := "UPDATE pull_request SET status = 'MERGED', merged_at = NOW() WHERE id = $1 AND merged_at <> 'MERGED' RETURNING status, merged_at"
	err := database.DB.QueryRow(query, PR.Id).Scan(&PR.Status, &PR.MergedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}

func (prs *PullRequestService) Reassign() error {
	query := "SELECT"
	rows, err := database.DB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (prs *PullRequestService) сhooseReviewers(users []models.User) ([]models.User, error) {
	if len(users) <= 2 {
		return users, nil
	}

	rand.Seed(time.Now().UnixNano())

	copyArr := append([]models.User(nil), users...)

	rand.Shuffle(len(copyArr), func(i, j int) {
		copyArr[i], copyArr[j] = copyArr[j], copyArr[i]
	})

	return copyArr[:2], nil
}
