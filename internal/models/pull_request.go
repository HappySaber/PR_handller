package models

type PullRequest struct {
	Id                int    `json:"id"`
	Name              string `json:"name,omitempty"`
	Author            string `json:"author,omitempty"`
	Status            string `json:"status,omitempty"`
	Reviewers         []User `json:"reviewers,omitempty"`
	NeedMoreReviewers bool   `json:"need_more_reviewers,omitempty"`
}
