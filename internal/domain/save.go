package domain

import (
	"mime/multipart"
	"time"
)

// Save 存档信息
type Save struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Save) TableName() string {
	return "mc_dashboard_saves"
}

// SaveRepository 存档仓储接口
type SaveRepository interface {
	BaseRepository[Save]
	FindByName(name string) (*Save, error)
}

// SaveService 存档服务接口
type SaveService interface {
	UploadSave(file multipart.File, header *multipart.FileHeader) error
	ListSaves(page, size int, order, orderBy string) ([]Save, int64, error)
	DeleteSave(id uint) error
	ApplySave(saveID, serverID uint) error
}
