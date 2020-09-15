package parser

import (
	"github.com/just-coding-0/orm"
	"reflect"
)

type Factory interface {
	Get() StringBuilder
	Put(builder StringBuilder)
}

type StringBuilder interface {
	SetDb(db string)
	CreateTable(model orm.Model) (string, error)
	Create(module orm.Model) (string, error)
	CreateBySlice(module []orm.Model) (string, error)
	Update()
	Delete()
	Select(query string)
	Count()
	Limit(limit uint64)
	Offset(offset uint64)
	Order(orderBy []string)
	Build(...interface{}) string
	Reset()
}

func NewMysqlFactory() Factory {
	return newMysqlParserSyncPool()
}

type Command uint8

const (
	CREATE Command = iota + 1
	INSERT
	SELECT
	DELETE
	UPDATE
)

func checkType(kind reflect.Kind) bool {
	switch kind {
	case reflect.Array:
		return false
	case reflect.Chan:
		return false
	case reflect.Func:
		return false
	case reflect.Interface:
		return false
	case reflect.Map:
		return false
	case reflect.Ptr:
		return false
	}
	return true
}
