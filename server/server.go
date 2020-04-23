package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/zhcppy/golib/common/ginn"
	"github.com/zhcppy/golib/common/util"
	"github.com/zhcppy/golib/db"
	"github.com/zhcppy/golib/logger"
	"github.com/zhcppy/golib/module"
	"github.com/zhcppy/golib/module/user"
)

type Config struct {
	Name     string     `json:"name"`
	DbConfig *db.Config `json:"dbConfig"`
}

func (c *Config) String() string {
	marshal, _ := json.Marshal(c)
	return string(marshal)
}

func NewServer(cfg *Config) (*gin.Engine, error) {
	log := logger.New(cfg.Name)

	mysql, err := db.NewMysql(cfg.DbConfig)
	if err != nil {
		return nil, err
	}

	engine := gin.Default()
	if util.RunEnv() == util.EnvProd {
		gin.SetMode(gin.ReleaseMode)
	}
	engine.Use(ginn.Cors())
	group := engine.Group("/api/v1")

	sv := module.New(mysql, log)

	manager := module.NewManager(
		user.New(sv), // 注册不同service的路由
	)
	manager.BindAll(group)
	return engine, nil
}

// 创建数据库并创建表结构
func InitDB(cfg *db.Config) error {
	err := db.CreateDB(cfg)
	if err != nil {
		return err
	}
	mysql, err := db.NewMysql(cfg)
	if err != nil {
		return err
	}
	defer mysql.Close()
	return mysql.AutoMigrate(&user.User{}).Error
}
