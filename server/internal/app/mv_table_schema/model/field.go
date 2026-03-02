package model

import (
	"mvtable/internal/pkg/constants"
)

type GeneralFieldConfig interface {
	GetBase() *FieldBase
	SetBase(v *FieldBase)
}

// BaseField 基础字段接口
type BaseField interface {
	GetBase() *FieldBase
	SetBase(v *FieldBase)
}

type FieldBase struct {
	ID       string              `json:"id"`
	Type     constants.FieldType `json:"type"`
	Title    string              `json:"title"`
	Width    int                 `json:"width"`
	Label    string              `json:"label"`
	Desc     string              `json:"desc"`
	IconName string              `json:"iconName"`
	Fixed    bool                `json:"fixed"`
	IsShow   bool                `json:"isShow"`
}

func (f *FieldBase) GetBase() *FieldBase {
	return f
}

func (f *FieldBase) SetBase(v *FieldBase) {
	f.ID = v.ID
	f.Type = v.Type
	f.Title = v.Title
	f.Width = v.Width
	f.Label = v.Label
	f.Desc = v.Desc
	f.IconName = v.IconName
	f.Fixed = v.Fixed
	f.IsShow = v.IsShow
}

// FieldTextConfig 文本字段配置
type FieldTextConfig struct {
	FieldBase
	DefaultValue string `json:"defaultValue"`
}

// FieldRadioConfig 单选字段配置
type FieldRadioConfig struct {
	FieldBase
	DefaultValue string             `json:"defaultValue"`
	Options      []RadioOption      `json:"options"`
	Settings     FieldRadioSettings `json:"settings"`
	Ref          FieldRef           `json:"ref"`
}

// FieldNumberConfig 数字字段配置
type FieldNumberConfig struct {
	FieldBase
	DefaultValue string              `json:"defaultValue"`
	Settings     FieldNumberSettings `json:"settings"`
}

// FieldSelectMultipleConfig 多选字段配置
type FieldSelectMultipleConfig struct {
	FieldBase
	DefaultValue []string           `json:"defaultValue"`
	Options      []RadioOption      `json:"options"`
	Settings     FieldRadioSettings `json:"settings"`
	Ref          FieldRef           `json:"ref"`
}

// FieldDateConfig 日期字段配置
type FieldDateConfig struct {
	FieldBase
	DefaultValue string            `json:"defaultValue"`
	Format       string            `json:"format"`
	Settings     FieldDateSettings `json:"settings"`
}

// FieldPersonConfig 人员字段配置
type FieldPersonConfig struct {
	FieldBase
	DefaultValue []any               `json:"defaultValue"`
	Settings     FieldPersonSettings `json:"settings"`
}

// FieldCurrencyConfig 货币字段配置
type FieldCurrencyConfig struct {
	FieldBase
	DefaultValue string                `json:"defaultValue"`
	Settings     FieldCurrencySettings `json:"settings"`
}

// FieldImageConfig 图片字段配置
type FieldImageConfig struct {
	FieldBase
	DefaultValue []any `json:"defaultValue"`
}

// FieldAttachmentConfig 附件字段配置
type FieldAttachmentConfig struct {
	FieldBase
}

// FieldProgressConfig 进度字段配置
type FieldProgressConfig struct {
	FieldBase
	DefaultValue any                   `json:"defaultValue"`
	Settings     FieldProgressSettings `json:"settings"`
}

// FieldLinkConfig 链接字段配置
type FieldLinkConfig struct {
	FieldBase
	DefaultValue FieldLinkDefaultValue `json:"defaultValue"`
	LinkTitle    string                `json:"linkTitle"`
	Settings     FieldLinkSettings     `json:"settings"`
}

// FieldFormulaConfig 公式字段配置
type FieldFormulaConfig struct {
	FieldBase
	Expression string `json:"expression"` // 公式表达式
}

// FieldCheckboxConfig 复选框字段配置
type FieldCheckboxConfig struct {
	FieldBase
	DefaultValue int `json:"defaultValue"`
}

// FieldRateConfig 评分字段配置
type FieldRateConfig struct {
	FieldBase
	DefaultValue any               `json:"defaultValue"`
	Settings     FieldRateSettings `json:"settings"`
}

// FieldEmailConfig 邮箱字段配置
type FieldEmailConfig struct {
	FieldBase
	DefaultValue string `json:"defaultValue"`
}

// FieldIdCardConfig 身份证字段配置
type FieldIdCardConfig struct {
	FieldBase
	DefaultValue string `json:"defaultValue"`
}

// FieldRichTextConfig 富文本字段配置
type FieldRichTextConfig struct {
	FieldBase
	DefaultValue string `json:"defaultValue"`
}

// FieldFindRefConfig 查找引用字段配置
type FieldFindRefConfig struct {
	FieldBase
	DefaultValue string               `json:"defaultValue"`
	Settings     FieldFindRefSettings `json:"settings"`
}

// FieldCreatorConfig 创建人字段配置
type FieldCreatorConfig struct {
	FieldBase
	DefaultValue any                  `json:"defaultValue"`
	Settings     FieldCreatorSettings `json:"settings"`
}

// FieldSingleAssociationConfig 单向关联字段配置
type FieldSingleAssociationConfig struct {
	FieldBase
	Settings FieldAssociationSettings `json:"settings"`
}

// FieldDoubleAssociationConfig 双向关联字段配置
type FieldDoubleAssociationConfig struct {
	FieldBase
	Settings FieldAssociationSettings `json:"settings"`
}

// FieldCreateTimeConfig 创建时间字段配置
type FieldCreateTimeConfig struct {
	FieldBase
	DefaultValue string `json:"defaultValue"`
	Format       string `json:"format"`
}

// FieldAutoNumberConfig 自动编号字段配置
type FieldAutoNumberConfig struct {
	FieldBase
	DefaultValue string                  `json:"defaultValue"`
	Settings     FieldAutoNumberSettings `json:"settings"`
}
