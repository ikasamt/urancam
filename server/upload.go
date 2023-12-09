package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var s3BucketName = os.Getenv("S3_BUCKET_NAME")
var s3Endpoint = os.Getenv("S3_ENDPOINT")
var s3AccessToken = os.Getenv("S3_ACCESS_TOKEN")
var s3SecretKey = os.Getenv("S3_SECRET_KEY")

func uploadToS3(filePath string) error {
	sess, _ := GetS3Session()
	s3Client := s3.New(sess)

	// ファイルを開く
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// S3にファイルをアップロード
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(filepath.Base(filePath)),
		Body:   file,
	})
	return err
}

func uploadFilesConcurrently(fns []string) []error {
	var wg sync.WaitGroup
	errs := make([]error, len(fns))

	for i, fn := range fns {
		wg.Add(1)
		go func(i int, fn string) {
			defer wg.Done()
			if err := uploadToS3(fn); err != nil {
				log.Println("File upload error: " + err.Error())
				errs[i] = err
			}
		}(i, fn)
	}

	wg.Wait()
	return errs
}
