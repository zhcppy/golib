package user

import (
	"github.com/jinzhu/gorm"
	"github.com/zhcppy/golib/module"
)

type Service struct {
	*module.Service
}

func New(ser *module.Service) *Service {
	return &Service{ser}
}

func (s *Service) FindByName(name string) (user *User, err error) {
	user = &User{}
	if err := s.DB.First(user, "name = ?", name).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, nil
		}
		s.Log.Error("user find by name", err.Error())
		return nil, err
	}
	return
}
