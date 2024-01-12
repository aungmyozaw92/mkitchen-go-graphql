package models

import (
	"context"
	"errors"
	"time"

	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/config"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/utils"
)

type Branch struct {
	ID        int       `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name" binding:"required"`
	Street1   string    `gorm:"size:255" json:"street1"`
	Street2   string    `gorm:"size:255" json:"street2"`
	City      string    `gorm:"size:255" json:"city"`
	State     string    `gorm:"size:255" json:"state"`
	Phone     string    `gorm:"size:255" json:"phone"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type NewBranch struct {
	Name    string `gorm:"size:255;not null" json:"name" binding:"required"`
	Street1 string `gorm:"size:255" json:"street1"`
	Street2 string `gorm:"size:255" json:"street2"`
	City    string `gorm:"size:255" json:"city"`
	State   string `gorm:"size:255" json:"state"`
	Phone   string `gorm:"size:255" json:"phone"`
}

type BranchEdge struct {
	Cursor string `json:"cursor"`
	Node   Branch `json:"node,omitempty"`
}

type BranchPagination struct {
	Edges    []*BranchEdge `json:"edges"`
	PageInfo *PageInfo     `json:"pageInfo"`
}

func CreateBranch(ctx context.Context, input *NewBranch) (*Branch, error) {

	db := config.GetDB()
	var count int64

	err := db.WithContext(ctx).Model(&Branch{}).Where("name = ?", input.Name).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("duplicate name")
	}

	branch := Branch{
		Name:    input.Name,
		Street1: input.Street1,
		Street2: input.Street2,
		City:    input.City,
		State:   input.State,
		Phone:   input.Phone,
	}

	err = db.WithContext(ctx).Create(&branch).Error
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

func UpdateBranch(ctx context.Context, id int, input *NewBranch) (*Branch, error) {

	db := config.GetDB()
	var count int64

	err := db.WithContext(ctx).Model(&Category{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count <= 0 {
		return nil, utils.ErrorRecordNotFound
	}

	if err = db.WithContext(ctx).Model(&Branch{}).
		Where("name = ?", input.Name).
		Not("id = ?", id).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("duplicate name")
	}

	branch := Branch{
		ID:      id,
		Name:    input.Name,
		Street1: input.Street1,
		Street2: input.Street2,
		City:    input.City,
		State:   input.State,
		Phone:   input.Phone,
	}

	err = db.WithContext(ctx).Model(&branch).Updates(map[string]interface{}{
		"Name":    input.Name,
		"Street1": input.Street1,
		"Street2": input.Street2,
		"City":    input.City,
		"State":   input.State,
		"Phone":   input.Phone,
	}).Error
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

func DeleteBranch(ctx context.Context, id int) (*Branch, error) {

	db := config.GetDB()
	var result Branch

	err := db.WithContext(ctx).First(&result, id).Error
	if err != nil {
		return nil, utils.ErrorRecordNotFound
	}

	err = db.WithContext(ctx).Delete(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetBranch(ctx context.Context, id int) (*Branch, error) {

	db := config.GetDB()
	var result Branch

	fieldNames, err := utils.GetQueryFields(ctx, &Branch{})
	if err != nil {
		return nil, err
	}

	err = db.WithContext(ctx).Select(fieldNames).First(&result, id).Error
	if err != nil {
		return nil, utils.ErrorRecordNotFound
	}
	return &result, nil
}

func GetBranches(ctx context.Context, name *string, city *string) ([]*Branch, error) {

	db := config.GetDB()
	var results []*Branch

	fieldNames, err := utils.GetQueryFields(ctx, &Branch{})
	if err != nil {
		return nil, err
	}

	dbCtx := db.WithContext(ctx)
	if name != nil && len(*name) > 0 {
		dbCtx = dbCtx.Where("name LIKE ?", "%"+*name+"%")
	}
	if city != nil && len(*city) > 0 {
		dbCtx = dbCtx.Where("city LIKE ?", "%"+*city+"%")
	}
	err = dbCtx.Order("name").Select(fieldNames).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func GetPaginatedBranches(ctx context.Context, first *int, after *string) (*BranchPagination, error) {

	decodedCursor, _ := DecodeCursor(after)
	edges := make([]*BranchEdge, *first)
	count := 0
	hasNextPage := false

	db := config.GetDB()
	var results []Branch
	var err error

	fieldNames, err := utils.GetPaginatedQueryFields(ctx, &Branch{})
	if err != nil {
		return nil, err
	}

	if decodedCursor == "" {
		err = db.WithContext(ctx).Select(fieldNames).Order("name").Limit(*first + 1).Find(&results).Error
	} else {
		err = db.WithContext(ctx).Select(fieldNames).Order("name").Limit(*first+1).Where("name > ?", decodedCursor).Find(&results).Error
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
			edges[count] = &BranchEdge{
				Cursor: EncodeCursor(result.Name),
				Node:   result,
			}
			count++
		}
	}

	pageInfo := PageInfo{
		StartCursor: EncodeCursor(edges[0].Node.Name),
		EndCursor:   EncodeCursor(edges[count-1].Node.Name),
		HasNextPage: &hasNextPage,
	}

	branches := BranchPagination{
		Edges:    edges[:count],
		PageInfo: &pageInfo,
	}

	return &branches, nil
}
