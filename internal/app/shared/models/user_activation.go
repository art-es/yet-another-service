package models

type UserActivation struct {
	Token  string
	UserID string
}

func (a UserActivation) Stored() bool {
	return a.Token == ""
}
