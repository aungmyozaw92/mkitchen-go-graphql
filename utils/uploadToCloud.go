package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func getClient() (*minio.Client, error) {
	endpoint := os.Getenv("SP_URL")
	accessKey := os.Getenv("SP_ACCESS_KEY_ID")
	secretKey := os.Getenv("SP_SECRET_ACCESS_KEY")
	

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func SaveImageToSpaces(objectName, imageData string) error {
    // Decode the base64 data
    decodedData, err := base64.StdEncoding.DecodeString(imageData)
    if err != nil {
        return err
    }
    bucketName := os.Getenv("SP_BUCKET")   

    // Get the Minio client
	client, err := getClient()
    if err != nil {
		return err
	}

    // Upload the decoded image data to the specified object name in your Space
    contentType := "image/jpeg" 

    _, err = client.PutObject(context.Background(), bucketName, objectName, bytes.NewReader(decodedData), int64(len(decodedData)), minio.PutObjectOptions{
        ContentType: contentType,
        UserMetadata: map[string]string{
            "x-amz-acl": "public-read", 
        },
    })
    if err != nil {
        return err
    }

    return nil
}

func DeleteImageFromSpaces(objectName string) error {
	// Get the Minio client
	client, err := getClient()
	if err != nil {
		return err
	}

    bucketName := os.Getenv("SP_BUCKET")  

	// Remove the specified object from your Space
	err = client.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		// Check if the error is due to the object not existing
		if strings.Contains(err.Error(), "The specified key does not exist") {
			fmt.Print("Object does not exist:", objectName)
			return nil
		}
		return err
	}

	fmt.Print("Object deleted successfully:", objectName)
	return nil
}

func ObjectExists(objectName string) (bool, error) {
	// Get the Minio client
	client, err := getClient()
	if err != nil {
		return false, err
	}

	bucketName := os.Getenv("SP_BUCKET")

	// StatObject is used to check the existence of an object without downloading its content
	_, err = client.StatObject(context.Background(), bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "The specified key does not exist") {
			return false, nil // Object does not exist
		}
		return false, err // Other error
	}

	return true, nil // Object exists
}