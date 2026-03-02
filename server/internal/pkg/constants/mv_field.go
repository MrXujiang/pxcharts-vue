package constants

type FieldType string

const (
	MvFieldTypeText           FieldType = "text"
	MvFieldTypeRadio          FieldType = "radio"
	MvFieldTypeNumber         FieldType = "number"
	MvFieldTypeSelectMultiple FieldType = "selectMultiple"
	MvFieldTypeDate           FieldType = "date"
	MvFieldTypePerson         FieldType = "person"
	MvFieldTypeCurrency       FieldType = "currency"
	MvFieldTypeImage          FieldType = "image"
	MvFieldTypeAttachment     FieldType = "attachment"
	MvFieldTypeProgress       FieldType = "progress"
	MvFieldTypeLink           FieldType = "link"
	MvFieldTypeFormula        FieldType = "formula"
	MvFieldTypeCheckbox       FieldType = "checkbox"
	MvFieldTypeRate           FieldType = "rate"
	MvFieldTypeEmail          FieldType = "email"
	MvFieldTypeIdCard         FieldType = "idCard"
	MvFieldTypeRichText       FieldType = "richText"
	MvFieldTypeFindRef        FieldType = "findRef"
	MvFieldTypeCreator        FieldType = "creator"
	// MvFieldTypeSingleAssociation 单向引用 mv_record中对应保存recordId的列表
	MvFieldTypeSingleAssociation FieldType = "singleAssociation"
	MvFieldTypeDoubleAssociation FieldType = "doubleAssociation"
	MvFieldTypeCreateTime        FieldType = "createTime"
	MvFieldTypeAutoNumber        FieldType = "autoNumber"
)

// SingleAssociationConfig 单向引用的配置
type SingleAssociationConfig struct {
	TableSchemaID    string
	AllowMultiSelect bool
	Description      string
}

// DoubleAssociationConfig 单向引用的配置
type DoubleAssociationConfig struct {
	TableSchemaID    string
	AllowMultiSelect bool
	Description      string
}
