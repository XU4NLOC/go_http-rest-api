package models

type Account struct {
	ID        int    `json: "id"`
	Name      string `json: "name"`
	Type      string `json: "type"`
	CreatedAt string `json: "created_at"`
}
