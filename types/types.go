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

type UserUpdateInput struct {
	Age int `json:"age"`
	Address string `json:"address"`
	Phone string `json:"phone"`
}

type User struct {
  Id string `json:"id"`;
  Name string `json:"name"`;
  Email string `json:"email"`;
  Password string `json:"password"`;
}

type User_Info struct {
  Id string `json:"id"`;
  Age int `json:"age"`
	Address string `json:"address"`
	Phone string `json:"phone"`
	UserId string `json:"user_id"`
}

type Total_User_Data struct {
	UserId string `json:"user_id"`
	Age int `json:"age"`
	Address string `json:"address"`
	Phone string `json:"phone"`
	Name string `json:"name"`
  Email string `json:"email"`
	UserInfoId string `json:"user_info_id"`
}

type Claims struct {
	UserId string 
	jwt.RegisteredClaims
}