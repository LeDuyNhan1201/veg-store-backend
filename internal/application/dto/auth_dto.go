package dto

type SignInRequest struct {
	Username string `json:"username" binding:"required,range=6-20" example:"admin"`
	Password string `json:"password" binding:"required" example:"password123"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
