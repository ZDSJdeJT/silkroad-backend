package cache

import (
	"fmt"
	"silkroad-backend/models"
	"silkroad-backend/utils"
	"strconv"
	"sync"
)

var Cache sync.Map

func LoadTextValue(key string) string {
	if value, ok := Cache.Load(key); ok {
		return fmt.Sprintf("%v", value)
	} else {
		return ""
	}
}

func LoadNumberValue(key string) uint64 {
	if value, ok := Cache.Load(key); ok {
		if intValue, err := strconv.Atoi(fmt.Sprintf("%v", value)); err == nil {
			return uint64(intValue)
		}
	}
	return 0
}

func InitCacheAndClientHTML(settings []models.Setting) error {
	for _, setting := range settings {
		err := HandleSetting(setting)
		if err != nil {
			return err
		}
	}
	return nil
}

func HandleSetting(setting models.Setting) error {
	switch setting.Key {
	case models.AdminName:
		Cache.Store(models.AdminName, setting.TextValue)
	case models.AdminPassword:
		Cache.Store(models.AdminPassword, setting.TextValue)
	case models.WebsiteDescription:
		clientHTML, err := utils.ReadClientHTML()
		if err != nil {
			return nil
		}
		clientHTML, err = utils.ReplaceHTMLMetaTag(clientHTML, utils.DescriptionMetaName, setting.TextValue)
		if err != nil {
			return err
		}
		err = utils.OverwriteClientHTML(clientHTML)
		if err != nil {
			return err
		}
	case models.WebsiteKeywords:
		clientHTML, err := utils.ReadClientHTML()
		if err != nil {
			return nil
		}
		clientHTML, err = utils.ReplaceHTMLMetaTag(clientHTML, utils.KeywordsMetaName, setting.TextValue)
		if err != nil {
			return err
		}
		err = utils.OverwriteClientHTML(clientHTML)
		if err != nil {
			return err
		}
	case models.KeepDays:
		Cache.Store(models.KeepDays, setting.NumberValue)
	case models.DownloadTimes:
		Cache.Store(models.DownloadTimes, setting.NumberValue)
	case models.UploadFileBytes:
		Cache.Store(models.UploadFileBytes, setting.NumberValue)
	case models.UploadChunkBytes:
		Cache.Store(models.UploadChunkBytes, setting.NumberValue)
	case models.UploadTextLength:
		Cache.Store(models.UploadTextLength, setting.NumberValue)
	}
	return nil
}
