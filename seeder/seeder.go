package seeder

import (
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) {
	// Seed data
	seedUser(db)
	
	// Create suppliers
	// suppliers := []models.Supplier{
	// 	{Name: "Supplier1", Email: "ex@gmail.com", Address: "Yangon, Kamayut", Phone: "09420118123"},
	// 	{Name: "Supplier2", Email: "exa2@gmail.com", Address: "Yangon, BaYintNaung", Phone: "09420118124"},
	// }

	// for i := range suppliers {
	// 	db.Create(&suppliers[i])
	// }

	// // Create product categories
	// productCategories := []models.ProductCategory{
	// 	{Name: "Food", NameMM: "Food Myanmar Name"},
	// 	{Name: "Oil", NameMM: "Oil Myanmar Name"},
	// }

	// for i := range productCategories {
	// 	db.Create(&productCategories[i])
	// }

}
