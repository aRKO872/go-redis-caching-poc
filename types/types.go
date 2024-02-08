package types

import "github.com/golang-jwt/jwt/v5"

type SignUpInput struct {
	Name string `json:"name"`
	Password string `json:"password"`
	Email string `json:"email"`
}

type LogInInput struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type RefreshTokenInput struct {
	InputAccessToken string `json:"input_access_token"`
	InputRefreshToken string `json:"input_refresh_token"`
}

type User struct {
  Id string `json:"id"`;
  Name string `json:"name"`;
  Email string `json:"email"`;
  Password string `json:"password"`;
}

type Claims struct {
	UserId string 
	jwt.RegisteredClaims
}