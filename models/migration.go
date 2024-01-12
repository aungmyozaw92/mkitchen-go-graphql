package models

import (
	"log"

	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/config"
)

func MigrateTable() {
	db := config.GetDB()

	err := db.AutoMigrate(
		&Branch{}, 
		&Category{}, 
		&User{}, 
		&Supplier{},
		&Product{},
		&ProductOption{},
		&ProductVariation{},
		&Tag{},
		&ProductTags{},
		&Image{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
