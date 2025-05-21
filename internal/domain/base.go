package domain

// BaseRepository 基础仓储接口
type BaseRepository[T any] interface {
	Create(entity *T) error
	FindBy(key, value string) (*T, error)
	FindByID(id uint) (*T, error)
	Update(entity *T) error
	Delete(id uint) error
	List(page, size int, order, orderBy string) ([]T, int64, error)
}
