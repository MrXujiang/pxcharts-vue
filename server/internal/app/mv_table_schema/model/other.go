package model

type RadioOption struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Color string `json:"color"`
}

// FieldRef 字段引用结构体
type FieldRef struct {
	TableSchemaID string `json:"tableSchemaId"`
	FieldID       string `json:"fieldId"`
}

// FieldLinkDefaultValue 链接字段默认值结构体
type FieldLinkDefaultValue struct {
	LinkTitle string `json:"linkTitle"`
	LinkUrl   string `json:"linkUrl"`
}

// FieldSourceRef 字段源引用结构体
type FieldSourceRef struct {
	TableSchemaID string `json:"tableSchemaId"`
	FieldID       string `json:"fieldId"`
}

// FieldFindCondition 查找条件结构体
type FieldFindCondition struct {
	SourceTableFieldID string `json:"sourceTableFieldId"`
	Operator           string `json:"operator"`
	TargetTableFieldID string `json:"targetTableFieldId"`
}

// FieldAssociationSourceRef 关联字段源引用结构体
type FieldAssociationSourceRef struct {
	TableSchemaID string `json:"tableSchemaId"`
}

// FieldAssociationConditions 关联条件结构体
type FieldAssociationConditions struct {
	AllowMultiSelect bool `json:"allowMultiSelect"`
}

// FieldAutoNumberRules 自动编号规则结构体
type FieldAutoNumberRules struct {
	Number int    `json:"number"`
	Text   string `json:"text"`
}

// Settings 结构体定义

// FieldRadioSettings 单选字段设置
type FieldRadioSettings struct {
	IsReferenced bool `json:"isReferenced"`
}

// FieldNumberSettings 数字字段设置
type FieldNumberSettings struct {
	UseThousandSeparator    bool   `json:"useThousandSeparator"`
	DisplayFormat           string `json:"displayFormat"`
	DecimalPlaces           int    `json:"decimalPlaces"`
	ThousandSeparator       string `json:"thousandSeparator"`
	LargeNumberAbbreviation string `json:"largeNumberAbbreviation"`
	DisallowNegative        bool   `json:"disallowNegative"`
}

// FieldDateSettings 日期字段设置
type FieldDateSettings struct {
	IsDefaultCreateTime bool `json:"isDefaultCreateTime"`
}

// FieldPersonSettings 人员字段设置
type FieldPersonSettings struct {
	ShowNickname    bool `json:"showNickname"`
	ShowExternalTag bool `json:"showExternalTag"`
	AllowMultiple   bool `json:"allowMultiple"`
}

// FieldCurrencySettings 货币字段设置
type FieldCurrencySettings struct {
	UseThousandSeparator    bool   `json:"useThousandSeparator"`
	Currency                string `json:"currency"`
	DecimalPlaces           int    `json:"decimalPlaces"`
	ThousandSeparator       string `json:"thousandSeparator"`
	LargeNumberAbbreviation string `json:"largeNumberAbbreviation"`
	DisallowNegative        bool   `json:"disallowNegative"`
}

// FieldProgressSettings 进度字段设置
type FieldProgressSettings struct {
	NumberFormat         string `json:"numberFormat"`
	DecimalPlaces        int    `json:"decimalPlaces"`
	Color                string `json:"color"`
	EnableCustomProgress bool   `json:"enableCustomProgress"`
}

// FieldLinkSettings 链接字段设置
type FieldLinkSettings struct {
	IsLinkTitle bool `json:"isLinkTitle"`
}

// FieldRateSettings 评分字段设置
type FieldRateSettings struct {
	Icon       string `json:"icon"`
	ScoreStart int    `json:"scoreStart"`
	ScoreEnd   int    `json:"scoreEnd"`
}

// FieldCreatorSettings 创建人字段设置
type FieldCreatorSettings struct {
	ShowNickname    bool `json:"showNickname"`
	ShowExternalTag bool `json:"showExternalTag"`
}

// FieldFindRefSettings 查找引用字段设置
type FieldFindRefSettings struct {
	SourceRef      FieldSourceRef       `json:"sourceRef"`
	MatchCondition string               `json:"matchCondition"`
	Conditions     []FieldFindCondition `json:"conditions"`
}

// FieldAssociationSettings 关联字段设置
type FieldAssociationSettings struct {
	SourceRef  FieldAssociationSourceRef  `json:"sourceRef"`
	Conditions FieldAssociationConditions `json:"conditions"`
}

// FieldAutoNumberSettings 自动编号字段设置
type FieldAutoNumberSettings struct {
	IDType string               `json:"idType"`
	Rules  FieldAutoNumberRules `json:"rules"`
}
