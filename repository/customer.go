package repository

import (
	"bot/models"
	"database/sql"
	"fmt"
)

type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// SaveUser saves a new user in the database
func (repo *UserRepository) SaveUser(user models.User) error {
	// Prepare the SQL statement
	query := `INSERT INTO users (username) VALUES ($1) RETURNING id`

	// Execute the SQL statement
	var userID int
	err := repo.DB.QueryRow(query, user.Name).Scan(&userID)
	if err != nil {
		// Handle specific database errors if necessary
		return fmt.Errorf("failed to save user: %w", err)
	}

	// Set the ID field in the user struct if needed
	// user.ID = userID
	return nil
}
