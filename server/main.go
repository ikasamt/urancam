package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024, // 100MBに設定
	})
	app.Get("/ping", func(c *fiber.Ctx) error {
		c.JSON("OK")
		return nil
	})
	app.Get("/", func(c *fiber.Ctx) error {
		data := map[string]interface{}{
			"cameraIDs": []string{"1", "2", "3"},
		}
		w := new(strings.Builder)
		tpl := template.Must(template.ParseFiles("index.html"))
		err := tpl.Execute(w, data)
		if err != nil {
			log.Println(err)
			return err
		}

		txt := w.String()
		c.Set("Content-Type", "text/html")
		c.SendString(txt)
		return nil
	})
	app.Get("/m3u8/:cameraID/:day", serveHLSContent)
	app.Get("/ts/:key", serveTs)
	app.Post("/upload", uploadHandler)
	app.Listen(":8080")
}

func uploadHandler(c *fiber.Ctx) error {
	// ファイルを取得
	file, err := c.FormFile("file")
	if err != nil {
		log.Println("File upload error: " + err.Error())
		return c.JSON(err)
	}

	// day
	cameraId := c.FormValue("cameraID")

	day := c.FormValue("day")
	if day == "" {
		log.Println("day", day)
		return c.JSON(fmt.Errorf("day error  empty "))
	}

	hour := c.FormValue("hour")
	if hour == "" {
		log.Println("hour", hour)
		return c.JSON(fmt.Errorf("hour error  empty "))
	}

	minute := file.Filename

	// ファイルをサーバーに一時保存
	destFileName := fmt.Sprintf("%s-%s-%s-%s", cameraId, day, hour, minute)
	tempFilePath := filepath.Join(os.TempDir(), destFileName)
	if err := c.SaveFile(file, tempFilePath); err != nil {
		log.Println("File save error: " + err.Error())
		return c.JSON(err)
	}

	// HLS形式に変換（仮の関数、実装が必要）
	var fns []string
	if fns, err = convertToHLS(tempFilePath); err != nil {
		log.Println("File convert error: " + err.Error())
		return c.JSON(err)
	}

	// 生成したファイルをS3にアップロード（仮の関数、実装が必要）
	uploadFilesConcurrently(fns)
	return c.SendString("File uploaded successfully")
}
