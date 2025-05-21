package handler

import (
	"net/http"
	"strconv"

	"github.com/candbright/mc-dashboard/internal/domain"
	"github.com/candbright/mc-dashboard/internal/pkg/errors"
	"github.com/gin-gonic/gin"
)

// ServerHandler 服务器处理器
type ServerHandler struct {
	service domain.ServerService
}

// NewServerHandler 创建服务器处理器实例
func NewServerHandler(service domain.ServerService) *ServerHandler {
	return &ServerHandler{service: service}
}

// RegisterRoutes 注册路由
func (h *ServerHandler) RegisterRoutes(r *gin.RouterGroup) {
	servers := r.Group("/servers")
	{
		servers.GET("", h.ListServerInfos)
		servers.GET("/:id", h.GetServerInfo)
		servers.POST("", h.CreateServer)
		servers.PUT("/:id", h.UpdateServer)
		servers.POST("/:id/start", h.StartServer)
		servers.POST("/:id/stop", h.StopServer)
		servers.DELETE("/:id", h.DeleteServer)
		servers.GET("/:id/server_properties", h.ListServerProperties)
		servers.PUT("/:id/server_properties", h.UpdateServerProperties)
		servers.GET("/:id/allowlist", h.ListAllowlist)
		servers.POST("/:id/allowlist/:username", h.AddAllowlist)
		servers.DELETE("/:id/allowlist/:username", h.DeleteAllowlist)
		servers.GET("/:id/console_log", h.GetConsoleLog)
	}
}

// ListServerInfos godoc
// @Summary 获取服务器列表
// @Description 获取所有服务器的列表，支持分页和排序
// @Tags 服务器
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码，从1开始" default(0)
// @Param size query int false "每页数量" default(0)
// @Param order query string false "排序方向，asc或desc" default(desc)
// @Param orderBy query string false "排序字段，支持id/name/version" default(created_at)
// @Success 200 {array} domain.ServerInfo "服务器列表"
// @Failure 500 {object} errors.AppError "服务器内部错误"
// @Router /servers [get]
func (h *ServerHandler) ListServerInfos(ctx *gin.Context) {
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

	servers, total, err := h.service.ListServerInfos(pageNum, sizeNum, order, orderBy)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": total,
		"items": servers,
	})
}

// GetServerInfo godoc
// @Summary 获取服务器信息
// @Description 根据服务器ID获取服务器详细信息
// @Tags 服务器
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "服务器ID"
// @Success 200 {object} domain.ServerInfo "服务器信息"
// @Failure 404 {object} errors.AppError "服务器不存在"
// @Failure 500 {object} errors.AppError "服务器内部错误"
// @Router /servers/{id} [get]
func (h *ServerHandler) GetServerInfo(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(errors.ErrInvalidRequestParam)
		return
	}
	serverID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.Error(errors.ErrParseRequestParam.Wrap(err))
		return
	}
	server, err := h.service.GetServerInfo(uint(serverID))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, server)
}

// CreateServer godoc
// @Summary 创建服务器
// @Description 创建一个新的Minecraft服务器
// @Tags 服务器
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateServerRequest true "服务器创建参数"
// @Success 200 {string} string "服务器创建成功"
// @Failure 400 {object} errors.AppError "请求参数错误"
// @Failure 500 {object} errors.AppError "服务器内部错误"
// @Router /servers [post]
func (h *ServerHandler) CreateServer(ctx *gin.Context) {
	var request CreateServerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(errors.ErrInvalidRequestBody.Wrap(err))
		return
	}
	if err := h.service.CreateServer(&domain.Server{
		Name:        request.Name,
		Description: request.Description,
		WorldName:   request.WorldName,
		Version:     request.Version,
	}); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "服务器创建成功"})
}

