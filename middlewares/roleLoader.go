package middlewares

import (
	"context"

	"github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"

	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/models"
)

type roleReader struct {
	db *gorm.DB
}

func (r *roleReader) getRoles(ctx context.Context, ids []int) []*dataloader.Result[*models.Role] {
	var results []*models.Role

	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&results).Error
	if err != nil {
		return handleError[*models.Role](len(ids), err)
	}

	loaderResults := make([]*dataloader.Result[*models.Role], 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			loaderResults = append(loaderResults, &dataloader.Result[*models.Role]{Data: &models.Role{}})
		} else {
			for _, result := range results {
				if result.ID == id {
					loaderResults = append(loaderResults, &dataloader.Result[*models.Role]{Data: result})
					break
				}
			}
		}
	}
	return loaderResults
}

func GetRole(ctx context.Context, id int) (*models.Role, error) {
	loaders := For(ctx)
	return loaders.RoleLoader.Load(ctx, id)()
}

func GetRoles(ctx context.Context, ids []int) ([]*models.Role, []error) {
	loaders := For(ctx)
	return loaders.RoleLoader.LoadMany(ctx, ids)()
}
