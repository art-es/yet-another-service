package models

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
}

func (u User) Stored() bool {
	return u.ID != ""
}
