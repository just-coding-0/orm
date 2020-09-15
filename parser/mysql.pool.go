package parser

import (
	"strings"
	"sync"
)

func newMysqlParse() StringBuilder {
	return &mysqlParse{
		command:          SELECT,
		builder:          strings.Builder{},
	}
}

// 使用前初始化
func (p *mysqlParserSyncPool) Get() StringBuilder {
	builder := p.Pool.Get().(StringBuilder)
	builder.Reset()
	return builder
}
func (p *mysqlParserSyncPool) Put(builder StringBuilder) {
	p.Pool.Put(builder)
}
func newMysqlParserSyncPool() Factory {
	pool := &mysqlParserSyncPool{}
	pool.Pool = sync.Pool{}
	pool.Pool.New = func() interface{} {
		return newMysqlParse()
	}
	return pool
}
