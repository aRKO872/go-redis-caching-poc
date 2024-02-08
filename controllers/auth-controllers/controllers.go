package auth_controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis-caching-poc/config"
	"github.com/go-redis-caching-poc/types"
)

func SignUp(c *gin.Context) {
	var input types.SignUpInput
	bindDataErr := c.Bind(&input);
	if bindDataErr != nil {
		c.JSON(400, gin.H{
			"status": false,
			"message": "field is missing",
		})
		return
	}

	valid := validName (input.Name) && validEmail (input.Email) && validPassword (input.Password)

	if !valid {
		c.JSON(400, gin.H{
			"status": false,
			"message": "valid data not provided",
		})
		return
	}

	// Checking if User trying to sign up already exists
	_, userExistsError := checkUserExists(input.Email, false)

	if userExistsError != nil {
		c.JSON(400, gin.H{
			"status": false,
			"message": userExistsError.Error(),
		})
		return
	}

	// Encrypting Provided Password and Storing User info in DB.
	savedUser, persistUserError := persistUser(&input);

	if persistUserError != nil {
		c.JSON(400, gin.H{
			"status": false,
			"message": persistUserError.Error(),
		})
		return
	}

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

	if accessTokenErr != nil {
		c.JSON(400, gin.H{
			"status": false,
			"message": accessTokenErr.Error(),
		})
		return
	}
	if refreshTokenErr != nil {
		c.JSON(400, gin.H{
			"status": false,
			"message": refreshTokenErr.Error(),
		})
		return
	}

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
	if bindDataErr != nil {
		c.JSON(400, gin.H{
			"status": false,
			"message": "field is missing",
		})
		return
	}

	valid := validEmail (input.Email) && validPassword (input.Password)

	if !valid {
		c.JSON(400, gin.H{
			"status": false,
			"message":  "has incorrect password or email type",
		})
		return
	}

	loginUserData, userExistsError := checkUserExists(input.Email, true)

	if userExistsError != nil {
		c.JSON(400, gin.H{
			"status": false,
			"message": userExistsError.Error(),
		})
		return
	}


	passwordsMatch := ComparePasswords(input.Password, loginUserData.Password);

	if !passwordsMatch {
		c.JSON(400, gin.H{
			"status": false,
			"message":  "user unauthorized. please check password or email",
		})
		return
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

	if accessTokenErr != nil {
		c.JSON(400, gin.H{
			"status": false,
			"message": accessTokenErr.Error(),
		})
		return
	}
	if refreshTokenErr != nil {
		c.JSON(400, gin.H{
			"status": false,
			"message": refreshTokenErr.Error(),
		})
		return
	}

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
	if bindDataErr != nil {
		c.JSON(400, gin.H{
			"status": false,
			"message": bindDataErr.Error(),
		})
		return
	}

	input_access_token, input_refresh_token := input.InputAccessToken, input.InputRefreshToken

	oldTokenDerivedPayload, oldTokenErr := ParseExpiredToken(input_access_token, config.GetConfig().JWT_SECRET_KEY);

	if oldTokenErr != nil {
		c.JSON(400, gin.H{
			"status": false,
			"message": oldTokenErr.Error(),
		})
		return
	}

	refreshTokenPayload, refreshTokenErr := ParseToken(input_refresh_token, config.GetConfig().JWT_SECRET_KEY)

	if refreshTokenErr != nil {
		c.JSON(400, gin.H{
			"status": false,
			"message": refreshTokenErr.Error(),
		})
		return
	}

	if oldTokenDerivedPayload.UserId != refreshTokenPayload.UserId {
		c.JSON(400, gin.H{
			"status": false,
			"message":  "invalid user",
		})
		return
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

	if accessTokenErr != nil {
		c.JSON(400, gin.H{
			"status": false,
			"message": accessTokenErr.Error(),
		})
		return
	}
	if refreshTokenErr != nil {
		c.JSON(400, gin.H{
			"status": false,
			"message": refreshTokenErr.Error(),
		})
		return
	}

	// Successful Refresh Returns Access and Refresh Tokens :
	c.JSON(200, gin.H{
		"status": true,
		"access_token": access_token,
		"refresh_token": refresh_token,
	})
}