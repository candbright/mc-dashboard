package usecase

import (
	"fmt"
	"mime/multipart"
	"strconv"

	"github.com/candbright/mc-dashboard/internal/app/infra/minecraft"
	"github.com/candbright/mc-dashboard/internal/domain"
	"github.com/candbright/mc-dashboard/internal/pkg/errors"
	"github.com/sirupsen/logrus"
)

type SaveService struct {
	log       *logrus.Logger
	mcManager *minecraft.Manager
	repo      domain.SaveRepository
}

// NewSaveService 创建存档服务实例
func NewSaveService(log *logrus.Logger, mcManager *minecraft.Manager, repo domain.SaveRepository) domain.SaveService {
	return &SaveService{
		log:       log,
		mcManager: mcManager,
		repo:      repo,
	}
}

// UploadSave 上传存档
func (s *SaveService) UploadSave(file multipart.File, header *multipart.FileHeader) error {
	// 检查存档名称是否已存在
	existing, err := s.repo.FindByName(header.Filename)
	if err == nil && existing != nil {
		return errors.ErrSaveAlreadyExists
	}

	fileInfo, err := s.mcManager.AddSave(file, header)
	if err != nil {
		return errors.ErrSaveFile
	}

	save := &domain.Save{
		Name:         fileInfo.Name(),
		Size:         fileInfo.Size(),
		LastModified: fileInfo.ModTime(),
	}

	// 保存到数据库
	if err := s.repo.Create(save); err != nil {
		return errors.ErrSaveFile
	}
	return nil
}

// ListSaves 获取存档列表
func (s *SaveService) ListSaves(page, size int, order, orderBy string) ([]domain.Save, int64, error) {
	saves, total, err := s.repo.List(page, size, order, orderBy)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list saves from repo: %w", err)
	}
	return saves, total, nil
}

// DeleteSave 删除存档
func (s *SaveService) DeleteSave(id uint) error {
	// 查找存档
	save, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("save not found: %w", err)
	}

	// 删除文件
	if err := s.mcManager.DeleteSave(save.Name); err != nil {
		return fmt.Errorf("failed to delete save file: %w", err)
	}

	// 从数据库中删除
	if err := s.repo.Delete(save.ID); err != nil {
		return fmt.Errorf("failed to delete save file from repo: %w", err)
	}
	return nil
}

// ApplySave 应用存档
func (s *SaveService) ApplySave(saveID, serverID uint) error {
	// 查找存档
	save, err := s.repo.FindByID(saveID)
	if err != nil {
		return fmt.Errorf("save not found: %w", err)
	}

	if err := s.mcManager.ApplySave(save.Name, strconv.FormatUint(uint64(serverID), 10)); err != nil {
		return fmt.Errorf("failed to apply save: %w", err)
	}
	return nil
}
