package parser

var (
	ModelTypeError      = NewError("必须为指针类型")
	UnSupportFieldType  = NewError("不支持的字段类型")
	ManyPrimaryKeyError = NewError("只能有一个主键")
	PrimaryMustBeUint64 = NewError("主键必须为uint64")
)

type Error struct {
	msg string
}

func (e *Error) Error() string {
	return e.msg
}

func NewError(msg string) error {
	return &Error{msg}
}
