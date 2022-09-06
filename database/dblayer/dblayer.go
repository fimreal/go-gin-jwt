package dblayer

import "github.com/fimreal/go-gin-jwt/database/models"

type DBLayer interface {
	AddUser(models.User) (models.User, error)
	// GetUser(username string) (models.User, error)
	GetUser(user interface{}) (models.User, error)
	// UpdateUser()
	// LockUser()
	// ReleaseUser()
	// DeleteUser()
	// FindPassword()
	// ResetPassword()
}
