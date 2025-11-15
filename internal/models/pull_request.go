package models

import "time"

type PullRequest struct {
	Id                string    `json:"id"`
	Title             string    `json:"title,omitempty"`
	AuthorId          string    `json:"author_id,omitempty"`
	Status            string    `json:"status,omitempty"`
	Reviewers         []User    `json:"reviewers,omitempty"`
	NeedMoreReviewers bool      `json:"need_more_reviewers,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	MergedAt          time.Time `json:"merged_at"`
}
