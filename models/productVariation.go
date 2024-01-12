package models

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/config"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/utils"
	"gorm.io/gorm"
)

type ProductVariation struct {
	ID        		int         		`gorm:"primary_key" json:"id"`
    ProductId 		int         		`gorm:"index;not null" json:"product_id"`
    VariantName     string      		`gorm:"size:255;not null" json:"variant_name" binding:"required,min=3,max=30"`
    Price   		float64   			`gorm:"type:decimal(10,2);not null;default:0.0" json:"price" binding:"required"`
    SKU             string    			`gorm:"index;size:100;not null;unique" json:"sku"  binding:"required,min=3,max=50"`
    Barcode         string    			`gorm:"index;size:100;default:null; unique" json:"barcode"  binding:"required,min=3,max=50"`
    ImageUrl        string    			`gorm:"size:100;default:null" json:"image_url"`
    CreatedAt 		time.Time 		    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt 		time.Time 			`gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   	gorm.DeletedAt 	 	`gorm:"index"`
}

type UpdateProductVariation struct {
	ID        		int         		`json:"id"`
    VariantName     string      		`json:"variant_name" binding:"required,min=3,max=30"`
    Price   		float64   			`json:"price" binding:"required"`
    SKU             string    			`json:"sku"  binding:"required,min=3,max=50"`
    Barcode         string    			`json:"barcode"  binding:"required,min=3,max=50"`
    ImageUrl        string    			`json:"image_url"`
}

type NewProductVariation struct {
    VariantName     string      `json:"variant_name" binding:"required,min=3,max=30"`
    Price   		float64   	`json:"price" binding:"required"`
    SKU             string    	`json:"sku"  binding:"required,min=3,max=50"`
    Barcode         string    	`json:"barcode"  binding:"required,min=3,max=50"`
    ImageUrl        string    	`json:"image_url" `
}

func GetProductIDsByVariationId(productID int) ([]int, error) {

	db := config.GetDB()
	var variationIDs []int
	result := db.Model(&ProductVariation{}).
				Where("product_id = ?", productID).
				Pluck("id", &variationIDs)

	if result.Error != nil {
		return nil, result.Error
	}

	return variationIDs, nil
}


func mapVariationInput(tx *gorm.DB, ctx context.Context, productId int, productVariationInput []NewProductVariation) error {

    for _, variationInput := range productVariationInput {
        var count int64

        err := tx.WithContext(ctx).Model(&ProductVariation{}).
                Where("sku = ?", variationInput.SKU).
                Or("barcode = ?",variationInput.Barcode).
            Count(&count).Error

        if err != nil {
            return err
        }

        if count > 0 {
            return errors.New("duplicate sku or barcode in product variants")
        }

        if variationInput.ImageUrl != "" {
            parsedURL, err := url.Parse(variationInput.ImageUrl)
            if err != nil {
                return err
            }

            path := parsedURL.Path
            path = strings.TrimLeft(path, "/")
            
            exists, err := utils.ObjectExists(path)

            if err != nil {
                return err
            }

            if !exists {
                return errors.New("Image does not exist in DigitalOcean Spaces")
            }
        }

        productVariation := ProductVariation{
            ProductId:   productId,
            VariantName: variationInput.VariantName,
            Price:       variationInput.Price,
            SKU:         variationInput.SKU,
            Barcode:     variationInput.Barcode,
            ImageUrl:    variationInput.ImageUrl,
        }

        err = tx.WithContext(ctx).Create(&productVariation).Error
        if err != nil {
            return err
        }
    }

    return nil
}

func updateVariations(tx *gorm.DB, ctx context.Context, productId int, productVariationInput []UpdateProductVariation) error {

    var existingVariation ProductVariation

    for _, variationInput := range productVariationInput {

        if err := tx.WithContext(ctx).Where("ID = ? AND product_id = ?", variationInput.ID, productId).First(&existingVariation).Error; err != nil {
			return err	
		}

        var count int64

        err := tx.WithContext(ctx).Model(&ProductVariation{}).
                Not("id = ?", variationInput.ID).
                Where("sku = ?", variationInput.SKU).
                Or("barcode = ?",variationInput.Barcode).
            Count(&count).Error

        if err != nil {
            return err
        }

        if count > 0 {
            return errors.New("duplicate sku or barcode in updated variants")
        }

        if variationInput.ImageUrl != "" {
            parsedURL, err := url.Parse(variationInput.ImageUrl)
            if err != nil {
                return err
            }

            path := parsedURL.Path
            path = strings.TrimLeft(path, "/")
            
            exists, err := utils.ObjectExists(path)

            if err != nil {
                return err
            }

            if !exists {
                return errors.New("Image does not exist in DigitalOcean Spaces")
            }
        }

        existingVariation.VariantName = variationInput.VariantName
        existingVariation.Price = variationInput.Price
        existingVariation.SKU = variationInput.SKU
        existingVariation.Barcode = variationInput.Barcode
        existingVariation.ImageUrl = variationInput.ImageUrl

        if err := tx.WithContext(ctx).Save(&existingVariation).Error; err != nil {
            return  err
        }
    }

    return nil
}

func deleteVariations(tx *gorm.DB, ctx context.Context, productId int, deleteVariations []int) error {

	for _, deleteID := range deleteVariations {

		var existingVariation ProductVariation

		if err := tx.WithContext(ctx).
					Where("ID = ? AND product_id = ?", deleteID, productId).
					First(&existingVariation).Error; err != nil {
			return errors.New("invalid id in delete Variations")
		}

		if err := tx.WithContext(ctx).Delete(&existingVariation).Error; err != nil {
			return err
		}
	}

	return nil
}