package auth

type SignupRequest struct {
	Name     string
	Email    string
	Password string
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
}

type LogoutRequest struct {
	AccessToken  *string
	RefreshToken string
}
