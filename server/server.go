package server

import (
	"encoding/json"

	"github.com/zhcppy/golib/module/user"

	"github.com/gin-gonic/gin"
	"github.com/zhcppy/golib/common/util"
	"github.com/zhcppy/golib/db"
	"github.com/zhcppy/golib/logger"
	"github.com/zhcppy/golib/module"
	"github.com/zhcppy/golib/module/middle"
)

type Config struct {
	DbConfig *db.Config `json:"dbConfig"`
}

func (c *Config) String() string {
	marshal, _ := json.Marshal(c)
	return string(marshal)
}
func NewServer(name string, cfg *Config) (*gin.Engine, error) {
	log := logger.New(name)

	mysql, err := db.NewMysql(cfg.DbConfig, func() []interface{} {
		return []interface{}{&user.User{}}
	})
	if err != nil {
		return nil, err
	}

	engine := gin.Default()
	if util.RunEnv() == util.EnvProd {
		gin.SetMode(gin.ReleaseMode)
	}
	engine.Use(middle.Cors())
	group := engine.Group("/api/v1")

	module.New(mysql, log)

	manager := module.NewManager()
	manager.BindAll(group)
	return engine, nil
}
