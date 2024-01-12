package models

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/config"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/utils"
	"gorm.io/gorm"
)

type Supplier struct {
	ID        int            `gorm:"primary_key" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name" binding:"required,min=3,max=50"`
	Email     string         `gorm:"size:255;unique" json:"email" binding:"required,email"`
	Address   string         `gorm:"type:text;not null" json:"address" binding:"required"`
	Phone     string         `gorm:"size:255;unique;not null" json:"phone" binding:"required,min=5,max=16"`
	Password  string         `gorm:"default:null;size:100" json:"password"`
	IsActive  *bool          `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type NewSupplier struct {
	Name     	string 			`json:"name" binding:"required,min=3,max=50"`
	Phone    	string 			`json:"phone" binding:"required"`
	Email      	string 			`json:"email" binding:"required,email"`
	Address   	string 			`json:"address" binding:"required,min=5,max=16"`
	Password 	string 			`json:"password"`
	IsActive 	*bool  			`json:"is_active"`
}

func (result *Supplier) PrepareGive() {
	result.Password = ""
}

func GetAllSuppliers(ctx context.Context, name *string) ([]*Supplier, error) {

	db := config.GetDB()
	var results []*Supplier

	fieldNames, err := utils.GetQueryFields(ctx, &Supplier{})
	if err != nil {
		return nil, err
	}

	if err := db.WithContext(ctx).Select(fieldNames).Find(&results).Error; err != nil {
		return results, errors.New("no supplier")
	}

	for i, u := range results {
		u.Password = ""
		results[i] = u
	}

	return results, nil
}

func GetSupplier(ctx context.Context, id int) (*Supplier, error) {

	db := config.GetDB()
	var result Supplier

	fieldNames, err := utils.GetQueryFields(ctx, &Supplier{})
	if err != nil {
		return nil, err
	}

	err = db.WithContext(ctx).Select(fieldNames).First(&result, id).Error

	if err != nil {
		return &result, utils.ErrorRecordNotFound
	}

	result.PrepareGive()

	return &result, nil
}

func CreateSupplier(ctx context.Context, input *NewSupplier) (*Supplier, error) {

	db := config.GetDB()
	var count int64

	if input.Email != "" && !utils.IsValidEmail(input.Email) {
		return &Supplier{}, errors.New("invalid email address")
	}

	if err := utils.ValidatePhoneNumber(input.Phone, utils.CountryCode); err != nil {
		return &Supplier{}, errors.New("phone is invalid")
	}

	err := db.WithContext(ctx).Model(&Supplier{}).Where("phone = ?", input.Phone).Or("email = ?", input.Email).Count(&count).Error
	if err != nil {
		return &Supplier{}, err
	}
	if count > 0 {
		return &Supplier{}, errors.New("duplicate phone or email")
	}

	if input.Password != "" {

		if len(input.Password) < 6 {
			return &Supplier{}, errors.New("password must be at least 6 characters long")
		}
		hashedPassword, err := utils.HashPassword(input.Password)
		if err != nil {
			return &Supplier{}, err
		}
		input.Password = string(hashedPassword)
	}

	supplier := Supplier{
		Name:     input.Name,
		Email:    strings.ToLower(input.Email),
		Phone:    input.Phone,
		Address:  input.Address,
		Password: input.Password,
		IsActive: input.IsActive,
	}

	err = db.WithContext(ctx).Create(&supplier).Error

	if err != nil {
		return &Supplier{}, err
	}

	supplier.Password = ""
	return &supplier, nil

}

func UpdateSupplier(ctx context.Context, id int, input *NewSupplier) (*Supplier, error) {

	db := config.GetDB()
	var count int64

	err := db.WithContext(ctx).Model(&Supplier{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count <= 0 {
		return nil, utils.ErrorRecordNotFound
	}

	if input.Email != "" && !utils.IsValidEmail(input.Email) {
		return &Supplier{}, errors.New("invalid email address")
	}

	if err := utils.ValidatePhoneNumber(input.Phone, utils.CountryCode); err != nil {
		return &Supplier{}, errors.New("phone is invalid")
	}

	if err = db.WithContext(ctx).Model(&Supplier{}).
		Where("email = ? OR phone = ?", input.Email, input.Phone).
		Not("id = ?", id).
		Count(&count).Error; err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("duplicate email or phone")
	}

	if input.Password != "" {
		if len(input.Password) < 6 {
			return &Supplier{}, errors.New("password must be at least 6 characters long")
		}
		hashedPassword, err := utils.HashPassword(input.Password)
		if err != nil {
			return &Supplier{}, err
		}
		input.Password = string(hashedPassword)
	}

	supplier := Supplier{
		ID:        id,
		Name:      input.Name,
		Email:     input.Email,
		Phone:     input.Phone,
		Address:   input.Address,
		IsActive:  input.IsActive,
		Password:  input.Password,
		UpdatedAt: time.Now(),
	}

	err = db.WithContext(ctx).Model(&supplier).Updates(map[string]interface{}{
		"Name":      input.Name,
		"Email":     input.Email,
		"Phone":     input.Phone,
		"Address":   input.Address,
		"IsActive":  input.IsActive,
		"Password":  input.Password,
		"UpdatedAt": time.Now(),
	}).Error

	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

func DeleteSupplier(ctx context.Context, id int) (*Supplier, error) {

	db := config.GetDB()
	var result Supplier

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
