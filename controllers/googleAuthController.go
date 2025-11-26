// controllers/googleAuthController.go
package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/Aaryansingh20/jwt/database"
	helper "github.com/Aaryansingh20/jwt/helpers"
	"github.com/Aaryansingh20/jwt/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GoogleAuthCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Query("email")
		name := c.Query("name")
		picture := c.Query("picture")

		if email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
			return
		}

		log.Printf("Processing Google auth for email: %s", email)

		// Check if user exists in PostgreSQL
		var foundUser models.User
		err := database.Client.Where("email = ?", email).First(&foundUser).Error

		if err != nil {
			// User doesn't exist, create new user
			log.Println("User not found, creating new user")
			
			firstName := name
			lastName := ""
			phone := ""
			userType := "USER"
			password := "google-oauth-user"
			userID := uuid.New().String()

			newUser := models.User{
				First_name: &firstName,
				Last_name:  &lastName,
				Email:      &email,
				Password:   &password,
				Phone:      &phone,
				User_type:  &userType,
				Created_at: time.Now(),
				Updated_at: time.Now(),
				User_id:    userID,
			}

			err = database.Client.Create(&newUser).Error
			if err != nil {
				log.Println("Error creating user:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}
			foundUser = newUser
			log.Println("New user created successfully")
		} else {
			log.Println("Existing user found")
		}

		// Generate JWT tokens
		token, refreshToken, err := helper.GenerateAllTokens(
			*foundUser.Email,
			*foundUser.First_name,
			*foundUser.Last_name,
			*foundUser.User_type,
			foundUser.User_id,
		)
		if err != nil {
			log.Println("Error generating tokens:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		log.Println("Tokens generated and updated successfully")

		c.JSON(http.StatusOK, gin.H{
			"token":         token,
			"refresh_token": refreshToken,
			"email":         *foundUser.Email,
			"first_name":    *foundUser.First_name,
			"last_name":     *foundUser.Last_name,
			"user_id":       foundUser.User_id,
			"user_type":     *foundUser.User_type,
			"picture":       picture,
		})
	}
}
