package auth_controllers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis-caching-poc/config"
	"github.com/go-redis-caching-poc/types"
)

func SignUp(c *gin.Context) {
	var input types.SignUpInput
	bindDataErr := c.Bind(&input);
	errorWriter(bindDataErr, 400, "field is missing", c)

	valid := validName (input.Name) && validEmail (input.Email) && validPassword (input.Password)

	if !valid {
		errorWriter(fmt.Errorf("valid data not provided"), 400, "valid data not provided", c)
	}

	// Checking if User trying to sign up already exists
	_, userExistsError := checkUserExists(input.Email, false)

	errorWriter(userExistsError, 400, userExistsError.Error(), c)

	// Encrypting Provided Password and Storing User info in DB.
	savedUser, persistUserError := persistUser(&input);

	errorWriter(persistUserError, 400, persistUserError.Error(), c)

	// Generate Access and Refresh Tokens 
	access_token, accessTokenErr := generateToken(
		config.GetConfig().JWT_SECRET_KEY, 
		time.Duration(config.GetConfig().ACCESS_TOKEN_EXPIRY),
		savedUser.Id,
	)

	refresh_token, refreshTokenErr := generateToken(
		config.GetConfig().JWT_SECRET_KEY, 
		time.Duration(config.GetConfig().REFRESH_TOKEN_EXPIRY),
		savedUser.Id,
	)

	errorWriter(accessTokenErr, 400, accessTokenErr.Error(), c)
	errorWriter(refreshTokenErr, 400, refreshTokenErr.Error(), c)

	// Successful Sign Up Returns Access and Refresh Tokens :
	c.JSON(200, gin.H{
		"status": true,
		"access_token": access_token,
		"refresh_token": refresh_token,
	})
}

func LogIn(c *gin.Context) {
	var input types.LogInInput
	bindDataErr := c.Bind(&input);
	errorWriter(bindDataErr, 400, "field is missing", c)

	valid := validEmail (input.Email) && validPassword (input.Password)

	if !valid {
		errorWriter(fmt.Errorf("has incorrect password or email type"), 400, "has incorrect password or email type", c)
	}

	loginUserData, userExistsError := checkUserExists(input.Email, true)

	errorWriter(userExistsError, 400, userExistsError.Error(), c)

	passwordsMatch := ComparePasswords(input.Password, loginUserData.Password);

	if !passwordsMatch {
		errorWriter(fmt.Errorf("user unauthorized. please check password or email"), 400, "user unauthorized. please check password or email", c)
	}

	// Generate Access and Refresh Tokens 
	access_token, accessTokenErr := generateToken(
		config.GetConfig().JWT_SECRET_KEY, 
		time.Duration(config.GetConfig().ACCESS_TOKEN_EXPIRY),
		loginUserData.Id,
	)

	refresh_token, refreshTokenErr := generateToken(
		config.GetConfig().JWT_SECRET_KEY, 
		time.Duration(config.GetConfig().REFRESH_TOKEN_EXPIRY),
		loginUserData.Id,
	)

	errorWriter(accessTokenErr, 400, accessTokenErr.Error(), c)
	errorWriter(refreshTokenErr, 400, refreshTokenErr.Error(), c)

	// Successful Sign In Returns Access and Refresh Tokens :
	c.JSON(200, gin.H{
		"status": true,
		"access_token": access_token,
		"refresh_token": refresh_token,
	})
}

func RefreshToken(c *gin.Context) {
	var input types.RefreshTokenInput
	bindDataErr := c.Bind(&input);
	errorWriter(bindDataErr, 400, "field is missing", c)

	input_access_token, input_refresh_token := input.InputAccessToken, input.InputRefreshToken

	oldTokenDerivedPayload, oldTokenErr := ParseExpiredToken(input_access_token, config.GetConfig().JWT_SECRET_KEY);

	errorWriter(oldTokenErr, 400, oldTokenErr.Error(), c)

	refreshTokenPayload, refreshTokenErr := ParseToken(input_refresh_token, config.GetConfig().JWT_SECRET_KEY)

	errorWriter(refreshTokenErr, 400, refreshTokenErr.Error(), c)

	if oldTokenDerivedPayload.UserId != refreshTokenPayload.UserId {
		errorWriter(fmt.Errorf("invalid user"), 400, "invalid user", c)
	}

	// Generate Access and Refresh Tokens 
	access_token, accessTokenErr := generateToken(
		config.GetConfig().JWT_SECRET_KEY, 
		time.Duration(config.GetConfig().ACCESS_TOKEN_EXPIRY),
		oldTokenDerivedPayload.UserId,
	)

	refresh_token, refreshTokenErr := generateToken(
		config.GetConfig().JWT_SECRET_KEY, 
		time.Duration(config.GetConfig().REFRESH_TOKEN_EXPIRY),
		oldTokenDerivedPayload.UserId,
	)

	errorWriter(accessTokenErr, 400, accessTokenErr.Error(), c)
	errorWriter(refreshTokenErr, 400, refreshTokenErr.Error(), c)

	// Successful Refresh Returns Access and Refresh Tokens :
	c.JSON(200, gin.H{
		"status": true,
		"access_token": access_token,
		"refresh_token": refresh_token,
	})
}