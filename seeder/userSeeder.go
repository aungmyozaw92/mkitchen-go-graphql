package seeder

import (
	"fmt"

	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/models"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/utils"
	"gorm.io/gorm"
)

func seedUser(db *gorm.DB) {

	// Seed Roles
	roles := []models.Role{
		{
			Name:     "SuperAdmin",
		},
	}

	err := db.Create(&roles).Error
	if err != nil {
		fmt.Println("Error seeding roles: " + err.Error())
		return
	}

	// Seed Users
	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		fmt.Println("Error hashing password: " + err.Error())
		return
	}

	var role models.Role

	err = db.First(&role).Error

	if err != nil {
		fmt.Println("Role does not exit: " + err.Error())
		return
	}

	users := []models.User{
		{
			Username: "super_admin",
			Name:     "SuperAdmin",
			Email:    "superadmin@example.com",
			RoleId:  role.ID,
			IsActive: new(bool),
			Password: string(hashedPassword),
		},
	}

	err = db.Create(&users).Error
	if err != nil {
		fmt.Println("Error seeding users: " + err.Error())
		return
	}

}