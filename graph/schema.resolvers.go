package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/middlewares"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/models"
)

// ParentCategory is the resolver for the parentCategory field.
func (r *categoryResolver) ParentCategory(ctx context.Context, obj *models.Category) (*models.Category, error) {
	return middlewares.GetCategory(ctx, obj.ParentCategoryId)
}

// Products is the resolver for the products field.
func (r *categoryResolver) Products(ctx context.Context, obj *models.Category) ([]*models.Product, error) {
	ids, err := models.GetProductIDsByCategoryID(obj.ID)

	if err != nil {
		return nil, err
	}

	productResults, errors := middlewares.GetProducts(ctx, ids)

	if len(errors) > 0 {
		return nil, errors[0]
	}

	return productResults, nil
}

// OwnerID is the resolver for the owner_id field.
func (r *imageResolver) OwnerID(ctx context.Context, obj *models.Image) (*int, error) {
	panic(fmt.Errorf("not implemented: OwnerID - owner_id"))
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, username string, password string) (*models.LoginInfo, error) {
	return models.Login(ctx, username, password)
}

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, input models.NewUser) (*models.User, error) {
	return models.CreateUser(ctx, &input)
}

// CreateBranch is the resolver for the createBranch field.
func (r *mutationResolver) CreateBranch(ctx context.Context, input models.NewBranch) (*models.Branch, error) {
	return models.CreateBranch(ctx, &input)
}

// UpdateBranch is the resolver for the updateBranch field.
func (r *mutationResolver) UpdateBranch(ctx context.Context, id int, input models.NewBranch) (*models.Branch, error) {
	return models.UpdateBranch(ctx, id, &input)
}

// DeleteBranch is the resolver for the deleteBranch field.
func (r *mutationResolver) DeleteBranch(ctx context.Context, id int) (*models.Branch, error) {
	return models.DeleteBranch(ctx, id)
}

// CreateRole is the resolver for the createRole field.
func (r *mutationResolver) CreateRole(ctx context.Context, input models.NewRole) (*models.Role, error) {
	return models.CreateRole(ctx, &input)
}

// UpdateRole is the resolver for the updateRole field.
func (r *mutationResolver) UpdateRole(ctx context.Context, id int, input models.NewRole) (*models.Role, error) {
	return models.UpdateRole(ctx, id, &input)
}

// DeleteRole is the resolver for the deleteRole field.
func (r *mutationResolver) DeleteRole(ctx context.Context, id int) (*models.Role, error) {
	return models.DeleteRole(ctx, id)
}

// CreateCategory is the resolver for the createCategory field.
func (r *mutationResolver) CreateCategory(ctx context.Context, input models.NewCategory) (*models.Category, error) {
	return models.CreateCategory(ctx, &input)
}

// UpdateCategory is the resolver for the updateCategory field.
func (r *mutationResolver) UpdateCategory(ctx context.Context, id int, input models.NewCategory) (*models.Category, error) {
	return models.UpdateCategory(ctx, id, &input)
}

// DeleteCategory is the resolver for the deleteCategory field.
func (r *mutationResolver) DeleteCategory(ctx context.Context, id int) (*models.Category, error) {
	return models.DeleteCategory(ctx, id)
}

// CreateSupplier is the resolver for the createSupplier field.
func (r *mutationResolver) CreateSupplier(ctx context.Context, input models.NewSupplier) (*models.Supplier, error) {
	return models.CreateSupplier(ctx, &input)
}

// UpdateSupplier is the resolver for the updateSupplier field.
func (r *mutationResolver) UpdateSupplier(ctx context.Context, id int, input models.NewSupplier) (*models.Supplier, error) {
	return models.UpdateSupplier(ctx, id, &input)
}

// DeleteSupplier is the resolver for the deleteSupplier field.
func (r *mutationResolver) DeleteSupplier(ctx context.Context, id int) (*models.Supplier, error) {
	return models.DeleteSupplier(ctx, id)
}

// UploadSingleImage is the resolver for the uploadSingleImage field.
func (r *mutationResolver) UploadSingleImage(ctx context.Context, file graphql.Upload) (string, error) {
	return models.UploadSingleImage(file)
}

// UploadMultipleImages is the resolver for the uploadMultipleImages field.
func (r *mutationResolver) UploadMultipleImages(ctx context.Context, files []*graphql.Upload) ([]string, error) {
	return models.UploadMultipleImages(files)
}

// CreateProduct is the resolver for the createProduct field.
func (r *mutationResolver) CreateProduct(ctx context.Context, input models.NewProduct) (*models.Product, error) {
	return models.CreateProduct(ctx, &input)
}

// UpdateProduct is the resolver for the updateProduct field.
func (r *mutationResolver) UpdateProduct(ctx context.Context, id int, input models.UpdateProductInput) (*models.Product, error) {
	return models.UpdateProduct(ctx, id, &input)
}

