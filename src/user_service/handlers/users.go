package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/thejasmeetsingh/go-ecommerce/src/user_service/internal/database"
	"github.com/thejasmeetsingh/go-ecommerce/src/user_service/utils"
	"github.com/thejasmeetsingh/go-ecommerce/src/user_service/validators"
)

// Common function for getting user object from the context, For all handlers in this module
func getDBUser(c *gin.Context) (database.User, error) {
	user, exists := c.Get("user")

	if !exists {
		return database.User{}, fmt.Errorf("authentication required")
	}

	dbUser, ok := user.(database.User)

	if !ok {
		return database.User{}, fmt.Errorf("invalid user")
	}

	return dbUser, nil
}

// Fetch user profile details
func (apiCfg *ApiConfig) GetUserProfile(c *gin.Context) {
	dbUser, err := getDBUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": DatabaseUserToUser(dbUser)})
}

// Update user profile details
func (apiCfg *ApiConfig) UpdateUserProfile(c *gin.Context) {
	dbUser, err := getDBUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	type Parameters struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	var params Parameters
	err = c.ShouldBindJSON(&params)

	if err != nil {
		log.Errorln(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error while parsing the request data"})
		return
	}

	// check if request body is empty
	if params.Name == "" && params.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	// Pre-fill name and email address with existing data from DB
	if params.Name == "" {
		params.Name = dbUser.Name.String
	}

	isEmailChanged := true
	if params.Email == "" {
		params.Email = dbUser.Email
		isEmailChanged = false
	}

	email := strings.ToLower(params.Email)

	// Validate the given email address
	if !validators.EmailValidator(params.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email address"})
		return
	}

	// Check if user any exists with the new email address
	if isEmailChanged && email != dbUser.Email {
		_, err = apiCfg.Queries.GetUserByEmail(c, email)
		if err == nil {
			log.Errorln(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "User with this email address already exists"})
			return
		}
	}

	// Update the profile details
	user, err := apiCfg.Queries.UpdateUserDetails(c, database.UpdateUserDetailsParams{
		Name: sql.NullString{
			String: params.Name,
			Valid:  true,
		},
		Email:      email,
		ModifiedAt: time.Now().UTC(),
		ID:         dbUser.ID,
	})

	if err != nil {
		log.Errorln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile details updated successfully!", "data": DatabaseUserToUser(user)})
}

// Delete user profile
func (apiCfg *ApiConfig) DeleteUserProfile(c *gin.Context) {
	dbUser, err := getDBUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	if err := apiCfg.Queries.DeleteUser(c, dbUser.ID); err != nil {
		log.Errorln(err)
		c.JSON(http.StatusForbidden, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted successfully!"})
}

// Change password API for authenticated user
func (apiCfg *ApiConfig) ChangePassword(c *gin.Context) {
	dbUser, err := getDBUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	type Parameters struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	var params Parameters
	err = c.ShouldBindJSON(&params)

	if err != nil {
		log.Errorln(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error while parsing the request data"})
		return
	}

	// Check wheather or not old password is correct or not
	match, err := utils.CheckPassowrdValid(params.OldPassword, dbUser.Password)
	if err != nil {
		log.Errorln(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid old password, Please try again."})
		return
	} else if !match {
		log.Errorln(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid old password, Please try again."})
		return
	}

	if params.OldPassword == params.NewPassword {
		c.JSON(http.StatusBadRequest, gin.H{"message": "New password should not be same as old password"})
		return
	}

	// Validate the new password
	if err = validators.PasswordValidator(params.NewPassword, dbUser.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Generate the new hashed password
	hashedPassword, err := utils.GetHashedPassword(params.NewPassword)

	if err != nil {
		log.Errorln(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Something went wrong"})
		return
	}

	// Update the password
	_, err = apiCfg.Queries.UpdateUserPassword(c, database.UpdateUserPasswordParams{
		Password:   hashedPassword,
		ModifiedAt: time.Now().UTC(),
		ID:         dbUser.ID,
	})

	if err != nil {
		log.Errorln(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully!"})
}
