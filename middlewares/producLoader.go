package middlewares

import (
	"context"

	"gorm.io/gorm"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/models"
)

type productReader struct {
	db *gorm.DB
}

func (r *productReader) getProducts(ctx context.Context, ids []int) []*dataloader.Result[*models.Product] {
	var results []*models.Product

	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&results).Error
	if err != nil {
		// Instead of returning []error, create a single error for the dataloader.Result
		return handleError[*models.Product](len(ids), err)
	}

	loaderResults := make([]*dataloader.Result[*models.Product], 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			loaderResults = append(loaderResults, &dataloader.Result[*models.Product]{Data: &models.Product{}})
		} else {
			for _, result := range results {
				if result.ID == id {
					loaderResults = append(loaderResults, &dataloader.Result[*models.Product]{Data: result})
					break
				}
			}
		}
	}
	return loaderResults
}

func GetProduct(ctx context.Context, id int) (*models.Product, error) {
	loaders := For(ctx)
	return loaders.ProductLoader.Load(ctx, id)()
}

func GetProducts(ctx context.Context, ids []int) ([]*models.Product, []error) {
	loaders := For(ctx)
	return loaders.ProductLoader.LoadMany(ctx, ids)()
}
