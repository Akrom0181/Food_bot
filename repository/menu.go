package repository

import "bot/models"

func GetMenuItems() ([]models.MenuItem, error) {
	// Logic to retrieve menu items from the database
	return []models.MenuItem{
		{ID: 1, Name: "Pizza", Price: 9.99},
		{ID: 2, Name: "Burger", Price: 5.49},
		{ID: 3, Name: "Soda", Price: 1.99},
	}, nil
}
