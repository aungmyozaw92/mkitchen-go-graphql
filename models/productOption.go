package models

import (
	"context"
	"errors"
	"time"

	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/config"
	"gorm.io/gorm"
)

type ProductOption struct {
	ID        		int   				`gorm:"primary_key" json:"id"`
    ProductId 		int   				`gorm:"index;not null" json:"product_id" binding:"required"`
    OptionName      string 				`gorm:"size:255;not null" json:"option_name" binding:"required,min=3,max=30"`
    OptionValue     string 				`gorm:"size:255;not null" json:"option_value" binding:"required,min=3,max=30"`
	CreatedAt 		time.Time 		    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt 		time.Time 			`gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   	gorm.DeletedAt 	 	`gorm:"index"`
}

type UpdateProductOption struct {
	ID        		int   				`json:"id"`
    OptionName      string 				`json:"option_name" binding:"required,min=3,max=30"`
    OptionValue     string 				`json:"option_value" binding:"required,min=3,max=30"`
}

type NewProductOption struct {
    OptionName      string 				`json:"option_name" binding:"required,min=3,max=30"`
    OptionValue     string 				`json:"option_value" binding:"required,min=3,max=30"`
}

func GetProductIDsByOptionId(productID int) ([]int, error) {

	db := config.GetDB()
	var optionIDs []int
	result := db.Model(&ProductOption{}).
				Where("product_id = ?", productID).
				Pluck("id", &optionIDs)

	if result.Error != nil {
		return nil, result.Error
	}

	return optionIDs, nil
}

func mapOptionInput(tx *gorm.DB, ctx context.Context, productId int, productOptionInput []NewProductOption) error {

	for _, optionRequest := range productOptionInput {

		productOption := ProductOption{
			ProductId:  productId,
			OptionName:  optionRequest.OptionName,
			OptionValue: optionRequest.OptionValue,
		}

		err := tx.WithContext(ctx).Create(&productOption).Error

        if err != nil {
            return err
        }
	}

    return nil
}

func updateOptions(tx *gorm.DB, ctx context.Context, productId int, productOptionInput []UpdateProductOption) error {

	for _, optionRequest := range productOptionInput {

		var existingOption ProductOption

		if err := tx.WithContext(ctx).Where("ID = ? AND product_id = ?", optionRequest.ID, productId).First(&existingOption).Error; err != nil {
			return err	
		}
		existingOption.OptionName = optionRequest.OptionName
		existingOption.OptionValue = optionRequest.OptionValue

		if err := tx.WithContext(ctx).Save(&existingOption).Error; err != nil {
			return err
		}
	}

    return nil
}

func deleteOptions(tx *gorm.DB, ctx context.Context, productId int, deleteOptions []int) error {

	for _, deleteID := range deleteOptions {

		var existingOption ProductOption

		if err := tx.WithContext(ctx).
					Where("ID = ? AND product_id = ?", deleteID, productId).
					First(&existingOption).Error; err != nil {
			return errors.New("invalid id in delete options")
		}

		if err := tx.WithContext(ctx).Delete(&existingOption).Error; err != nil {
			return err
		}
	}

	return nil
}

// func mapOptionInput(productOptionInput []NewProductOption) ([]ProductOption) {

// 	var productOptions []ProductOption

// 	for _, optionRequest := range productOptionInput {
// 		productOption := ProductOption{
// 			OptionName:  optionRequest.OptionName,
// 			OptionValue: optionRequest.OptionValue,
// 		}
// 		productOptions = append(productOptions, productOption)
// 	}

// 	return productOptions
// }


