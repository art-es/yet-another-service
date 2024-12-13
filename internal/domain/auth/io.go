package auth

type SignupIn struct {
	Name     string
	Email    string
	Password string
}

type LoginIn struct {
	Email    string
	Password string
}

type LoginOut struct {
	AccessToken  string
	RefreshToken string
}

type LogoutIn struct {
	AccessToken  *string
	RefreshToken string
}

type PasswordRecoverIn struct {
	Token       string
	OldPassword string
	NewPassword string
}