package handler

import (
	"net/http"
	"strconv"

	"github.com/candbright/mc-dashboard/internal/domain"
	"github.com/candbright/mc-dashboard/internal/pkg/errors"
	"github.com/gin-gonic/gin"
)

// SaveHandler 存档处理器
type SaveHandler struct {
	service domain.SaveService
}

// NewSaveHandler 创建存档处理器实例
func NewSaveHandler(service domain.SaveService) *SaveHandler {
	return &SaveHandler{service: service}
}

// RegisterRoutes 注册路由
func (h *SaveHandler) RegisterRoutes(r *gin.RouterGroup) {
	saves := r.Group("/saves")
	{
		saves.POST("", h.UploadSave)
		saves.GET("", h.ListSaves)
		saves.DELETE("/:id", h.DeleteSave)
		saves.POST("/apply", h.ApplySave)
	}
}

// UploadSave godoc
// @Summary 上传存档
// @Description 上传新的存档文件
// @Tags 存档管理
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param file formData file true "存档文件"
// @Success 200 {object} domain.Save
// @Failure 400 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /saves [post]
func (h *SaveHandler) UploadSave(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.Error(err)
		return
	}

	defer file.Close()

	if err := h.service.UploadSave(file, header); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "存档上传成功"})
}

// ListSaves godoc
// @Summary 获取存档列表
// @Description 获取所有存档文件列表
// @Tags 存档管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码，从1开始" default(0)
// @Param size query int false "每页数量" default(0)
// @Param order query string false "排序方向，asc或desc" default(desc)
// @Param orderBy query string false "排序字段，支持id/name/version" default(create_at)
// @Success 200 {array} domain.Save
// @Failure 500 {object} errors.AppError
// @Router /saves [get]
func (h *SaveHandler) ListSaves(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "0")
	size := ctx.DefaultQuery("size", "0")
	order := ctx.DefaultQuery("order", "desc")            // 默认降序
	orderBy := ctx.DefaultQuery("order_by", "created_at") // 默认按创建时间排序

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 0 {
		pageNum = 0
	}

	sizeNum, err := strconv.Atoi(size)
	if err != nil || sizeNum < 0 {
		sizeNum = 0
	}

	saves, total, err := h.service.ListSaves(pageNum, sizeNum, order, orderBy)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": total,
		"items": saves,
	})
}

// DeleteSave godoc
// @Summary 删除存档
// @Description 删除指定的存档文件
// @Tags 存档管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "存档id"
// @Success 200 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /saves/{id} [delete]
func (h *SaveHandler) DeleteSave(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(errors.ErrInvalidRequestParam)
		return
	}
	saveID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.Error(errors.ErrParseRequestParam.Wrap(err))
		return
	}
	if err := h.service.DeleteSave(uint(saveID)); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ApplySave godoc
// @Summary 应用存档
// @Description 将指定的存档应用到服务器
// @Tags 存档管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body ApplySaveRequest true "存档应用参数"
// @Success 200 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /saves/apply [post]
func (h *SaveHandler) ApplySave(ctx *gin.Context) {
	var request ApplySaveRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(errors.ErrInvalidRequestBody)
		return
	}
	if err := h.service.ApplySave(request.SaveID, request.ServerID); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "应用成功"})
}

type ApplySaveRequest struct {
	SaveID   uint `json:"save_id" binding:"required"`
	ServerID uint `json:"server_id" binding:"required"`
}
