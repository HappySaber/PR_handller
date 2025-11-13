package models

type PullRequest struct {
	Id                int
	Name              string
	Author            string
	Status            string
	Reviewers         []User
	NeedMoreReviewers bool
}
