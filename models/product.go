package models

import (
	"context"
	"errors"
	"time"

	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/config"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/utils"
	"gorm.io/gorm"
)

type Product struct {
	ID                              int                `gorm:"primary_key" json:"id"`
	Title                           string             `gorm:"size:255;not null:unique" json:"title" binding:"required,min=3,max=200"`
	Description                     string             `gorm:"type:text;not null" json:"description" binding:"required,min=3"`
	Price                           float64            `gorm:"type:decimal(10,2);not null;default:0.0" json:"price" binding:"required"`
	ComparePrice                    float64            `gorm:"type:decimal(10,2);default:0.0" json:"compare_price"`
	Cost                            float64            `gorm:"type:decimal(10,2);default:0.0" json:"cost"`
	SKU                             string             `gorm:"size:100;unique;default:null" json:"sku"`
	Barcode                         string             `gorm:"size:100;unique;default:null" json:"barcode"`
	IsQtyTracked                    bool               `gorm:"default:false" json:"is_qty_tracked"`
	IsPhysicalProduct               bool               `gorm:"default:false" json:"is_physical_product"`
	IsContinueSellingOutOfStock 	bool               `gorm:"default:false" json:"is_continue_selling_out_of_stock"`
	Weight                          float64            `gorm:"type:decimal(10,2);default:0.0" json:"weight"`
	// Category                 		*Category   	   `gorm:"foreignKey:CategoryId" json:"category"`
	CategoryId               		int                 `gorm:"index;not null" json:"category_id" binding:"required"`
	// Supplier                        *Supplier          `gorm:"foreignKey:SupplierId" json:"supplier"`
	SupplierId                      int                 `gorm:"index;not null" json:"supplier_id" binding:"required"`
	Images                          []Image             `gorm:"polymorphic:Owner" json:"images"`
	// ProductOptions                  []ProductOption     `json:"product_options" binding:"required,dive,required"`
	// ProductVariations               []ProductVariation  `json:"product_variations" binding:"required,dive,required"`
	Tags                            []Tag               `gorm:"many2many:product_tags;"`
	CreatedAt 						time.Time 		    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt 						time.Time 			`gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   					gorm.DeletedAt 	 	`gorm:"index"`
}

type NewProduct struct {
	Title     						string 					`json:"title" binding:"required,min=3,max=50"`
	Description    					string 	    			`json:"description" binding:"required"`
	Price    						float64 	    		`json:"price" binding:"required"`
	ComparePrice    				float64 	    		`json:"compare_price" `
	Cost    						float64 	    		`json:"cost" `
	SKU      						string 					`json:"sku"`
	Barcode      					string 					`json:"barcode"`
	IsQtyTracked                    bool					`json:"is_qty_tracked"`
	IsPhysicalProduct               bool					`json:"is_physical_product"`
	IsContinueSellingOutOfStock 	bool					`json:"is_continue_selling_out_of_stock"`
	Weight                          float64					`json:"weight"`
	CategoryId               		int                		`json:"category_id" binding:"required"`
	SupplierId                      int                		`json:"supplier_id" binding:"required"`
	ProductOptions                  []NewProductOption		`json:"product_options" binding:"required,dive,required"`
	ProductVariations               []NewProductVariation	`json:"product_variations" binding:"required,dive,required"`
	Tags                            []NewTag				`json:"tags"`
	Images                          []NewImage				`json:"image_urls"`
}

type UpdateProductInput struct {
	Title     						string 					`json:"title" binding:"required,min=3,max=50"`
	Description    					string 	    			`json:"description" binding:"required"`
	Price    						float64 	    		`json:"price" binding:"required"`
	ComparePrice    				float64 	    		`json:"compare_price" `
	Cost    						float64 	    		`json:"cost" `
	SKU      						string 					`json:"sku"`
	Barcode      					string 					`json:"barcode"`
	IsQtyTracked                    bool					`json:"is_qty_tracked"`
	IsPhysicalProduct               bool					`json:"is_physical_product"`
	IsContinueSellingOutOfStock 	bool					`json:"is_continue_selling_out_of_stock"`
	Weight                          float64					`json:"weight"`
	CategoryId               		int                		`json:"category_id" binding:"required"`
	SupplierId                      int                		`json:"supplier_id" binding:"required"`
	AddOptions                  	[]NewProductOption		`json:"add_options" binding:"required,dive,required"`
	UpdateOptions                  	[]UpdateProductOption	`json:"update_options" binding:"required,dive,required"`
	DeleteOptions     				[]int               	`json:"delete_options"`
	AddVariations               	[]NewProductVariation	`json:"add_variations" binding:"required,dive,required"`
	UpdateVariations               	[]UpdateProductVariation`json:"update_variations" binding:"required,dive,required"`
	DeleteVariations     			[]int               	`json:"delete_variations"`
	Tags                            []NewTag				`json:"tags"`
	Images                          []NewImage				`json:"image_urls"`
}

type ProductEdge struct {
	Cursor string `json:"cursor"`
	Node   Product `json:"node,omitempty"`
}

type ProductPagination struct {
	Edges    []*ProductEdge `json:"edges"`
	PageInfo *PageInfo     `json:"pageInfo"`
}

func GetPaginatedProducts(ctx context.Context, first *int, after *string) (*ProductPagination, error) {

	decodedCursor, _ := DecodeCursor(after)
	edges := make([]*ProductEdge, *first)
	count := 0
	hasNextPage := false

	db := config.GetDB()
	var results []Product
	var err error

	if decodedCursor == "" {
		err = db.WithContext(ctx).Order("title").Limit(*first + 1).Find(&results).Error
	} else {
		err = db.WithContext(ctx).Order("title").Limit(*first+1).Where("title > ?", decodedCursor).Find(&results).Error
	}

	if err != nil {
		return nil, err
	}

	for _, result := range results {
		// If there are any elements left after the current page
		// we indicate that in the response
		if count == *first {
			hasNextPage = true
		}

		if count < *first {
			edges[count] = &ProductEdge{
				Cursor: EncodeCursor(result.Title),
				Node:   result,
			}
			count++
		}
	}

	pageInfo := PageInfo{
		StartCursor: EncodeCursor(edges[0].Node.Title),
		EndCursor:   EncodeCursor(edges[count-1].Node.Title),
		HasNextPage: &hasNextPage,
	}

	products := ProductPagination{
		Edges:    edges[:count],
		PageInfo: &pageInfo,
	}

	return &products, nil
}

func GetProducts(ctx context.Context, name *string) ([]*Product, error) {

	db := config.GetDB()
	var results []*Product

	if err := db.WithContext(ctx).
				Preload("Images").
				Preload("Tags").
				Find(&results).Error; err != nil {
		return results, errors.New("no Product")
	}

	return results, nil
}

func GetProduct(ctx context.Context, id int) (*Product, error) {

	db := config.GetDB()
	var result Product

	err := db.WithContext(ctx).
				Preload("Tags").
				Preload("Images").
				First(&result, id).Error

	if err != nil {
		return &result, utils.ErrorRecordNotFound
	}

	return &result, nil
}

func CreateProduct(ctx context.Context, input *NewProduct) (*Product, error) {

	db := config.GetDB()

	var count int64

	isValidId := utils.IsRecordValidByID(input.CategoryId, &Category{}, db)

	if !isValidId {
		return &Product{}, errors.New("invalid category id")
	}

	isValidSupplierId := utils.IsRecordValidByID(input.SupplierId, &Supplier{}, db)

	if !isValidSupplierId {
		return &Product{}, errors.New("invalid supplier id")
	}

	err :=  db.WithContext(ctx).Model(&Product{}).Where("sku = ? OR barcode = ? OR title = ? ", input.SKU, input.Barcode, input.Title).
		Count(&count).Error
	if err != nil {
		return &Product{}, err
	}
	if count > 0 {
		return &Product{}, errors.New("duplicate sku or barcode or product title")
	}

	tx := db.Begin()
	
	images, err  :=  mapImageInput(input.Images)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	
	tags, err := createOrUpdateTags(ctx, input.Tags)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	product := Product{
		Title:     						input.Title,
		Description:    				input.Description,
		Price:    						input.Price,
		ComparePrice:     				input.ComparePrice,
		Cost:     						input.Cost,
		SKU:     						input.SKU,
		Barcode:     					input.Barcode,
		IsQtyTracked:     				input.IsQtyTracked,
		IsPhysicalProduct:     			input.IsPhysicalProduct,
		IsContinueSellingOutOfStock:    input.IsContinueSellingOutOfStock,
		Weight:     					input.Weight,
		CategoryId:    					input.CategoryId,
		SupplierId:     				input.SupplierId,
		// ProductVariations:     			productVariations,
		// ProductOptions:     			productOptions,
		Images:     				    images,
		Tags:     						tags,
	}

	err = tx.WithContext(ctx).Create(&product).Error

	if err != nil {
		tx.Rollback()
		return nil, err

	}
	
	if input.ProductVariations != nil {
		if err := mapVariationInput(tx, ctx, product.ID, input.ProductVariations); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if input.ProductOptions != nil {
		if err := mapOptionInput(tx, ctx, product.ID, input.ProductOptions); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err = tx.Commit().Error; err != nil {
       return nil, err

    }

	return &product, nil

}

func UpdateProduct(ctx context.Context, id int, input *UpdateProductInput) (*Product, error) {

	db := config.GetDB()
	tx := db.Begin()

	var count int64

	isValidId := utils.IsRecordValidByID(input.CategoryId, &Category{}, db)

	if !isValidId {
		return &Product{}, errors.New("invalid category id")
	}

	isValidSupplierId := utils.IsRecordValidByID(input.SupplierId, &Supplier{}, db)

	if !isValidSupplierId {
		return &Product{}, errors.New("invalid supplier id")
	}
	

	err :=  tx.WithContext(ctx).Model(&Product{}).
			Where("sku = ? OR barcode = ? OR title = ? ", input.SKU, input.Barcode, input.Title).
			Not("id = ?", id).
			Count(&count).Error

	if err != nil {
		return &Product{}, err
	}

	if count > 0 {
		return &Product{}, errors.New("duplicate sku or barcode or product title")
	}


	var product Product
	if err := tx.WithContext(ctx).First(&product, id).Error; err != nil {
		tx.Rollback()
		return &Product{}, errors.New("error fetching product")
	}

	images, err  :=  mapImageInput(input.Images)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	
	tags, err := createOrUpdateTags(ctx, input.Tags)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	product.Title = input.Title
	product.Description = input.Description
	product.Price = input.Price
	product.ComparePrice = input.ComparePrice
	product.Cost = input.Cost
	product.SKU = input.SKU
	product.Barcode = input.Barcode
	product.IsQtyTracked = input.IsQtyTracked
	product.IsPhysicalProduct = input.IsPhysicalProduct
	product.IsContinueSellingOutOfStock = input.IsContinueSellingOutOfStock
	product.Weight = input.Weight
	product.CategoryId = input.CategoryId
	product.SupplierId = input.SupplierId
	product.Tags = tags
	product.Images = images

	
	if input.AddOptions != nil {
		if err = mapOptionInput(tx, ctx, product.ID, input.AddOptions); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if input.AddVariations != nil {
		if err = mapVariationInput(tx, ctx, product.ID, input.AddVariations); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if input.DeleteOptions != nil {
		err := deleteOptions(tx, ctx, product.ID, input.DeleteOptions)
		if err != nil {
			return nil, err
		}
	}

	if input.DeleteVariations != nil {
		err := deleteVariations(tx, ctx, product.ID, input.DeleteVariations)
		if err != nil {
			return nil, err
		}
	}

	if input.UpdateOptions != nil {
		err := updateOptions(tx, ctx, product.ID, input.UpdateOptions)
		if err != nil {
			return nil, err
		}
	}

	if input.UpdateVariations != nil {
		err := updateVariations(tx, ctx, product.ID, input.UpdateVariations)
		if err != nil {
			return nil, err
		}
	}

	if err =  tx.WithContext(ctx).Save(&product).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err = tx.Commit().Error; err != nil {
       return nil, err
    }

	return &product, nil
}

func DeleteProduct(ctx context.Context, id int) (*Product, error) {

	db := config.GetDB()

	var result Product

	err := db.WithContext(ctx).First(&result, id).Error
	if err != nil {
		return nil, utils.ErrorRecordNotFound
	}

	tx := db.Begin()

	if err := tx.WithContext(ctx).Model(&result).Association("ProductOptions").Unscoped().Clear(); err != nil {
    	tx.Rollback()
        return nil, err
    }
	if err := tx.WithContext(ctx).Model(&result).Association("ProductVariations").Unscoped().Clear(); err != nil {
        tx.Rollback()
		return nil,err
    }
	if err := tx.WithContext(ctx).Model(&result).Association("Images").Unscoped().Clear(); err != nil {
        tx.Rollback()
		return nil,err
    }
	if err := tx.WithContext(ctx).Model(&result).Association("Tags").Unscoped().Clear(); err != nil {
        tx.Rollback()
		return nil,err
    }

	err = db.WithContext(ctx).Delete(&result).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
        return nil, err
    }

	return &result, nil
}

func GetProductIDsByCategoryID(categoryID int) ([]int, error) {

	db := config.GetDB()
	var productIDs []int
	result := db.Model(&Product{}).
				Where("category_id = ?", categoryID).
				Pluck("id", &productIDs)

	if result.Error != nil {
		return nil, result.Error
	}

	return productIDs, nil
}


