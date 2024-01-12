package utils

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm/schema"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/validator/v10"
	"github.com/ttacon/libphonenumber"
	"gorm.io/gorm"
)

var CountryCode = "MM"

func IsValidEmail(email string) bool {
	// Basic email validation regex pattern
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

func IsRecordValidByID(id int, model interface{}, db *gorm.DB) bool {

	modelType := reflect.TypeOf(model).Elem() // Get the type of the element (struct)
	record := reflect.New(modelType).Interface()
	// Construct a query using the model's primary key
	query := db.Where("id = ?", id)

	// Perform the query
	if err := query.First(record).Error; err != nil {
		return false // Record with the given ID does not exist
	}

	return true
}

func ValidatePhoneNumber(phoneNumber, countryCode string) error {
	p, err := libphonenumber.Parse(phoneNumber, countryCode)
	if err != nil {
		return err // Phone number is invalid
	}

	if !libphonenumber.IsValidNumber(p) {
		return fmt.Errorf("phone number is not valid")
	}

	return nil // Phone number is valid for the specified country code
}

func GenerateUniqueFilename() string {

	timestamp := time.Now().UnixNano()

	random := rand.Intn(1000)

	uniqueFilename := fmt.Sprintf("%d_%d", timestamp, random)

	return uniqueFilename
}

func ProcessValidationErrors(err error) map[string]string {

	validationErrors := err.(validator.ValidationErrors)

	errorResponse := make(map[string]string)

	for _, ve := range validationErrors {
		errorResponse[ve.Field()] = ve.Tag()
	}

	return errorResponse
}

func NewFalse() *bool {
	b := false
	return &b
}

func GetQueryFields(ctx context.Context, model interface{}) (fieldNames []string, err error) {
	s, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return
	}
	m := make(map[string]string)
	for _, field := range s.Fields {
		dbName := field.DBName
		modelName := strings.ToLower(field.Name)
		m[modelName] = dbName
	}

	fields := graphql.CollectFieldsCtx(ctx, nil)
	for _, column := range fields {
		if !strings.HasPrefix(column.Name, "__") {
			colName := strings.ToLower(column.Name)
			if len(column.Selections) == 0 {
				fieldNames = append(fieldNames, m[colName])
			} else {
				colName += "id"
				fieldNames = append(fieldNames, m[colName])
			}
		}
	}
	return
}

func GetPaginatedQueryFields(ctx context.Context, model interface{}) (fieldNames []string, err error) {
	s, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return
	}
	m := make(map[string]string)
	for _, field := range s.Fields {
		dbName := field.DBName
		modelName := strings.ToLower(field.Name)
		m[modelName] = dbName
	}

	fields := graphql.CollectFieldsCtx(ctx, nil)
	for _, column := range fields {
		if column.Name == "edges" {
			edgesFields := graphql.CollectFields(graphql.GetOperationContext(ctx), column.Selections, nil)
			nodeFields := graphql.CollectFields(graphql.GetOperationContext(ctx), edgesFields[0].Selections, nil)
			for _, nodeColumn := range nodeFields {
				if !strings.HasPrefix(nodeColumn.Name, "__") {
					colName := strings.ToLower(nodeColumn.Name)
					if len(nodeColumn.Selections) == 0 {
						fieldNames = append(fieldNames, m[colName])
					} else {
						colName += "id"
						fieldNames = append(fieldNames, m[colName])
					}
				}
			}
			break
		}
	}
	return
}

// func GetAllQueryFields(ctx context.Context) []string {
// 	return GetNestedPreloads(
// 		graphql.GetOperationContext(ctx),
// 		graphql.CollectFieldsCtx(ctx, nil),
// 		"",
// 	)
// }

// func GetNestedPreloads(ctx *graphql.OperationContext, fields []graphql.CollectedField, prefix string) (preloads []string) {
// 	for _, column := range fields {
// 		prefixColumn := GetPreloadString(prefix, column.Name)
// 		preloads = append(preloads, prefixColumn)
// 		preloads = append(preloads, GetNestedPreloads(ctx, graphql.CollectFields(ctx, column.Selections, nil), prefixColumn)...)
// 	}
// 	return
// }

// func GetPreloadString(prefix, name string) string {
// 	if len(prefix) > 0 {
// 		return prefix + "." + name
// 	}
// 	return name
// }

func CheckImageExistInCloud(imageURL string) (error) {
	
	resp, err := http.Head(imageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return nil
}