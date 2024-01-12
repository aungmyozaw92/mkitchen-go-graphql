package middlewares

import (
	"context"

	"gorm.io/gorm"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/models"
)

type productOptionReader struct {
	db *gorm.DB
}


func (r *productOptionReader) GetProductOptions(ctx context.Context, ids []int) []*dataloader.Result[*models.ProductOption] {
	var results []*models.ProductOption

	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&results).Error
	if err != nil {
		return handleError[*models.ProductOption](len(ids), err)
	}

	loaderResults := make([]*dataloader.Result[*models.ProductOption], 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			loaderResults = append(loaderResults, &dataloader.Result[*models.ProductOption]{Data: &models.ProductOption{}})
		} else {
			for _, result := range results {
				if result.ID == id {
					loaderResults = append(loaderResults, &dataloader.Result[*models.ProductOption]{Data: result})
					break
				}
			}
		}
	}
	return loaderResults
}


func GetProductOption(ctx context.Context, id int) (*models.ProductOption, error) {
	loaders := For(ctx)
	return loaders.ProductOptionLoader.Load(ctx, id)()
}

func GetProductOptions(ctx context.Context, ids []int) ([]*models.ProductOption, []error) {
	loaders := For(ctx)
	return loaders.ProductOptionLoader.LoadMany(ctx, ids)()
}


