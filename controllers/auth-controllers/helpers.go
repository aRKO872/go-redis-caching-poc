package auth_controllers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	
	"github.com/go-redis-caching-poc/config"
	db_controllers "github.com/go-redis-caching-poc/controllers/db-controllers"
	"github.com/go-redis-caching-poc/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)


func validName (name string) bool {
	name = strings.TrimSpace(name);
	return len(name) > 3 
}

func validEmail (email string) bool {
	email = strings.TrimSpace(email);
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func validPassword(password string) bool {
	password = strings.TrimSpace(password)

	// Check minimum length
	if len(password) < 8 {
		return false
	}

	// Check for at least one lowercase letter
	hasLower := false
	for _, char := range password {
		if 'a' <= char && char <= 'z' {
			hasLower = true
			break
		}
	}

	if !hasLower {
		return false
	}

	// Check for at least one uppercase letter
	hasUpper := false
	for _, char := range password {
		if 'A' <= char && char <= 'Z' {
			hasUpper = true
			break
		}
	}

	if !hasUpper {
		return false
	}

	// Check for at least one digit
	hasDigit := false
	for _, char := range password {
		if '0' <= char && char <= '9' {
			hasDigit = true
			break
		}
	}

	if !hasDigit {
		return false
	}

	// Check for at least one special character
	specialCharRegex := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
	return specialCharRegex.MatchString(password)
}

func checkUserExists (email string, fromLogin bool) (*types.User, error) {
	var user types.User;
	result := db_controllers.DB.Where("email = ?", email).First (&user);

	fmt.Println(result == nil, " is true...")

	if result.Error != nil {
		// User not found
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			if fromLogin {
				return nil, errors.New ("user does not exist. please sign up")
			}
			return nil, nil
		}
		return nil, errors.New ("error occured while checking existence of user")
	}

	// For Registration
	if !fromLogin && user.Id != "" {
		return nil, errors.New ("user already exists! please try logging in")
	}

	if fromLogin {
		return &user, nil
	}
	return nil, nil;
}

func persistUser (
	req *types.SignUpInput,
) (*types.User, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil && len(hashedPasswordBytes) > 0 {
		return nil, errors.New("server error, unable to create your account")
	}

	hashedPassword := string (hashedPasswordBytes[:]);

	saveUser := &types.User {
		Id : uuid.NewString(),
		Name : req.Name,
		Email : req.Email,
		Password : hashedPassword,
	}

	result := db_controllers.DB.Create(&saveUser)

	if result.RowsAffected == 0 {
		return nil, errors.New("error persisting user")
	}

	return saveUser, nil;
}

func generateToken(key string, expiration time.Duration, user_id string) (string, error) {
	claims := &types.Claims{
		UserId: user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration*time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(key[:]))
	if err != nil {
		return "", errors.New("error occured while generating token")
	}

	return signedToken, nil
}

func ComparePasswords(inputPassword string, hashedPassword string) bool {
	passErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))

	return passErr == nil
}

func ParseToken(tokenString string, key string) (*types.Claims, error) {
	claims := &types.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return nil, errors.New("error occurred while parsing token")
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	return claims, nil
}

func ParseExpiredToken(expiredToken string, key string) (*types.Claims, error) {
	claims := &types.Claims{}

	oldTokenPayload, oldTokenErr := jwt.ParseWithClaims(expiredToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().JWT_SECRET_KEY), nil
	})

	claims, ok := oldTokenPayload.Claims.(*types.Claims)
	if !ok {
		return nil, errors.New("error getting custom claims")
	}

	if oldTokenErr != nil {
		expirationTime := claims.ExpiresAt.Time

		if time.Now().After(expirationTime) {
			return claims, nil
		} else {
			return nil, errors.New("invalid token")
		}
	}

	return claims, nil
}