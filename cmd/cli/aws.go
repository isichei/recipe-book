package main

// From Chat

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const bucketFolder = "app-data/"

// Sync local static files from an aws bucket
func syncFromAws(bucket string, dataPath string) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	// Download db
	if err := downloadObject(context.TODO(), client, bucket, bucketFolder+"recipes.db", filepath.Join(dataPath, "/recipes.db")); err != nil {
		log.Fatalf("css download failed: %v", err)
	}

	// Download css
	if err := downloadObject(context.TODO(), client, bucket, bucketFolder+"static/css/styles.css", filepath.Join(dataPath, "/static/css/styles.css")); err != nil {
		log.Fatalf("css download failed: %v", err)
	}

	// Download images
	err = downloadFolderFromS3(context.TODO(), client, bucket, bucketFolder+"static/img", filepath.Join(dataPath, "static/img/"))
	if err != nil {
		log.Fatalf("images download failed: %v", err)
	}

	fmt.Println("Download complete!")
}

// downloadFolderFromS3 lists files in an S3 folder for a given prefix and downloads them to the specified local directory.
func downloadFolderFromS3(ctx context.Context, client *s3.Client, bucket string, objPrefix string, localDir string) error {

	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(objPrefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to list objects: %v\n", err)
		}

		for _, obj := range page.Contents {
			// define localFilePath based on object key
			key := *obj.Key
			strings.TrimPrefix(key, "/")
			remainingFilepath := strings.TrimPrefix(strings.TrimPrefix(key, "/"), objPrefix)

			// Download each object
			if strings.Contains(remainingFilepath, ".") {
				if err := downloadObject(ctx, client, bucket, key, filepath.Join(localDir, remainingFilepath)); err != nil {
					return fmt.Errorf("failed to download object %s: %v to %s\n", key, err, localDir)
				}
			}
		}
	}
	return nil
}

// downloadObject downloads a single S3 object to a local file.
func downloadObject(ctx context.Context, client *s3.Client, bucket string, key string, localFilePath string) error {
	// Get the object from S3
	resp, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to get object %s: %v\n", key, err)
	}
	defer resp.Body.Close()

	// Create local directory if incase file doesn't exist
	localFilePathDir := filepath.Dir(localFilePath)

	if err := os.MkdirAll(localFilePathDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v\n", localFilePathDir, err)
	}

	// Create or truncate local file
	file, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v\n", localFilePath, err)
	}
	defer file.Close()

	// Copy object content to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %v\n", localFilePath, err)
	}

	fmt.Printf("Downloaded %s to %s\n", key, localFilePath)
	return nil
}