// DeleteProduct is the resolver for the deleteProduct field.
func (r *mutationResolver) DeleteProduct(ctx context.Context, id int) (*models.Product, error) {
	return models.DeleteProduct(ctx, id)
}

// Category is the resolver for the category field.
func (r *productResolver) Category(ctx context.Context, obj *models.Product) (*models.Category, error) {
	return middlewares.GetCategory(ctx, obj.CategoryId)
}

// Supplier is the resolver for the supplier field.
func (r *productResolver) Supplier(ctx context.Context, obj *models.Product) (*models.Supplier, error) {
	return middlewares.GetSupplier(ctx, obj.SupplierId)
}

// ProductOptions is the resolver for the product_options field.
func (r *productResolver) ProductOptions(ctx context.Context, obj *models.Product) ([]*models.ProductOption, error) {
	ids, err := models.GetProductIDsByVariationId(obj.ID)

	if err != nil {
		return nil, err
	}

	productResults, errors := middlewares.GetProductOptions(ctx, ids)

	if len(errors) > 0 {
		return nil, errors[0]
	}

	return productResults, nil
}

// ProductVariations is the resolver for the product_variations field.
func (r *productResolver) ProductVariations(ctx context.Context, obj *models.Product) ([]*models.ProductVariation, error) {
	ids, err := models.GetProductIDsByVariationId(obj.ID)

	if err != nil {
		return nil, err
	}

	productResults, errors := middlewares.GetProductVariations(ctx, ids)

	if len(errors) > 0 {
		return nil, errors[0]
	}

	return productResults, nil
}

// Branch is the resolver for the branch field.
func (r *queryResolver) Branch(ctx context.Context, id int) (*models.Branch, error) {
	return models.GetBranch(ctx, id)
}

// Branches is the resolver for the branches field.
func (r *queryResolver) Branches(ctx context.Context, name *string, city *string) ([]*models.Branch, error) {
	return models.GetBranches(ctx, name, city)
}

// BranchPagination is the resolver for the branchPagination field.
func (r *queryResolver) BranchPagination(ctx context.Context, first *int, after *string) (*models.BranchPagination, error) {
	return models.GetPaginatedBranches(ctx, first, after)
}

// Role is the resolver for the role field.
func (r *queryResolver) Role(ctx context.Context, id int) (*models.Role, error) {
	return models.GetRole(ctx, id)
}

// Roles is the resolver for the roles field.
func (r *queryResolver) Roles(ctx context.Context, name *string) ([]*models.Role, error) {
	return models.GetRoles(ctx, name)
}

// Category is the resolver for the category field.
func (r *queryResolver) Category(ctx context.Context, id int) (*models.Category, error) {
	return models.GetCategory(ctx, id)
}

// Categories is the resolver for the categories field.
func (r *queryResolver) Categories(ctx context.Context, name *string) ([]*models.Category, error) {
	return models.GetCategories(ctx, name)
}

// Supplier is the resolver for the supplier field.
func (r *queryResolver) Supplier(ctx context.Context, id int) (*models.Supplier, error) {
	return models.GetSupplier(ctx, id)
}

// Suppliers is the resolver for the suppliers field.
func (r *queryResolver) Suppliers(ctx context.Context, name *string) ([]*models.Supplier, error) {
	return models.GetAllSuppliers(ctx, name)
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id int) (*models.User, error) {
	return models.GetUser(ctx, id)
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, name *string) ([]*models.User, error) {
	return models.GetAllUsers(ctx)
}

// Product is the resolver for the product field.
func (r *queryResolver) Product(ctx context.Context, id int) (*models.Product, error) {
	return models.GetProduct(ctx, id)
}

// Products is the resolver for the products field.
func (r *queryResolver) Products(ctx context.Context, name *string) ([]*models.Product, error) {
	return models.GetProducts(ctx, name)
}

// ProductPagination is the resolver for the productPagination field.
func (r *queryResolver) ProductPagination(ctx context.Context, first *int, after *string) (*models.ProductPagination, error) {
	return models.GetPaginatedProducts(ctx, first, after)
}

// Role is the resolver for the role field.
func (r *userResolver) Role(ctx context.Context, obj *models.User) (*models.Role, error) {
	return middlewares.GetRole(ctx, obj.RoleId)
}

// Category returns CategoryResolver implementation.
func (r *Resolver) Category() CategoryResolver { return &categoryResolver{r} }

// Image returns ImageResolver implementation.
func (r *Resolver) Image() ImageResolver { return &imageResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Product returns ProductResolver implementation.
func (r *Resolver) Product() ProductResolver { return &productResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type categoryResolver struct{ *Resolver }
type imageResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type productResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
