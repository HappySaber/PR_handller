package models

type Team struct {
	Name    string `json:"name,omitempty"`
	Members []User `json:"members,omitempty"`
}
