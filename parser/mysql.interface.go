package parser

import (
	"reflect"
)

const (
	char     = "char"
	varChar  = "varchar"
	longText = "longtext"
	date     = "date"
	datetime = "datetime"
)

func mysqlCheckCharFieldType(t string) bool {
	switch t {
	case char:
		return true
	case varChar:
		return true
	case longText:
		return true
	case date:
		return  true
	case datetime:
		return true
	}
	return false
}

func mysqlReflect(kind reflect.Kind) string {
	switch kind {
	case reflect.Bool:
		return "tinyint unsigned"
	case reflect.Int:
		return "integer"
	case reflect.Int8:
		return "tinyint"
	case reflect.Int16:
		return "smallint"
	case reflect.Int32:
		return "integer"
	case reflect.Int64:
		return "bigint"
	case reflect.Uint:
		return "integer unsigned"
	case reflect.Uint8:
		return "tinyint unsigned"
	case reflect.Uint16:
		return "smallint unsigned"
	case reflect.Uint32:
		return "integer unsigned"
	case reflect.Uint64:
		return "bigint unsigned"
	case reflect.Float32, reflect.Float64:
		return "double"
	default:
		return ""
	}

}

func boolToUint8(bl bool) uint8 {
	if bl {
		return 1
	}
	return 0
}
