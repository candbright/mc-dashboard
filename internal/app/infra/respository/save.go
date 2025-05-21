package repository

import (
	"github.com/candbright/mc-dashboard/internal/domain"
	"gorm.io/gorm"
)

type SaveRepository struct {
	BaseRepositoryImpl[domain.Save]
}

// NewSaveRepository 创建存档仓储实例
func NewSaveRepository(db *gorm.DB) domain.SaveRepository {
	return &SaveRepository{
		BaseRepositoryImpl: BaseRepositoryImpl[domain.Save]{db: db},
	}
}

func (r *SaveRepository) FindByName(name string) (*domain.Save, error) {
	return r.FindBy("name", name)
}
