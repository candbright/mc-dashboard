package repository

import (
	"github.com/candbright/mc-dashboard/internal/domain"
	"gorm.io/gorm"
)

type ServerRepository struct {
	BaseRepositoryImpl[domain.Server]
}

// NewServerRepository 创建服务器仓储实例
func NewServerRepository(db *gorm.DB) domain.ServerRepository {
	return &ServerRepository{
		BaseRepositoryImpl: BaseRepositoryImpl[domain.Server]{db: db},
	}
}

func (r *ServerRepository) FindByName(name string) (*domain.Server, error) {
	return r.FindBy("name", name)
}