// UpdateServer godoc
// @Summary 编辑服务器
// @Description 修改服务器的配置信息
// @Tags 服务器
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "服务器ID"
// @Param request body UpdateServerRequest true "服务器更新参数"
// @Success 200 {string} string "服务器编辑成功"
// @Failure 404 {object} errors.AppError "服务器不存在"
// @Failure 500 {object} errors.AppError "服务器内部错误"
// @Router /servers/{id} [put]
func (h *ServerHandler) UpdateServer(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(errors.ErrInvalidRequestParam)
		return
	}
	serverID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.Error(errors.ErrParseRequestParam.Wrap(err))
		return
	}

	var request UpdateServerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(errors.ErrInvalidRequestBody)
		return
	}
	if err := h.service.UpdateServer(&domain.Server{
		ID:          uint(serverID),
		Name:        request.Name,
		Description: request.Description,
	}); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "服务器编辑成功"})
}

// StartServer godoc
// @Summary 启动服务器
// @Description 启动指定的Minecraft服务器
// @Tags 服务器
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "服务器ID"
// @Success 200 {string} string "服务器启动成功"
// @Failure 404 {object} errors.AppError "服务器不存在"
// @Failure 400 {object} errors.AppError "服务器已经在运行"
// @Failure 500 {object} errors.AppError "服务器内部错误"
// @Router /servers/{id}/start [post]
func (h *ServerHandler) StartServer(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(errors.ErrInvalidRequestParam)
		return
	}
	serverID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.Error(errors.ErrParseRequestParam.Wrap(err))
		return
	}
	if err := h.service.StartServer(uint(serverID)); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "服务器启动成功"})
}

// StopServer godoc
// @Summary 停止服务器
// @Description 停止指定的Minecraft服务器
// @Tags 服务器
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "服务器ID"
// @Success 200 {string} string "服务器停止成功"
// @Failure 404 {object} errors.AppError "服务器不存在"
// @Failure 400 {object} errors.AppError "服务器未在运行"
// @Failure 500 {object} errors.AppError "服务器内部错误"
// @Router /servers/{id}/stop [post]
func (h *ServerHandler) StopServer(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(errors.ErrInvalidRequestParam)
		return
	}
	serverID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.Error(errors.ErrParseRequestParam.Wrap(err))
		return
	}
	if err := h.service.StopServer(uint(serverID)); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "服务器停止成功"})
}

// DeleteServer godoc
// @Summary 删除服务器
// @Description 删除指定的Minecraft服务器
// @Tags 服务器
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "服务器ID"
// @Success 200 {string} string "服务器删除成功"
// @Failure 404 {object} errors.AppError "服务器不存在"
// @Failure 500 {object} errors.AppError "服务器内部错误"
// @Router /servers/{id} [delete]
func (h *ServerHandler) DeleteServer(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(errors.ErrInvalidRequestParam)
		return
	}
	serverID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.Error(errors.ErrParseRequestParam.Wrap(err))
		return
	}
	if err := h.service.DeleteServer(uint(serverID)); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "服务器删除成功"})
}

// ListServerProperties godoc
// @Summary 获取服务器属性
// @Description 获取指定服务器的所有属性配置
// @Tags 服务器
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "服务器ID"
// @Success 200 {object} map[string]string "服务器属性"
// @Failure 404 {object} errors.AppError "服务器不存在"
// @Failure 500 {object} errors.AppError "服务器内部错误"
// @Router /servers/{id}/server_properties [get]
func (h *ServerHandler) ListServerProperties(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(errors.ErrInvalidRequestParam)
		return
	}
	serverID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.Error(errors.ErrParseRequestParam.Wrap(err))
		return
	}

	properties, err := h.service.GetServerProperties(uint(serverID))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, properties)
}

// UpdateServerProperties godoc
// @Summary 更新服务器属性
// @Description 更新指定服务器的属性配置
// @Tags 服务器
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "服务器ID"
// @Param request body map[string]string true "服务器属性"
// @Success 200 {string} string "服务器属性更新成功"
// @Failure 400 {object} errors.AppError "请求参数错误"
// @Failure 404 {object} errors.AppError "服务器不存在"
// @Failure 500 {object} errors.AppError "服务器内部错误"
// @Router /servers/{id}/server_properties [put]
func (h *ServerHandler) UpdateServerProperties(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(errors.ErrInvalidRequestParam)
		return
	}
	serverID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.Error(errors.ErrParseRequestParam.Wrap(err))
		return
	}

	var properties map[string]string
	if err := ctx.ShouldBindJSON(&properties); err != nil {
		ctx.Error(errors.ErrInvalidRequestBody.Wrap(err))
		return
	}

	if err := h.service.UpdateServerProperties(uint(serverID), properties); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "服务器属性更新成功"})
}

