package common_controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis-caching-poc/config"
	auth_controllers "github.com/go-redis-caching-poc/controllers/auth-controllers"
	db_controllers "github.com/go-redis-caching-poc/controllers/db-controllers"
	redis_controllers "github.com/go-redis-caching-poc/controllers/redis-controllers"
	"github.com/go-redis-caching-poc/types"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

func ExtractUserId (c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	// Check if the Authorization header starts with "Bearer"
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("invalid Authorization header format")
	}

	// Extract the token from the Authorization header
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Now you have the Bearer token in the 'token' variable
	tokenDerivedPayload, tokenErr := auth_controllers.ParseToken(
		token,
		config.GetConfig().JWT_SECRET_KEY,
	)

	if tokenErr != nil {
		c.JSON(
			400,
			gin.H{
				"status": false,
				"message": tokenErr.Error(),
			},
		)
	}

	return tokenDerivedPayload.UserId, nil
}

func GetUserInfo(input *types.UserUpdateInput, userId string) (*types.User_Info, error) {
	var fetchedUserInfo types.User_Info

	if err := db_controllers.DB.Select("id, age, address, phone").Where("user_id = ?", userId).First(&fetchedUserInfo).Error; err != nil {
			return nil, err
	}

	return &fetchedUserInfo, nil
}

func CreateNewUserInfo(loggedInUserId string, input *types.UserUpdateInput) (string, error) {
	userInfo := types.User_Info{
		UserId:  loggedInUserId,
		Age:     input.Age,
		Address: input.Address,
		Phone:   input.Phone,
		Id:			 uuid.NewString(),
}

	// Insert the new UserInfo record into the database
	if err := db_controllers.DB.Create(&userInfo).Error; err != nil {
			return "", err
	}

	// Return the ID of the newly created UserInfo record
	return userInfo.Id, nil
}

func UpdateUserInfo(userInfoID string, input *types.UserUpdateInput) error {
	// Create an instance of UserInfo with the provided userInfoID
	query := "UPDATE \"user_infos\" SET age = ?, address = ?, phone = ? WHERE id = ?;"

	// Perform the update
	if err := db_controllers.DB.Exec(query, input.Age, input.Address, input.Phone, userInfoID).Error; err != nil {
			return err
	}

	return nil
}

func FetchAllUserData(loggedInUserId string) (*types.Total_User_Data, error) {
	var totalUserData types.Total_User_Data

	query := "select u.id as user_id, u.name, u.email, ui.id as user_info_id, ui.age, ui.address, ui.phone FROM users u left join user_infos ui on u.id = ui.user_id WHERE u.id = ?"

	if err := db_controllers.DB.Raw(query, loggedInUserId).Scan(&totalUserData).Error; err != nil {
		return nil, err
	}

	fmt.Println(totalUserData)

	return &totalUserData, nil
}

func SetTotalUserData(ctx context.Context, data *types.Total_User_Data) error {
	// Serialize the struct to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
			return err
	}

	// Set the JSON data in Redis with key as user_id
	err = redis_controllers.Rdb.Set(ctx, data.UserId, jsonData, 0).Err()
	if err != nil {
			return err
	}

	return nil
}

func GetTotalUserData(ctx context.Context, userID string) (*types.Total_User_Data, error) {
	var data *types.Total_User_Data

	// Get the JSON data from Redis with key as user_id
	jsonData, err := redis_controllers.Rdb.Get(ctx, userID).Result()
	if err != nil && err == redis.Nil {
		fmt.Println("Slow")
		data, err = FetchAllUserData(userID)
		if err != nil {
			return nil, err
		} else {
			SetTotalUserData(context.Background(), data)
			return data, nil
		}
	}

	// Deserialize the JSON data into struct
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
			return nil, err
	}
	fmt.Println("Fast")

	return data, nil
}
