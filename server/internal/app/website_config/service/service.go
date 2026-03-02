package service

import (
	"encoding/json"
	"mvtable/internal/app/website_config/model"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"
)

type WebsiteConfigService struct{}

func NewWebsiteConfigService() *WebsiteConfigService {
	return &WebsiteConfigService{}
}

func (s *WebsiteConfigService) Get() (any, error) {
	list, _, err := db.List[model.WebsiteConfig](db.GetDB(), 1, 1, nil, nil)
	if err != nil {
		log.Error("get website config error: " + err.Error())
		return nil, errorx.InternalServerError("获取失败")
	}

	var defaultWebsiteConfig model.WebsiteConfig
	if len(list) == 0 {
		defaultWebsiteConfig = model.WebsiteConfig{}
		schema := map[string]any{
			"logo":            "",
			"slogan":          "",
			"registerMethod":  []string{"inviteCode"},
			"communityRQCode": "",
			"inviteRQCode":    "",
			"copyright":       "",
			"uploadMethod":    "aliyun",
		}
		bytes, _ := json.Marshal(schema)
		defaultWebsiteConfig.Schema = bytes
		if err = db.Create(db.GetDB(), &defaultWebsiteConfig); err != nil {
			log.Error("create default website config error: " + err.Error())
			return nil, errorx.InternalServerError("获取失败")
		}

		return defaultWebsiteConfig.Schema, nil
	}

	return list[0].Schema, nil
}

func (s *WebsiteConfigService) AdminSave(req *model.WebsiteConfigSchema) error {
	bytes, _ := json.Marshal(req)

	if err := db.Create(db.GetDB(), &model.WebsiteConfig{
		Schema: bytes,
	}); err != nil {
		log.Error("save website config error: " + err.Error())
		return errorx.InternalServerError("保存失败")
	}

	return nil
}
