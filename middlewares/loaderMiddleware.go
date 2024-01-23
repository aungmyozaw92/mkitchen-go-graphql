package middlewares

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/config"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/models"
	"gorm.io/gorm"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

// Loaders wrap your data loaders to inject via middleware
type Loaders struct {
	RoleLoader *dataloader.Loader[int, *models.Role]
	CatgoryLoader *dataloader.Loader[int, *models.Category]
	SupplierLoader *dataloader.Loader[int, *models.Supplier]
	ProductLoader *dataloader.Loader[int, *models.Product]
	ProductVariationLoader *dataloader.Loader[int, *models.ProductVariation]
	ProductOptionLoader *dataloader.Loader[int, *models.ProductOption]
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders(conn *gorm.DB) *Loaders {
	// define the data loader
	
	ar := &categoryReader{db: conn}
	role := &roleReader{db: conn}
	supplier := &supplierReader{db: conn}
	product := &productReader{db: conn}
	productV := &productVariationReader{db: conn}
	productOpt := &productOptionReader{db: conn}

	return &Loaders{
		RoleLoader: dataloader.NewBatchedLoader(role.getRoles, dataloader.WithWait[int, *models.Role](time.Millisecond)),
		CatgoryLoader: dataloader.NewBatchedLoader(ar.getCategories, dataloader.WithWait[int, *models.Category](time.Millisecond)),
		SupplierLoader: dataloader.NewBatchedLoader(supplier.getSuppliers, dataloader.WithWait[int, *models.Supplier](time.Millisecond)),
		ProductLoader: dataloader.NewBatchedLoader(product.getProducts, dataloader.WithWait[int, *models.Product](time.Millisecond)),
		ProductVariationLoader: dataloader.NewBatchedLoader(productV.GetProductVariations, dataloader.WithWait[int, *models.ProductVariation](time.Millisecond)),
		ProductOptionLoader: dataloader.NewBatchedLoader(productOpt.GetProductOptions, dataloader.WithWait[int, *models.ProductOption](time.Millisecond)),
	}
}

func LoaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		loader := NewLoaders(config.GetDB())
		ctx := context.WithValue(c.Request.Context(), loadersKey, loader)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

func handleError[T any](count int, err error) []*dataloader.Result[T] {
	results := make([]*dataloader.Result[T], count)
	for i := 0; i < count; i++ {
		results[i] = &dataloader.Result[T]{Error: err}
	}
	return results
}
