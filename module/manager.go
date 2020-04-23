package module

import "github.com/gin-gonic/gin"

type Router interface {
	Bind(group *gin.RouterGroup)
}

type Manager struct {
	routers []Router
}

func NewManager(ser ...Router) *Manager {
	return &Manager{routers: ser}
}

func (m *Manager) BindAll(group *gin.RouterGroup) {
	for _, router := range m.routers {
		router.Bind(group)
	}
}
