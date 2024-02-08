package config

import (
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

type ConfigType struct {
	DATABASE_URL string
	JWT_SECRET_KEY string
	ACCESS_TOKEN_EXPIRY int
	REFRESH_TOKEN_EXPIRY int
}

func GetConfig () ConfigType {
	godotenv.Load(".env")

	var DATABASE_URL = os.Getenv("DATABASE_URL")
	var JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
	var ACCESS_TOKEN_EXPIRY, _ = strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY"))
	var REFRESH_TOKEN_EXPIRY, _ = strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY"))

	return ConfigType{
		DATABASE_URL: DATABASE_URL,
		JWT_SECRET_KEY: JWT_SECRET_KEY,
		ACCESS_TOKEN_EXPIRY: ACCESS_TOKEN_EXPIRY,
		REFRESH_TOKEN_EXPIRY: REFRESH_TOKEN_EXPIRY,
	}
}