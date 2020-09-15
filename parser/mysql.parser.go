package parser

import (
	"fmt"
	"github.com/just-coding-0/orm"
	"github.com/just-coding-0/orm/trick/camel"
	"github.com/just-coding-0/orm/trick/stringpool"
	"reflect"
	"strings"
	"sync"
)

type mysqlParserSyncPool struct {
	Pool sync.Pool
}

type mysqlParse struct {
	PrimaryKey       string
	command          Command // SELECT UPDATE DELETE
	dbName           string  // 数据库名
	tableName        string  // 表名
	body             string  // * OR  id,name .....
	values           string
	limit            uint64 // limit 10
	offset           uint64
	order            []string // ([]string{"id DESC","name DESC" })
	builder          strings.Builder
	characterSet     string // utf8
	characterCOLLATE string // utf8_general_ci
	engine           string // InnoDb
}

const (
	mysqlPrimaryKeyFormat      = "`%s` bigint NOT NULL AUTO_INCREMENT"
	mysqlPrimaryKeyIndexFormat = "PRIMARY KEY (`%s`) USING BTREE"
)

func (p *mysqlParse) CreateTable(model orm.Model) (string, error) {
	Val := reflect.ValueOf(model)
	Type := reflect.TypeOf(model)

	// 如果是指针类型,先解指针
	if Val.Kind() == reflect.Ptr {
		Val = Val.Elem()
		Type = Type.Elem()
	}

	p.command = CREATE
	p.tableName = model.TableName()
	if p.tableName == "" { // 驼峰转json格式
		p.tableName = camel.CamelToSnake(Type.Name())
	}

	// 使用sync pool
	SubFieldArr := stringpool.GetStringSlice(Type.NumField())

	for i := 0; i < Type.NumField(); i++ {
		SubField, err := p.ParseFieldToCreateTable(Type.Field(i))
		if err != nil {
			return "", err
		}
		SubFieldArr = append(SubFieldArr, SubField)
	}

	buildStr := p.Build(SubFieldArr)

	// 放回池中
	stringpool.PutStringSlice(SubFieldArr)

	return buildStr, nil
}
func (p *mysqlParse) Create(model orm.Model) (string, error) {
	p.command = INSERT
	Val := reflect.ValueOf(model)

	// 如果是指针类型,先解指针
	if Val.Kind() == reflect.Ptr {
		Val = Val.Elem()
	}

	p.tableName = model.TableName()
	if p.tableName == "" { // 驼峰转json格式
		p.tableName = camel.CamelToSnake(Val.Type().Name())
	}

	values := stringpool.GetStringSlice(Val.Type().NumField())
	filedNames := stringpool.GetStringSlice(Val.Type().NumField())

	// 尝试读取已反射的字段结构
	for i := 0; i < Val.Type().NumField(); i++ {

		if isBlank(Val.Field(i)) {
			continue
		}

		fieldType := Val.Type().Field(i)
		field := Val.Field(i)
		kind := field.Type().Kind()

		if !checkType(kind) { // 检查类型
			return "", UnSupportFieldType
		}
		if kind != reflect.Struct {
			filedNames = append(filedNames, camel.CamelToSnake(fieldType.Name))
			values = append(values, p.ParseFieldValue(field))
		} else { // 最多支持一层嵌套
			for i := 0; i < fieldType.Type.NumField(); i++ {
				filedNames = append(filedNames, camel.CamelToSnake(field.Type().Field(i).Name))
				values = append(values, p.ParseFieldValue(field.Field(i)))
			}
		}

	}

	p.body = strings.Join(filedNames, ",")
	p.values = fmt.Sprintf("(%s)", strings.Join(values, ","))

	stringpool.PutStringSlice(filedNames)
	stringpool.PutStringSlice(values)

	return p.Build(nil), nil
}
func (p *mysqlParse) CreateBySlice(module []orm.Model) (string, error) {
	p.command = INSERT
	return "", nil
}
func (p *mysqlParse) Update() {
	p.command = UPDATE
}
func (p *mysqlParse) Delete() {
	p.command = DELETE
}
func (p *mysqlParse) Select(query string) {
	p.body = query
}
func (p *mysqlParse) Count() {
	p.body = `COUNT(1)`
}
func (p *mysqlParse) Limit(limit uint64) {
	p.limit = limit
}
func (p *mysqlParse) Offset(offset uint64) {
	p.offset = offset
}
func (p *mysqlParse) SetDb(dbName string) {
	p.dbName = dbName
}
func (p *mysqlParse) Order(orderBy []string) {
	p.order = orderBy
}
func (p *mysqlParse) Build(args ...interface{}) string { // 构建query
	switch p.command {
	case CREATE:
		// 拼接btree
		SubFieldArr := args[0].([]string)
		SubFieldArr = append(SubFieldArr, fmt.Sprintf(mysqlPrimaryKeyIndexFormat, p.PrimaryKey))
		// 构建sql
		p.builder.WriteString(fmt.Sprintf("CREATE TABLE `%s`.`%s` (", p.dbName, p.tableName))
		// 字段
		p.builder.WriteString(strings.Join(SubFieldArr, ","))
		// 尾部
		p.builder.WriteString(fmt.Sprintf(") ENGINE = %s  CHARACTER SET = %s COLLATE = %s ROW_FORMAT = Dynamic;", p.engine, p.characterSet, p.characterCOLLATE))
	case INSERT:
		p.builder.WriteString(fmt.Sprintf("INSERT INTO %s.%s(%s) VALUES %s", p.dbName, p.tableName, p.body, p.values))
	case SELECT:
	case UPDATE:
	case DELETE:

	}

	return p.builder.String()
}
func (p *mysqlParse) Reset() { // reset函数
	p.command = SELECT
	p.dbName = p.dbName[:0]
	p.tableName = p.tableName[:0]
	p.body = p.body[:0]
	p.limit = 0
	p.offset = 0
	p.order = p.order[:0]
	p.builder.Reset()
	p.characterSet = "utf8"
	p.characterCOLLATE = "utf8_general_ci"
	p.engine = "InnoDb"
}
func (p *mysqlParse) ParseFieldToCreateTable(field reflect.StructField) (str string, err error) {
	if val, ok := field.Tag.Lookup(orm.Tag); ok && val == orm.Primary {
		if len(p.PrimaryKey) > 0 {
			return "", ManyPrimaryKeyError
		}
		if field.Type.Kind() != reflect.Uint64 {
			return "", PrimaryMustBeUint64
		}

		p.PrimaryKey = camel.CamelToSnake(field.Name)
		return fmt.Sprintf(mysqlPrimaryKeyFormat, p.PrimaryKey), nil
	}

	_default, ok := field.Tag.Lookup(orm.FieldDefaultTag)

	integerDefault := `0`
	charDefault := `''`
	dateTimeDefault := `CURRENT_TIMESTAMP(0)`
	dateDefault := `CURRENT_DATE(0)`
	fieldLength := `255`

	if ok {
		integerDefault = _default
		charDefault = _default
		dateTimeDefault = _default
		dateDefault = _default
	}

	kind := field.Type.Kind()
	// 简单的类型映射
	switch kind {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str = fmt.Sprintf("`%s` %s NOT NULL DEFAULT %s ", camel.CamelToSnake(field.Name), mysqlReflect(kind), integerDefault)
	case reflect.Float32, reflect.Float64:
		str = fmt.Sprintf("`%s` %s  NOT NULL DEFAULT %s ", camel.CamelToSnake(field.Name), mysqlReflect(kind), integerDefault)
		// price float64 `decimal:"5,3"` 前面五位 小数小3位
		if val, ok := field.Tag.Lookup(orm.DecimalTag); ok {
			str = fmt.Sprintf("`%s` decimal(%s)  NOT NULL DEFAULT %s ", camel.CamelToSnake(field.Name), val, integerDefault)
		}
	case reflect.String:
		_fieldLength, ok := field.Tag.Lookup(orm.FieldLengthTag)
		if ok { // 默认255字节
			fieldLength = _fieldLength
		}

		str = fmt.Sprintf("`%s` varchar(%s) CHARACTER SET %s COLLATE %s NOT NULL DEFAULT %s",
			camel.CamelToSnake(field.Name), fieldLength, p.characterSet, p.characterCOLLATE, charDefault)

		if fieldType, ok := field.Tag.Lookup(orm.TypeTag); ok {
			if !mysqlCheckCharFieldType(fieldType) { // 检查是否支持该自定义类型
				err = UnSupportFieldType
			} else if fieldType == date {
				str = fmt.Sprintf("`%s` date NOT NULL DEFAULT %s", camel.CamelToSnake(field.Name), dateDefault)
			} else if fieldType == datetime {
				str = fmt.Sprintf("`%s` datetime NOT NULL DEFAULT %s", camel.CamelToSnake(field.Name), dateTimeDefault)
			} else {
				str = fmt.Sprintf("`%s` %s CHARACTER SET %s COLLATE %s NOT NULL ",
					camel.CamelToSnake(field.Name), fieldType, p.characterSet, p.characterCOLLATE)
			}
		}
	case reflect.Struct:

		fieldArr := stringpool.GetStringSlice(field.Type.NumField())

		for i := 0; i < field.Type.NumField(); i++ {
			subfield, err := p.ParseFieldToCreateTable(field.Type.Field(i))
			if err != nil {
				return "", err
			}
			fieldArr = append(fieldArr, subfield)
		}

		str = strings.Join(fieldArr, ",")

		stringpool.PutStringSlice(fieldArr)
	case reflect.Ptr:
		_type := field.Type.Elem()

		fieldArr := stringpool.GetStringSlice(_type.NumField())

		for i := 0; i < _type.NumField(); i++ {
			subfield, err := p.ParseFieldToCreateTable(_type.Field(i))
			if err != nil {
				return "", err
			}
			fieldArr = append(fieldArr, subfield)
		}

		str = strings.Join(fieldArr, ",")
		stringpool.PutStringSlice(fieldArr)

	default:
		// 这里使用switch 最后也会被优化成if 就那么个东西。
		switch field.Type.String() {
		case "time.Time":
			if field.Tag.Get(orm.TypeTag) == date {
				str = fmt.Sprintf("`%s` %s NOT NULL DEFAULT %s", camel.CamelToSnake(field.Name), date, dateDefault)
			} else { // 默认Datetime
				str = fmt.Sprintf("`%s` %s NOT NULL DEFAULT %s", camel.CamelToSnake(field.Name), datetime, dateTimeDefault)
			}
		default:
			err = UnSupportFieldType
		}
	}
	return
}
func (p *mysqlParse) ParseFieldToInsert(fieldType reflect.StructField) string {
	return camel.CamelToSnake(fieldType.Name)
}
func (p *mysqlParse) ParseFieldValue(field reflect.Value) (str string) {
	kind := field.Kind()

	if kind <= reflect.Bool { // 布尔转成uint8
		str = fmt.Sprintf("%d", boolToUint8(field.Bool()))
	} else if kind <= reflect.Int64 {
		str = fmt.Sprintf("%d", field.Int())
	} else if kind <= reflect.Uint64 {
		str = fmt.Sprintf("%d", field.Uint())
	} else if kind <= reflect.Float64 {
		str = fmt.Sprintf("%f", field.Float())
	} else if kind == reflect.String || field.Type().Name() == "time.Time" { // 拼接字符串 或者是date datetime
		str = `'` + field.String() + `'`
	}

	return
}

