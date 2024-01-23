package models

import (
	"context"
	"errors"
	"time"

	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/config"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/utils"
)


type Role struct {
	ID         int       `gorm:"primary_key" json:"id"`
	Name       string    `gorm:"index;size:100;not null" json:"name" binding:"required"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type NewRole struct {
	Name       string `json:"name" binding:"required"`
}

func CreateRole(ctx context.Context, input *NewRole) (*Role, error) {

	db := config.GetDB()
	var count int64

	err := db.WithContext(ctx).Model(&Role{}).Where("name = ?", input.Name).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("duplicate name")
	}

	role := Role{
		Name:       input.Name,
	}

	err = db.WithContext(ctx).Create(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func UpdateRole(ctx context.Context, id int, input *NewRole) (*Role, error) {

	db := config.GetDB()
	var count int64

	err := db.WithContext(ctx).Model(&Role{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count <= 0 {
		return nil, utils.ErrorRecordNotFound
	}

	if err = db.WithContext(ctx).Model(&Role{}).
		Where("name = ?", input.Name).
		Not("id = ?", id).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("duplicate name")
	}

	role := Role{
		ID:         id,
		Name:       input.Name,
	}

	err = db.WithContext(ctx).Model(&role).Updates(map[string]interface{}{
		"Name":       input.Name,
	}).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func DeleteRole(ctx context.Context, id int) (*Role, error) {

	db := config.GetDB()
	var result Role

	err := db.WithContext(ctx).First(&result, id).Error
	if err != nil {
		return nil, utils.ErrorRecordNotFound
	}

	// Delete role module if any
	// err = db.WithContext(ctx).Where("role_id = ?", id).Delete(&RoleModule{}).Error
	// if err != nil {
	// 	return nil, err
	// }

	err = db.WithContext(ctx).Delete(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetRole(ctx context.Context, id int) (*Role, error) {

	db := config.GetDB()
	var result Role

	fieldNames, err := utils.GetQueryFields(ctx, &Role{})
	if err != nil {
		return nil, err
	}

	err = db.WithContext(ctx).Select(fieldNames).First(&result, id).Error
	if err != nil {
		return nil, utils.ErrorRecordNotFound
	}
	return &result, nil
}

func GetRoles(ctx context.Context, name *string) ([]*Role, error) {

	db := config.GetDB()
	var results []*Role

	fieldNames, err := utils.GetQueryFields(ctx, &Role{})
	if err != nil {
		return nil, err
	}

	dbCtx := db.WithContext(ctx)
	if name != nil && len(*name) > 0 {
		dbCtx = dbCtx.Where("name LIKE ?", "%"+*name+"%")
	}
	err = dbCtx.Select(fieldNames).Order("name").Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}