package middlewares

import (
	"context"

	"github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"

	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/models"
)

type supplierReader struct {
	db *gorm.DB
}

func (r *supplierReader) getSuppliers(ctx context.Context, ids []int) []*dataloader.Result[*models.Supplier] {
	var results []*models.Supplier

	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&results).Error
	if err != nil {
		return handleError[*models.Supplier](len(ids), err)
	}

	loaderResults := make([]*dataloader.Result[*models.Supplier], 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			loaderResults = append(loaderResults, &dataloader.Result[*models.Supplier]{Data: &models.Supplier{}})
		} else {
			for _, result := range results {
				if result.ID == id {
					loaderResults = append(loaderResults, &dataloader.Result[*models.Supplier]{Data: result})
					break
				}
			}
		}
	}
	return loaderResults
}

func GetSupplier(ctx context.Context, id int) (*models.Supplier, error) {
	loaders := For(ctx)
	return loaders.SupplierLoader.Load(ctx, id)()
}

func GetSuppliers(ctx context.Context, ids []int) ([]*models.Supplier, []error) {
	loaders := For(ctx)
	return loaders.SupplierLoader.LoadMany(ctx, ids)()
}
