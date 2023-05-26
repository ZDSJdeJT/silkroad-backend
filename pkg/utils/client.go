package utils

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/fs"
	"os"
	"silkroad-backend/app/models"
	"strings"
)

func InitClientHTML() error {
	// 打开数据库连接
	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_DSN")), &gorm.Config{})
	if err != nil {
		return err
	}
	var websiteSettings []models.Setting
	db.Where("key IN ?", []string{"WEBSITE_TITLE", "WEBSITE_DESCRIPTION", "WEBSITE_KEYWORDS"}).Find(&websiteSettings)

	for _, setting := range websiteSettings {
		switch setting.Key {
		case "WEBSITE_TITLE":
			websiteName := setting.Value
			var memo map[string]interface{}
			if err := json.Unmarshal(websiteName, &memo); err != nil {
				return err
			}
			err = ReplaceClientHTMLTitle(memo["data"].(string))
			if err != nil {
				return err
			}
		case "WEBSITE_DESCRIPTION":
			websiteDescription := setting.Value
			var memo map[string]interface{}
			if err := json.Unmarshal(websiteDescription, &memo); err != nil {
				return err
			}
			err = ReplaceClientHTMLMetaDescription(memo["data"].(string))
			if err != nil {
				return err
			}
		case "WEBSITE_KEYWORDS":
			websiteKeywords := setting.Value
			var memo map[string]interface{}
			if err := json.Unmarshal(websiteKeywords, &memo); err != nil {
				return err
			}
			err = ReplaceClientHTMLMetaKeywords(memo["data"].(string))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ReplaceClientHTMLTitle(newTitle string) error {
	// 读取原始文件
	fileBytes, err := os.ReadFile("./client/index.html")
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	htmlString := string(fileBytes)

	// 替换 HTML 中的 title 标签内容
	var newHtmlString string
	token := html.NewTokenizer(strings.NewReader(htmlString))
	for {
		tt := token.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt == html.StartTagToken {
			tn, _ := token.TagName()
			if string(tn) == "title" {
				_ = token.Next()
				newHtmlString += fmt.Sprintf("<title>%s", newTitle)
				continue
			}
		}
		newHtmlString += string(token.Raw())
	}

	// 将修改后的 HTML 内容写回到原始文件中
	if err := os.WriteFile("./client/index.html", []byte(newHtmlString), fs.ModePerm); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

func ReplaceClientHTMLMetaDescription(newContent string) error {
	// 读取原始文件
	fileBytes, err := os.ReadFile("./client/index.html")
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	htmlString := string(fileBytes)

	// 替换 HTML 中的 meta 标签内容
	var newHtmlString string
	token := html.NewTokenizer(strings.NewReader(htmlString))
	for {
		tt := token.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt == html.SelfClosingTagToken || tt == html.StartTagToken {
			tn, hasAttr := token.TagName()
			if string(tn) == "meta" && hasAttr {
				key, val, _ := token.TagAttr()
				if string(key) == "name" && string(val) == "description" {
					newHtmlString += fmt.Sprintf(`<meta name="description" content="%s"/>`, newContent)
					continue
				}
			}
		}
		newHtmlString += string(token.Raw())
	}

	// 将修改后的 HTML 内容写回到原始文件中
	if err := os.WriteFile("./client/index.html", []byte(newHtmlString), fs.ModePerm); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

func ReplaceClientHTMLMetaKeywords(newContent string) error {
	// 读取原始文件
	fileBytes, err := os.ReadFile("./client/index.html")
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	htmlString := string(fileBytes)

	// 替换 HTML 中的 meta 标签内容
	var newHtmlString string
	token := html.NewTokenizer(strings.NewReader(htmlString))
	for {
		tt := token.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt == html.SelfClosingTagToken || tt == html.StartTagToken {
			tn, hasAttr := token.TagName()
			if string(tn) == "meta" && hasAttr {
				key, val, _ := token.TagAttr()
				if string(key) == "name" && string(val) == "keywords" {
					newHtmlString += fmt.Sprintf(`<meta name="keywords" content="%s"/>`, newContent)
					continue
				}
			}
		}
		newHtmlString += string(token.Raw())
	}

	// 将修改后的 HTML 内容写回到原始文件中
	if err := os.WriteFile("./client/index.html", []byte(newHtmlString), fs.ModePerm); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
