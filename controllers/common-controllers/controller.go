package common_controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis-caching-poc/types"
	"gorm.io/gorm"
)

func PrintHelloWorld(c *gin.Context) {
	c.JSON (200, gin.H{
		"status": true,
		"message": "Hello Word!",
	})
}

func UpdateDetails(c *gin.Context) {
	var input types.UserUpdateInput
	bindDataErr := c.Bind(&input);

	if bindDataErr != nil {
		c.JSON(400, gin.H{
			"status": false,
			"message": "field is missing",
		})
		return
	}

	loggedInUserId, userExtractionErr := ExtractUserId(c)

	if userExtractionErr != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"status": false,
				"message": userExtractionErr.Error(),
			},
		)
		return
	}

	fetchedUserInfo, fetchedUserInfoErr := GetUserInfo(&input, loggedInUserId)

	if fetchedUserInfoErr != nil {
		// Creation of User Info!
		if fetchedUserInfoErr == gorm.ErrRecordNotFound {
			// No user Info Existing Beforehand
			fmt.Println("No Records!")
			newUserInfoId, userInfoCreationErr := CreateNewUserInfo(loggedInUserId, &input)

			if userInfoCreationErr != nil {
				c.JSON(
					400,
					gin.H{
						"status": false,
						"message": userInfoCreationErr.Error(),
					},
				)
			} else {
				// Successfully created user Info
				c.JSON(200, gin.H{
					"status": true,
					"user_id": loggedInUserId,
					"user_info_id": newUserInfoId,
				})
			}
		} else {
			// Some Other Error Encountered
			c.JSON(
				400,
				gin.H{
					"status": false,
					"message": fetchedUserInfoErr.Error(),
				},
			)
		}
		return
	} else {
		// Updation of Existing User Info!
		userInfoUpdationErr := UpdateUserInfo(fetchedUserInfo.Id, &input)

		if userInfoUpdationErr != nil {
			c.JSON(
				400,
				gin.H{
					"status": false,
					"message": userInfoUpdationErr.Error(),
				},
			)
		} else {
			// Successfully created user Info
			c.JSON(200, gin.H{
				"status": true,
				"user_id": loggedInUserId,
				"user_info_id": fetchedUserInfo.Id,
			})
		}

		return
	}
}

func GetAllUserDetails(c *gin.Context) {
	var input types.UserUpdateInput
	bindDataErr := c.Bind(&input);

	if bindDataErr != nil {
		c.JSON(400, gin.H{
			"status": false,
			"message": "field is missing",
		})
		return
	}

	loggedInUserId, userExtractionErr := ExtractUserId(c)

	if userExtractionErr != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"status": false,
				"message": userExtractionErr.Error(),
			},
		)
		return
	}

	totalUserData, fetchTotalUserDataErr := GetTotalUserData(context.Background(), loggedInUserId)

	if fetchTotalUserDataErr != nil {
		c.JSON(
			400,
			gin.H{
				"status": false,
				"message": fetchTotalUserDataErr.Error(),
			},
		)
		return
	}

	c.JSON(
		200,
		gin.H{
			"status": true,
			"user_info": totalUserData,
		},
	)
}