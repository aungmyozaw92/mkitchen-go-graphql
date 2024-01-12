package middlewares

import (
	"context"
	"time"

	"github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"

	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/models"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/utils"
)

type categoryReader struct {
	db *gorm.DB
}

func (r *categoryReader) getCategories(ctx context.Context, ids []int) []*dataloader.Result[*models.Category] {
	var results []*models.Category

	fieldNames, err := utils.GetQueryFields(ctx, &models.Category{})
	if err != nil {
		return handleError[*models.Category](len(ids), err)
	}

	err = r.db.WithContext(ctx).Select(fieldNames).Where("id IN ?", ids).Find(&results).Error
	if err != nil {
		return handleError[*models.Category](len(ids), err)
	}

	loaderResults := make([]*dataloader.Result[*models.Category], 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			loaderResults = append(loaderResults, &dataloader.Result[*models.Category]{Data: &models.Category{CreatedAt: time.Now(), UpdatedAt: time.Now()}})
		} else {
			for _, result := range results {
				if result.ID == id {
					loaderResults = append(loaderResults, &dataloader.Result[*models.Category]{Data: result})
					break
				}
			}
		}
	}
	return loaderResults
}

func GetCategory(ctx context.Context, id int) (*models.Category, error) {
	loaders := For(ctx)
	return loaders.CatgoryLoader.Load(ctx, id)()
}

func GetCategories(ctx context.Context, ids []int) ([]*models.Category, []error) {
	loaders := For(ctx)
	return loaders.CatgoryLoader.LoadMany(ctx, ids)()
}
