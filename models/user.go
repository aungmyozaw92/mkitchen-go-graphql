package models

import (
	"context"
	"errors"
	"html"
	"strings"
	"time"

	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/config"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/utils"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `gorm:"primary_key" json:"id"`
	Username  string    `gorm:"size:255;not null;unique" json:"username" binding:"required"`
	Name      string    `gorm:"size:255;not null" json:"name" binding:"required"`
	Email     string    `gorm:"default:null;size:255;unique" json:"email"`
	Password  string    `gorm:"size:100;not null" json:"password"`
	IsActive  *bool     `gorm:"not null;default:false" json:"is_active"`
	RoleId     int       `gorm:"not null;default:0" json:"role_id" binding:"required"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type NewUser struct {
	Username string `json:"username" binding:"required"`
	Name     string `binding:"required"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
	IsActive *bool  `json:"is_active" binding:"required"`
	RoleId    int   `json:"role_id" binding:"required"`
}

type LoginInfo struct {
	Token    string `json:"token"`
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	Name     string `json:"name"`
	// Role     string `json:"role"`
}

func (result *User) PrepareGive() {
	result.Password = ""
}

func Login(ctx context.Context, username string, password string) (*LoginInfo, error) {

	db := config.GetDB()
	var err error
	var result LoginInfo

	u := User{}

	err = db.WithContext(ctx).Model(User{}).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return &result, errors.New("invalid username or password")
	}
	err = utils.ComparePassword(u.Password, password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return &result, errors.New("invalid username or password")
	}

	isActive := *u.IsActive
	if !isActive {
		return &result, errors.New("user is disabled")
	}
	token, err := utils.JwtGenerate(u.ID)
	result.Token = token
	result.Name = u.Name
	result.UserId = u.ID
	result.Username = u.Username

	if err != nil {
		return &result, err
	}

	return &result, nil
}

func GetAllUsers(ctx context.Context) ([]*User, error) {

	db := config.GetDB()
	var results []*User

	if err := db.WithContext(ctx).Find(&results).Error; err != nil {
		return results, errors.New("no user")
	}

	for i, u := range results {
		u.Password = ""
		results[i] = u
	}

	return results, nil
}

func CreateUser(ctx context.Context, input *NewUser) (*User, error) {

	db := config.GetDB()
	var count int64

	if input.Email != "" && !utils.IsValidEmail(input.Email) {
		return &User{}, errors.New("invalid email address")
	}

	err := db.WithContext(ctx).Model(&User{}).Where("username = ?", input.Username).Or("email = ?", input.Email).Count(&count).Error
	if err != nil {
		return &User{}, err
	}
	if count > 0 {
		return &User{}, errors.New("duplicate username or email")
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return &User{}, err
	}

	user := User{
		Username: html.EscapeString(strings.TrimSpace(input.Username)),
		Name:     input.Name,
		Email:    strings.ToLower(input.Email),
		Password: string(hashedPassword),
		IsActive: input.IsActive,
		RoleId: input.RoleId,
	}

	err = db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return &User{}, err
	}
	user.Password = ""
	return &user, nil
}

func GetUser(ctx context.Context, id int) (*User, error) {

	db := config.GetDB()
	var result User

	err := db.WithContext(ctx).First(&result, id).Error

	if err != nil {
		return &result, utils.ErrorRecordNotFound
	}

	result.PrepareGive()

	return &result, nil
}

func (input *User) UpdateUser(id int) (*User, error) {

	db := config.GetDB()
	var count int64

	err := db.Model(&User{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return &User{}, err
	}
	if count <= 0 {
		return nil, utils.ErrorRecordNotFound
	}

	if err = db.Model(&User{}).
		Where("username = ? OR email = ?", input.Username, input.Email).
		Not("id = ?", id).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return &User{}, errors.New("duplicate email or username")
	}

	err = db.Model(&input).Updates(User{Name: input.Name, Email: input.Email, Username: input.Username, IsActive: input.IsActive, RoleId: input.RoleId}).Error
	if err != nil {
		return &User{}, err
	}
	return input, nil
}

func (input *User) DeleteUser(id int) (*User, error) {

	db := config.GetDB()

	err := db.Model(&User{}).Where("id = ?", id).First(&input).Error
	if err != nil {
		return nil, utils.ErrorRecordNotFound
	}

	err = db.Delete(&input).Error
	if err != nil {
		return &User{}, err
	}
	return input, nil
}

func (input *User) ChangeUserPassword() (*User, error) {

	db := config.GetDB()
	//turn password into hash
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return &User{}, err
	}
	input.Password = string(hashedPassword)

	err = db.Model(&User{}).Where("id = ?", input.ID).First(&input).Error
	if err != nil {
		return nil, utils.ErrorRecordNotFound
	}

	err = db.Model(&input).Updates(User{Password: input.Password}).Error
	if err != nil {
		return &User{}, err
	}
	return input, nil
}
