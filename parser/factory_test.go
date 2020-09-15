package parser

import (
	"fmt"
	"github.com/just-coding-0/orm"
	"testing"
)

func TestNewMysqlParserSyncPool(b *testing.T) {

	p := NewMysqlFactory()


	type s struct {
		orm.BaseObject
		Name string
		Age  uint64
	}
	builder := p.Get()
	builder.SetDb("test")

	var r s

	r.Name = "阿里郎"
	r.Age = 15

	buf, _ := builder.Create(r)
	fmt.Println(buf)


	p.Put(builder)
	builder = p.Get()
	builder.SetDb("test")

	buf,_ = builder.CreateTable(r)
	fmt.Println(buf)

}
