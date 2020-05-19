package module

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type Service struct {
	DB  *gorm.DB
	Log *zap.SugaredLogger
}

func New(db *gorm.DB, log *zap.SugaredLogger) *Service {
	return &Service{
		DB:  db,
		Log: log,
	}
}

func GinHandler(s *Service, handler func(s *Service, c *gin.Context)) func(c *gin.Context) {
	return func(c *gin.Context) {
		handler(s, c)
	}
}
