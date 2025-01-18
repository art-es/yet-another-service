package dto

type User struct {
	ID           string
	DisplayName  string
	NickName     string
	Email        string
	PasswordHash string
}

func (u User) Stored() bool {
	return u.ID != ""
}
