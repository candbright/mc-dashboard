package repository

import (
	"github.com/candbright/mc-dashboard/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	BaseRepositoryImpl[domain.User]
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &UserRepository{
		BaseRepositoryImpl: BaseRepositoryImpl[domain.User]{db: db},
	}
}

func (r *UserRepository) FindByPhone(phone string) (*domain.User, error) {
	return r.FindBy("phone", phone)
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	return r.FindBy("email", email)
}
