package models

type Users struct {
	UserID   int    `db:"user_id"`
	Metadata string `db:"metadata"`
}
