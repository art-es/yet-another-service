package models

type PasswordRecovery struct {
	Token  string
	UserID string
}

func (r *PasswordRecovery) Stored() bool {
	return r.Token != ""
}
