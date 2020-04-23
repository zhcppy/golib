package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Mysql_DB(t *testing.T) {
	cfg := DefaultCfg()
	// 创建数据库
	assert.NoError(t, CreateDB(cfg))
	// 创建数据库连接
	db, err := NewMysql(cfg)
	assert.NoError(t, err)
	defer db.Close()               // 关闭数据库连接
	defer DropDB(db, cfg.Database) // 删除数据库
	defer ClearAllData(db)         // 删除所有表中的数据

	// 数据库模型
	type tester struct {
		Name string `gorm:"type:varchar(32);unique_index"`
	}
	// 删除表结构
	assert.NoError(t, db.DropTableIfExists(&tester{}).Error)
	// 新建表结构
	assert.NoError(t, db.AutoMigrate(&tester{}).Error)
	// 插入数据
	test1 := &tester{Name: "test"}
	assert.NoError(t, db.Create(test1).Error)
	// 查询数据
	var test2 = tester{}
	assert.NoError(t, db.First(&test2).Error)
	assert.Equal(t, test1.Name, test2.Name)
}
