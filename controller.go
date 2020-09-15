package orm

import "database/sql"

// 接口定义
type controller interface {
	CreateTable(Model) (sql.Result, error)     // 建表  db.CreateTable(object)  		  直接解析对象,拼接query
	Create(Model) (sql.Result, error)          // 插入  db.Create(object)			  拼接query,插入数据库
	CreateBySlice([]Model) (sql.Result, error) // 插入  db.CreateBySlice([]object)     拼接query,插入数据库
	Delete(Model, uint64) (sql.Result, error)  // 删除  db.delete(object,primaryId)    拼接uqery,删除数据

	Find(Model) controller          // 查询  db.find(object).where().Execute()
	FindBySlice([]Model) controller // 查询 db.FindBySlice([]object).where().Execute()
	Update(column string, value interface{}) error
	// 更新  db.model(object).update(column,value).Execute()
	// db.module(object).update(column,value).where().Execute() 更新多条

	Updates(map[string]interface{}) error // 更新  db.model(object).update(map[string]interface{}{column:value})
	Count(count uint64) error             // count db.module(object).count()  db.module(object).where.count()
	Where() controller                    // 条件  db.model
	Model(Model) controller
	Limit(uint64) controller
	Offset(uint64) controller
	Execute() (sql.Result, error)
	Row() (*sql.Row, error)
	Rows() (*sql.Rows, error)
}