// ListAllowlist godoc
// @Summary 获取白名单列表
// @Description 获取指定服务器的白名单用户列表
// @Tags 服务器
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "服务器ID"
// @Success 200 {object} domain.AllowList "白名单用户列表"
// @Failure 404 {object} errors.AppError "服务器不存在"
// @Failure 500 {object} errors.AppError "服务器内部错误"
// @Router /servers/{id}/allowlist [get]
func (h *ServerHandler) ListAllowlist(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(errors.ErrInvalidRequestParam)
		return
	}
	serverID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.Error(errors.ErrParseRequestParam.Wrap(err))
		return
	}

	allowList, err := h.service.GetAllowList(uint(serverID))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, allowList)
}

// AddAllowlist godoc
// @Summary 添加白名单用户
// @Description 将指定用户添加到服务器的白名单中
// @Tags 服务器
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "服务器ID"
// @Param username path string true "用户名"
// @Success 200 {string} string "添加白名单用户成功"
// @Failure 400 {object} errors.AppError "请求参数错误"
// @Failure 404 {object} errors.AppError "服务器不存在"
// @Failure 500 {object} errors.AppError "服务器内部错误"
// @Router /servers/{id}/allowlist/{username} [post]
func (h *ServerHandler) AddAllowlist(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(errors.ErrInvalidRequestParam)
		return
	}
	username := ctx.Param("username")
	if username == "" {
		ctx.Error(errors.ErrInvalidRequestParam)
		return
	}
	serverID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.Error(errors.ErrParseRequestParam.Wrap(err))
		return
	}

	if err := h.service.AddAllowListUser(uint(serverID), username); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "添加白名单用户成功"})
}

// DeleteAllowlist godoc
// @Summary 删除白名单用户
// @Description 从服务器的白名单中删除指定用户
// @Tags 服务器
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "服务器ID"
// @Param username path string true "用户名"
// @Success 200 {string} string "删除白名单用户成功"
// @Failure 404 {object} errors.AppError "服务器不存在或用户不在白名单中"
// @Failure 500 {object} errors.AppError "服务器内部错误"
// @Router /servers/{id}/allowlist/{username} [delete]
func (h *ServerHandler) DeleteAllowlist(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(errors.ErrInvalidRequestParam)
		return
	}
	serverID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.Error(errors.ErrParseRequestParam.Wrap(err))
		return
	}

	username := ctx.Param("username")
	if username == "" {
		ctx.Error(errors.ErrInvalidRequestParam)
		return
	}

	if err := h.service.DeleteAllowListUser(uint(serverID), username); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "删除白名单用户成功"})
}

// GetConsoleLog godoc
// @Summary 获取控制台日志
// @Description 获取指定服务器的控制台日志
// @Tags 服务器
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "服务器ID"
// @Param line query int false "获取最后几行日志" default(100)
// @Success 200 {object} ConsoleLogResponse "日志内容"
// @Failure 404 {object} errors.AppError "服务器不存在"
// @Failure 500 {object} errors.AppError "服务器内部错误"
// @Router /servers/{id}/console_log [get]
func (h *ServerHandler) GetConsoleLog(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(errors.ErrInvalidRequestParam)
		return
	}
	serverID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.Error(errors.ErrParseRequestParam.Wrap(err))
		return
	}

	// 获取行数参数，默认为100行
	lineStr := ctx.DefaultQuery("line", "100")
	line, err := strconv.Atoi(lineStr)
	if err != nil || line <= 0 {
		line = 100
	}

	// 获取日志内容
	logContent, err := h.service.GetConsoleLog(uint(serverID), line)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, ConsoleLogResponse{
		Content: logContent,
	})
}

type CreateServerRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"omitempty"`
	WorldName   string `json:"world_name" binding:"omitempty"`
	Version     string `json:"version" binding:"omitempty"`
}

type UpdateServerRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"omitempty"`
}

type ConsoleLogResponse struct {
	Content string `json:"content"`
}
