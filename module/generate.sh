#!/usr/bin/env bash

# shell is funny

set -e -x

function uppercase() {
   echo "$(echo ${1:0:1} | tr  '[a-z]' '[A-Z]')${1:1}"
}

if [[ -f go.mod ]]; then
    cd module
fi

module_name=$1
if [[ $module_name == ""  ]]; then
    read -p "请输入模块名称: " module_name
fi

if [[ -d $module_name ]]; then
    echo "$module_name 模块已存在！";exit 1
fi

mkdir ${module_name}

cat << EOF > ${module_name}/models.go
package $module_name

import (
	"github.com/zhcppy/golib/db"
)

type $(uppercase $module_name) struct {
	db.BaseModel
}
EOF

cat << EOF > ${module_name}/service.go
package ${module_name}

import (
	"github.com/zhcppy/golib/module"
)

type Service struct {
	*module.Service
}

func New(ser *module.Service) *Service {
	return &Service{Service: ser}
}

EOF


cat << EOF > ${module_name}/handler.go
package ${module_name}

import (
	"github.com/zhcppy/golib/common/ginn"
	"github.com/zhcppy/golib/module"
	"github.com/gin-gonic/gin"
)

func (s *Service) Bind(router *gin.RouterGroup) {
	router.POST("/register", module.GinHandler(s.Service, RegisterHandler))
}

func RegisterHandler(s *module.Service, c *gin.Context) {
	ginn.RespSuccess(c, nil)
}

EOF

git add -v ${module_name}

