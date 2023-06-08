package utils

import (
	"fmt"
	"golang.org/x/net/html"
	"io/fs"
	"os"
	"strings"
)

const ClientHTMLPath = "./client/index.html"

const (
	DescriptionMetaName = "description"
	KeywordsMetaName    = "keywords"
)

func ReplaceHTMLMetaTag(htmlString, metaName, newContent string) (string, error) {
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
				if string(key) == "name" && string(val) == metaName {
					newHtmlString += fmt.Sprintf(`<meta name="%s" content="%s"/>`, metaName, newContent)
					continue
				}
			}
		}
		newHtmlString += string(token.Raw())
	}
	return newHtmlString, nil
}

func ReadClientHTML() (string, error) {
	// 读取原始文件
	fileBytes, err := os.ReadFile(ClientHTMLPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}
	return string(fileBytes), nil
}

func OverwriteClientHTML(newHTMLString string) error {
	// 将修改后的 HTML 内容写回到原始文件中
	if err := os.WriteFile("./client/index.html", []byte(newHTMLString), fs.ModePerm); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	return nil
}
