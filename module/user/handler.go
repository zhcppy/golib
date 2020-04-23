package user

import (
	"github.com/gin-gonic/gin"
	"github.com/zhcppy/golib/common/ginn"
	"github.com/zhcppy/golib/module"
)

func (s *Service) Bind(router *gin.RouterGroup) {
	router.POST("/register", module.GinHandler(s.Service, RegisterHandler))
}

func RegisterHandler(s *module.Service, c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		ginn.RespFail(c, ginn.CodeFailed, err.Error())
		return
	}
	if err := s.DB.Create(&user).Error; err != nil {
		ginn.RespFail(c, ginn.CodeFailed, err.Error())
		return
	}
	ginn.RespSuccess(c, nil)
}
