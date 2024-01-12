package models

import (
	"context"
	"errors"
	"time"

	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/config"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/utils"
)

type Category struct {
	ID               int       `gorm:"primary_key" json:"id"`
	Name             string    `gorm:"size:100;not null" json:"name" binding:"required"`
	ParentCategoryId int       `gorm:"not null" json:"parentCategoryId"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type NewCategory struct {
	Name             string `json:"name" binding:"required"`
	ParentCategoryId int    `json:"parentCategoryId" binding:"required"`
}

func CreateCategory(ctx context.Context, input *NewCategory) (*Category, error) {

	db := config.GetDB()
	var count int64

	err := db.WithContext(ctx).Model(&Category{}).Where("name = ?", input.Name).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("duplicate name")
	}

	if input.ParentCategoryId > 0 {
		err := db.WithContext(ctx).Model(&Category{}).Where("id = ?", input.ParentCategoryId).Count(&count).Error
		if err != nil {
			return nil, err
		}
		if count <= 0 {
			return nil, errors.New("parent not found")
		}
	}

	category := Category{
		Name:             input.Name,
		ParentCategoryId: input.ParentCategoryId,
	}

	err = db.WithContext(ctx).Create(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func UpdateCategory(ctx context.Context, id int, input *NewCategory) (*Category, error) {

	db := config.GetDB()
	var count int64

	if id == input.ParentCategoryId {
		return nil, errors.New("self-parent not allowed")
	}

	err := db.WithContext(ctx).Model(&Category{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count <= 0 {
		return nil, utils.ErrorRecordNotFound
	}

	if err = db.WithContext(ctx).Model(&Category{}).
		Where("name = ?", input.Name).
		Not("id = ?", id).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("duplicate name")
	}

	if input.ParentCategoryId > 0 {
		err := db.WithContext(ctx).Model(&Category{}).Where("id = ?", input.ParentCategoryId).Count(&count).Error
		if err != nil {
			return nil, err
		}
		if count <= 0 {
			return nil, errors.New("parent not found")
		}
	}

	category := Category{
		ID:               id,
		Name:             input.Name,
		ParentCategoryId: input.ParentCategoryId,
	}

	err = db.WithContext(ctx).Model(&category).Updates(map[string]interface{}{
		"Name":             input.Name,
		"ParentCategoryId": input.ParentCategoryId,
	}).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func DeleteCategory(ctx context.Context, id int) (*Category, error) {

	db := config.GetDB()
	var result Category

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

func GetCategory(ctx context.Context, id int) (*Category, error) {

	db := config.GetDB()
	var result Category

	fieldNames, err := utils.GetQueryFields(ctx, &Category{})
	if err != nil {
		return nil, err
	}

	err = db.WithContext(ctx).Select(fieldNames).First(&result, id).Error
	if err != nil {
		return nil, utils.ErrorRecordNotFound
	}
	return &result, nil
}

func GetCategories(ctx context.Context, name *string) ([]*Category, error) {

	db := config.GetDB()
	var results []*Category

	fieldNames, err := utils.GetQueryFields(ctx, &Category{})
	if err != nil {
		return nil, err
	}

	dbCtx := db.WithContext(ctx)
	if name != nil && len(*name) > 0 {
		dbCtx = dbCtx.Where("name LIKE ?", "%"+*name+"%")
	}
	err = dbCtx.Order("name").Select(fieldNames).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
