package server

import (
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
	"github.com/zhcppy/golib/db"
	"testing"
)

func Test_Config_Toml(t *testing.T) {
	marshal, err := toml.Marshal(Config{DbConfig: db.DefaultCfg()})
	assert.NoError(t, err)
	t.Log("\n", string(marshal))
}
