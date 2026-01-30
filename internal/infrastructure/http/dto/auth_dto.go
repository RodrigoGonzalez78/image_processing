package dto

type RegisterRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type LoginRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
