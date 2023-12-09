package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
	"html/template"
	"log"
	"strings"
	"time"
)

func GetS3Session() (*session.Session, error) {
	// AWSセッションの初期化
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("auto"),
		Endpoint:    aws.String(s3Endpoint),
		Credentials: credentials.NewStaticCredentials(s3AccessToken, s3SecretKey, ""),
	})

	// S3サービスクライアントの作成
	return sess, err
}

func serveTs(c *fiber.Ctx) error {
	day := c.Params("day")
	log.Println("day: " + day)

	key := c.Params("key")
	log.Println("key: " + key)

	// s3 から presigned url を取得
	// presigned url からファイルを取得
	sess, _ := GetS3Session()
	s3Client := s3.New(sess)

	req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(day + "/" + key),
	})

	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		log.Println("Failed to sign request", err)
		return err
	}

	c.Redirect(urlStr, 302)

	//// ファイルを開く
	//file, err := s3Client.GetObject(&s3.GetObjectInput{
	//	Bucket: aws.String(s3BucketName),
	//	Key:    aws.String(fileName),
	//})
	//if err != nil {
	//	return err
	//}
	//
	//defer file.Body.Close()
	//
	//bytes_, err := io.ReadAll(file.Body)
	//if err != nil {
	//	return err
	//}
	//c.Set("Content-Type", "video/mp2t")
	//c.Send(bytes_)
	return nil
}

func serveHLSContent(c *fiber.Ctx) error {
	cameraID := c.Params("cameraID")
	day := c.Params("day")

	// m3u8 txt をテンプレートとして作成

	// 0..59
	//var minutes []int
	//for i := 0; i < 60; i++ {
	//	minutes = append(minutes, i)
	//}
	minutes := []int{19, 20, 21, 22, 23}
	//var hours []int
	//for i := 0; i < 24; i++ {
	//	hours = append(hours, i)
	//}
	hours := []int{15}

	data := map[string]interface{}{
		"cameraID": cameraID,
		"day":      day,
		"hours":    hours,
		"minutes":  minutes,
	}

	w := new(strings.Builder)
	tpl := template.Must(template.ParseFiles("m3u8.txt"))
	err := tpl.Execute(w, data)
	if err != nil {
		log.Println(err)
		return err
	}

	txt := w.String()

	c.Set("Content-Type", "application/x-mpegURL")
	c.SendString(txt)
	return nil
}
