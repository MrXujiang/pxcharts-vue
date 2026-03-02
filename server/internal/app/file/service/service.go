package service

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mime/multipart"
	"mvtable/internal/app/file/model"
	model2 "mvtable/internal/app/user/model"
	model3 "mvtable/internal/app/website_config/model"
	"mvtable/internal/pkg/constants"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"
	"mvtable/pkg/oss"
	"slices"
)

type FileService struct{}

func NewFileService() *FileService {
	return &FileService{}
}

func (s *FileService) Upload(userID string, file *multipart.FileHeader, remark string) (*model.UploadRes, error) {
	stream, err := file.Open()
	if err != nil {
		log.Error("open file error: " + err.Error())
		return nil, errorx.InternalServerError("文件解析失败")
	}
	defer stream.Close()

	genUUID := uuid.New().String()
	path := constants.UploadFileDir + "/" + genUUID + "_" + file.Filename

	ossProvider := oss.GetClient().Provider
	// 获取网站配置
	websiteConfigList, _, err := db.List[model3.WebsiteConfig](db.GetDB(), 1, 1, nil, nil)
	if err != nil {
		log.Error("query website config error: " + err.Error())
		return nil, errorx.InternalServerError("文件上传失败")
	}
	if len(websiteConfigList) != 0 {
		config := websiteConfigList[0]
		var schema model3.WebsiteConfigSchema
		if err = json.Unmarshal(config.Schema, &schema); err != nil {
			log.Error("unmarshal website config schema error: " + err.Error())
			return nil, errorx.InternalServerError("文件上传失败")
		}
		if schema.UploadMethod != "" {
			ossProvider = schema.UploadMethod
		}
	}

	ossAccessURL := oss.GetClient().MinioConf.AccessURL
	createFile := &model.File{
		UserID:   userID,
		Filename: file.Filename,
		Filesize: file.Size,
		Filetype: file.Header.Get("Content-Type"),
		OSS:      ossProvider,
		Path:     path,
		Remark:   remark,
	}

	if err = db.GetDB().Transaction(func(tx *gorm.DB) error {
		if err = db.Create(tx, createFile); err != nil {
			return err
		}

		file.Filename = genUUID + "_" + file.Filename

		if ossProvider == constants.OSSMinio {
			ossAccessURL = oss.GetClient().MinioConf.AccessURL
			return oss.GetClient().Minio.UploadFile(file, constants.UploadFileDir)
		}
		if ossProvider == constants.OSSAliyun {
			ossAccessURL = oss.GetClient().AliyunConf.AccessURL
			return oss.GetClient().Aliyun.UploadFile(file, constants.UploadFileDir)
		}
		return fmt.Errorf("unknown oss type: %s", ossProvider)
	}); err != nil {
		log.Error("upload file error: " + err.Error())
		return nil, errorx.InternalServerError("文件上传失败")
	}

	return &model.UploadRes{
		ID:       createFile.ID,
		URL:      ossAccessURL + "/" + createFile.Path,
		Filename: createFile.Filename,
		Size:     createFile.Filesize,
	}, nil
}

func (s *FileService) AdminGetList(req *model.AdminGetListReq) (*model.AdminGetListRes, error) {
	conditions := make(map[string]any)
	if req.Filename != "" {
		conditions["filename"] = []any{"LIKE", "%" + req.Filename + "%"}
	}
	if req.OSS != "" {
		conditions["oss"] = req.OSS
	}
	if req.Remark != "" {
		conditions["remark"] = []any{"LIKE", "%" + req.Remark + "%"}
	}

	list, total, err := db.List[model.File](db.GetDB(), req.Page, req.Size, conditions, nil)
	if err != nil {
		log.Error("admin get file list error: " + err.Error())
		return nil, errorx.InternalServerError("获取文件列表失败")
	}

	var (
		userIds = make([]string, 0, len(list))
		userMap = make(map[string]string)
		resList = make([]*model.AdminGetListItem, len(list))
	)
	for _, v := range list {
		if !slices.Contains(userIds, v.UserID) {
			userIds = append(userIds, v.UserID)
		}
	}

	users, _, err := db.List[model2.User](db.GetDB(), 0, 0, map[string]any{"id": []any{"IN", userIds}}, nil)
	if err != nil {
		log.Error("query user list error: " + err.Error())
		return nil, errorx.InternalServerError("获取文件列表失败")
	}

	for _, v := range users {
		userMap[v.ID] = v.Nickname
	}

	for i, v := range list {
		accessURL := ""
		if v.OSS == constants.OSSMinio {
			accessURL = oss.GetClient().MinioConf.AccessURL
		} else {
			accessURL = oss.GetClient().AliyunConf.AccessURL
		}
		resList[i] = &model.AdminGetListItem{
			File:     *v,
			Uploader: userMap[v.UserID],
			URL:      accessURL + "/" + v.Path,
		}
	}

	return &model.AdminGetListRes{
		List:  resList,
		Total: total,
	}, nil
}
