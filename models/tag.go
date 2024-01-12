package models

import (
	"context"

	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/config"
)

type Tag struct {
    ID     		int        `gorm:"primary_key" json:"id"`
    Name   		string	   `gorm:"size:50;not null:unique" json:"name"`
    Products 	[]Product  `gorm:"many2many:product_tags;"`
}

type NewTag struct {
    Name   		string	   `json:"name"`
}

type ProductTags struct {
    ProductID int
    TagID     int
}


func createOrUpdateTags(ctx context.Context, tagInput []NewTag) ([]Tag, error) {
    var tags []Tag

    db := config.GetDB()
    tx := db.Begin()

    for i := range tagInput {
        var existingTag Tag

        // Check if the tag already exists
        if err := tx.WithContext(ctx).Where("name = ?", tagInput[i].Name).First(&existingTag).Error; err != nil {
            // Tag doesn't exist, so create it
            newTag := Tag{Name: tagInput[i].Name}
            if err := tx.WithContext(ctx).Create(&newTag).Error; err != nil {
                tx.Rollback()
                return nil, err
            }
            tags = append(tags, newTag)
        } else {
            // Tag already exists, don't create a new one
            tags = append(tags, existingTag)
        }
    }

    // Commit the transaction
    if err := tx.Commit().Error; err != nil {
        return nil, err
    }

    return tags, nil
}

// func CreateOrUpdateTags(tags []Tag) ([]Tag, error) {
// 	db := config.GetDB()
//     for i := range tags {
//         var existingTag Tag
//         if err := db.Where("name = ?", tags[i].Name).First(&existingTag).Error; err != nil {
// 			// db.WithContext(ctx).Model(&Tag{})Where("name = ?", tags[i].Name).First(&existingTag).Error.
//             // Tag doesn't exist, so create it
//             newTag := Tag{Name: tags[i].Name}
//             if err := db.Create(&newTag).Error; err != nil {
//                 return nil, err
//             }
//             tags[i] = newTag
//         } else {
//             // Tag already exists, associate it
//             tags[i] = existingTag
//         }
//     }
//     return tags, nil
// }

