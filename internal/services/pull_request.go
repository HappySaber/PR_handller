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

func NewPullRequestService() *PullRequestService {
	return &PullRequestService{}

}

func (prs *PullRequestService) Create(PR *models.PullRequest) error {
	var teamName string
	query := "SELECT t.name FROM teams t LEFT JOIN users u ON  t.id = u.team_id WHERE u.id = $1"
	if err := database.DB.QueryRow(query, PR.AuthorID).Scan(&teamName); err != nil {
		return fmt.Errorf("author not found: %w", err)
	}

	var TeamMembers []models.TeamMember
	query = "SELECT u.id, u.name, u.is_active FROM users u LEFT JOIN teams t ON u.team_id = t.id WHERE t.name = $1 AND u.id <> $2 AND u.is_active = true"
	rows, err := database.DB.Query(query, teamName, PR.AuthorID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var u models.TeamMember
		if err := rows.Scan(&u.UserID, &u.Username, &u.IsActive); err != nil {
			return err
		}
		TeamMembers = append(TeamMembers, u)
	}

	reviewers, err := prs.сhooseReviewers(TeamMembers)
	if err != nil {
		return err
	}

	reviewerIDs := []string{}
	for _, r := range reviewers {
		reviewerIDs = append(reviewerIDs, r.UserID)
	}

	query = "INSERT INTO pull_requests (id, title, author_id, status, reviewer_ids) VALUES ($1, $2, $3, $4, $5) RETURNING created_at"

	err = database.DB.QueryRow(
		query,
		PR.PullRequestID,
		PR.PullRequestName,
		PR.AuthorID,
		"OPEN",
		pq.Array(reviewerIDs),
	).Scan(&PR.CreatedAt)
	if err != nil {
		return err
	}

	PR.Status = "OPEN"
	assigned := make([]string, 0, len(reviewers))
	for _, r := range reviewers {
		assigned = append(assigned, r.UserID)
	}
	PR.AssignedReviewers = assigned
	return nil
}

func (prs *PullRequestService) Merge(pullRequestID string) (*models.PullRequest, error) {

	PR := &models.PullRequest{}
	var reviewers pq.StringArray
	query := `UPDATE pull_requests 
			SET status = 'MERGED', merged_at = NOW() 
			WHERE id = $1 AND status <> 'MERGED' 
			RETURNING id, title, author_id, status, reviewer_ids, created_at, merged_at`
	err := database.DB.QueryRow(query, pullRequestID).Scan(&PR.PullRequestID, &PR.PullRequestName, &PR.AuthorID, &PR.Status, &reviewers, &PR.CreatedAt, &PR.MergedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return PR, nil
}

func (prs *PullRequestService) Reassign(prID, oldUserID string) (*models.PullRequest, string, error) {
	PR := &models.PullRequest{}
	query := `SELECT id, title, author_id, status, reviewer_ids, created_at, merged_at
			  FROM pull_requests
			  WHERE id = $1`
	var reviewers pq.StringArray
	err := database.DB.QueryRow(query, prID).Scan(
		&PR.PullRequestID,
		&PR.PullRequestName,
		&PR.AuthorID,
		&PR.Status,
		&reviewers,
		&PR.CreatedAt,
		&PR.MergedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", fmt.Errorf(models.ErrNotFound)
		}
		return nil, "", err
	}

	if PR.Status == "MERGED" {
		return nil, "", fmt.Errorf(models.ErrPRMerged)
	}

	found := false
	for _, r := range reviewers {
		if r == oldUserID {
			found = true
			break
		}
	}
	if !found {
		return nil, "", fmt.Errorf(models.ErrNotAssigned)
	}
	var teamName string
	query = `SELECT t.name 
			 FROM teams t
			 LEFT JOIN users u ON u.team_id = t.id
			 WHERE u.id = $1`
	if err := database.DB.QueryRow(query, PR.AuthorID).Scan(&teamName); err != nil {
		return nil, "", fmt.Errorf(models.ErrNotFound)
	}

	rows, err := database.DB.Query(`
		SELECT u.id
		FROM users u
		JOIN teams t ON u.team_id = t.id
		WHERE t.name = $1 AND u.is_active = TRUE AND u.id <> $2 AND NOT u.id = ANY($3) AND u.id <> $4
	`, teamName, oldUserID, pq.Array(reviewers), PR.AuthorID)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	candidates := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, "", err
		}
		candidates = append(candidates, id)
	}

	if len(candidates) == 0 {
		return nil, "", fmt.Errorf(models.ErrNoCandidate)
	}

	newUser := candidates[0]
	for i, r := range reviewers {
		if r == oldUserID {
			reviewers[i] = newUser
			break
		}
	}

	_, err = database.DB.Exec(`UPDATE pull_requests SET reviewer_ids = $1 WHERE id = $2`, pq.Array(reviewers), PR.PullRequestID)
	if err != nil {
		return nil, "", err
	}

	PR.AssignedReviewers = reviewers

	return PR, newUser, nil
}

func (prs *PullRequestService) сhooseReviewers(members []models.TeamMember) ([]models.TeamMember, error) {
	if len(members) <= 2 {
		return members, nil
	}

	rand.Seed(time.Now().UnixNano())

	shuffled := append([]models.TeamMember(nil), members...)

	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled[:2], nil
}
