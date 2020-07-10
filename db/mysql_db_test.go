package db

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func Test_Mysql_DB(t *testing.T) {
	cfg := DefMysqlCfg()
	// 创建数据库
	assert.NoError(t, CreateDB(cfg))
	// 创建数据库连接
	db, err := NewMysql(cfg)
	assert.NoError(t, err)
	defer db.Close() // 关闭数据库连接

	// 数据库模型
	type User struct {
		gorm.Model
		Name string `gorm:"type:varchar(32);unique_index"`
	}
	// 删除表结构
	assert.NoError(t, db.DropTableIfExists(&User{}).Error)
	// 新建表结构
	assert.NoError(t, db.AutoMigrate(&User{}).Error)
	// 插入数据
	test1 := &User{Name: "test"}
	assert.NoError(t, db.Create(test1).Error)
	// 查询数据
	var test2 = User{}
	assert.NoError(t, db.First(&test2).Error)
	assert.Equal(t, test1.Name, test2.Name)
}
