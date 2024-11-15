// handlers/users_handlers.gho

package handlers

import (
	"ekeberg.com/messaging-api-postgresql-go/models"
	"ekeberg.com/messaging-api-postgresql-go/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignUp(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "users_h::signup()::Could not parse request data."})
		return
	}

	err = user.SignUpUser()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "users_h::signup()::Could not save user."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "users_h::signup()::User created successfully"})
}

// handlers/users_h.go::Login()
func Login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "users_h::Login()::Could not parse request data: " + err.Error(),
		})
		return
	}

	// Attempt to login the user
	err = user.LoginUser()
	if err != nil {
		if err.Error() == "user is not approved" {
			context.JSON(http.StatusForbidden, gin.H{
				"message": "users_h::Login()::User is not approved. Please contact support for assistance.",
			})
		} else if err.Error() == "user not found or credentials invalid" {
			context.JSON(http.StatusUnauthorized, gin.H{
				"message": "users_h::Login()::Invalid email or password. Please try again.",
			})
		} else {
			context.JSON(http.StatusUnauthorized, gin.H{
				"message": "users_h::Login()::Could not authenticate user: " + err.Error(),
			})
		}
		return
	}

	// Generate token if login is successful
	token, err := utils.GenerateToken(user.Email, user.ID, user.HumanOrService) // Include human_or_service
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "users_h::Login()::Could not generate authentication token: " + err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "users_h::Login()::Login successful!",
		"token":   token,
	})
}
