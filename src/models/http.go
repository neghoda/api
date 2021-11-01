package models

// swagger:model
type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// swagger:model
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// swagger:model
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// swagger:model
type LoginResponse struct {
	TokenPair `json:"tokens_pair"`
}
