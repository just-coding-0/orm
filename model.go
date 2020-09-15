package orm

type Model interface {
	TableName() string
}

type BaseObject struct {
	Id        uint64 `orm:"primary"`
	CreatedAt string `type:"datetime" default:"CURRENT_TIMESTAMP(0)"`
	UpdatedAt string `type:"datetime" default:"CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0)"`
}

const (
	Tag     = `orm`
	Primary = `primary`

	TypeTag        = `type`
	FieldLengthTag = `length`
	DecimalTag     = `decimal`
	FieldDefaultTag   = `default`
)

// 如果TableName为空,则按照struct的Name建立Table(驼峰)
func (t BaseObject) TableName() string {

	return ""
}

