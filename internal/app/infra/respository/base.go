package repository

import (
	"fmt"

	"github.com/candbright/mc-dashboard/internal/domain"
	"gorm.io/gorm"
)

// BaseRepositoryImpl 基础仓储实现
type BaseRepositoryImpl[T any] struct {
	db *gorm.DB
}

// NewBaseRepository 创建基础仓储实例
func NewBaseRepository[T any](db *gorm.DB) domain.BaseRepository[T] {
	return &BaseRepositoryImpl[T]{db: db}
}

func (r *BaseRepositoryImpl[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

func (r *BaseRepositoryImpl[T]) FindBy(key, value string) (*T, error) {
	var entity T
	err := r.db.Where(fmt.Sprintf("%s = ?", key), value).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepositoryImpl[T]) FindByID(id uint) (*T, error) {
	var entity T
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepositoryImpl[T]) Update(entity *T) error {
	return r.db.Model(entity).Omit("created_at", "updated_at").Updates(entity).Error
}

func (r *BaseRepositoryImpl[T]) Delete(id uint) error {
	var entity T
	return r.db.Delete(&entity, "id = ?", id).Error
}

func (r *BaseRepositoryImpl[T]) List(page, size int, order, orderBy string) ([]T, int64, error) {
	var entities []T
	var total int64

	// 构建查询
	query := r.db.Model(new(T))

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 设置排序
	if orderBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", orderBy, order))
	}

	// 设置分页
	if page > 0 && size > 0 {
		offset := (page - 1) * size
		query = query.Offset(offset).Limit(size)
	}

	// 执行查询
	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}
