package models

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/99designs/gqlgen/graphql"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/utils"
)

type Image struct {
	ID         int   
    ImageUrl   string `json:"image_url"`
    OwnerType  string `json:"owner_type"`
    OwnerID    uint   `json:"owner_id"`
}

type NewImage struct {
    ImageUrl   string  `json:"image_url"`
}

func UploadSingleImage(file graphql.Upload) (string, error) {

    storagePath := "products/"
	uniqueFilename := utils.GenerateUniqueFilename()
	imageObjectURL := filepath.Join(storagePath, uniqueFilename)

	// Read the uploaded file
	data, err := ioutil.ReadAll(file.File)
	if err != nil {
		return "", err
	}

	// Encode the file data to base64
	imageData := base64.StdEncoding.EncodeToString(data)

	// Save the image to Minio
	err = utils.SaveImageToSpaces(imageObjectURL, imageData)
	if err != nil {
		return "", err
	}

    cloudURL := "https://" + os.Getenv("SP_BUCKET") + "." + os.Getenv("SP_URL") + "/" + imageObjectURL

	return cloudURL, nil
}

func UploadMultipleImages(files []*graphql.Upload) ([]string, error) {
	var fileURLs []string

	for _, file := range files {
		// Process each file (similar to the single file upload logic)
		storagePath := "products/"
		uniqueFilename := utils.GenerateUniqueFilename()
		imageObjectURL := filepath.Join(storagePath, uniqueFilename)

		// Read the uploaded file
		data, err := ioutil.ReadAll(file.File)
		if err != nil {
			return nil, err
		}
		// Encode the file data to base64 or save directly to Minio or other storage
		imageData := base64.StdEncoding.EncodeToString(data)

        
		// Save the image to Minio or other storage
		err = utils.SaveImageToSpaces(imageObjectURL, imageData)
		if err != nil {
			return nil, err
		}

        cloudURL := "https://" + os.Getenv("SP_BUCKET") + "." + os.Getenv("SP_URL") + "/" + imageObjectURL
		fileURLs = append(fileURLs, cloudURL)
	}

	return fileURLs, nil
}

func mapImageInput(imageInput []NewImage) ([]Image, error) {

	var images []Image

	for _, image := range imageInput {
		err := utils.CheckImageExistInCloud(image.ImageUrl)
		
		if err != nil {
			fmt.Println("Error checking image existence:", err)
			return nil, err 
		}

		images = append(images, Image{ImageUrl: image.ImageUrl})
		
	}


	return images, nil
}
