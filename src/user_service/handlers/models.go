// Custom user model for converting the raw data into a desired JSON data
//
// With keys as formatted as snake_case rather than TitleCase

package handlers

import (
	"github.com/google/uuid"
	"github.com/thejasmeetsingh/go-ecommerce/src/user_service/internal/database"
)

type user struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func DatabaseUserToUser(dbUser database.User) user {
	return user{
		ID:    dbUser.ID,
		Name:  dbUser.Name.String,
		Email: dbUser.Email,
	}
}
