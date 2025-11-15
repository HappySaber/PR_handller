package models

const (
	ErrTeamExists  = "TEAM_EXISTS"
	ErrPRExistst   = "PR_EXISTS"
	ErrPRMerged    = "PR_MERGED"
	ErrNotAssigned = "NOT_ASSIGNED"
	ErrNoCandidate = "NO_CANDIDATE"
	ErrNotFound    = "NOT_FOUND"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponce struct {
	Error Error `json:"error"`
}

func NewErrorResponce(code, message string) ErrorResponce {
	return ErrorResponce{
		Error: Error{
			Code:    code,
			Message: message,
		},
	}
}
