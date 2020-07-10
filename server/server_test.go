package server

import (
	"testing"

	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
	"github.com/zhcppy/golib/db"
)

func Test_Config_Toml(t *testing.T) {
	marshal, err := toml.Marshal(Config{DbConfig: db.DefMysqlCfg()})
	assert.NoError(t, err)
	t.Log("\n", string(marshal))
}
