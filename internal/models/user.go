package models

type User struct {
	Id       int
	Name     string `json:"name,omitempty"`
	IsActive bool   `json:"is_active,omitempty"`
}
