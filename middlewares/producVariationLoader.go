package middlewares

import (
	"context"

	"gorm.io/gorm"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/models"
)

type productVariationReader struct {
	db *gorm.DB
}


func (r *productVariationReader) GetProductVariations(ctx context.Context, ids []int) []*dataloader.Result[*models.ProductVariation] {
	var results []*models.ProductVariation

	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&results).Error
	if err != nil {
		return handleError[*models.ProductVariation](len(ids), err)
	}

	loaderResults := make([]*dataloader.Result[*models.ProductVariation], 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			loaderResults = append(loaderResults, &dataloader.Result[*models.ProductVariation]{Data: &models.ProductVariation{}})
		} else {
			for _, result := range results {
				if result.ID == id {
					loaderResults = append(loaderResults, &dataloader.Result[*models.ProductVariation]{Data: result})
					break
				}
			}
		}
	}
	return loaderResults
}


func GetProductVariation(ctx context.Context, id int) (*models.ProductVariation, error) {
	loaders := For(ctx)
	return loaders.ProductVariationLoader.Load(ctx, id)()
}

func GetProductVariations(ctx context.Context, ids []int) ([]*models.ProductVariation, []error) {
	loaders := For(ctx)
	return loaders.ProductVariationLoader.LoadMany(ctx, ids)()
}


