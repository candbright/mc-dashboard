package handler

import (
	"net/http"

	"github.com/candbright/mc-dashboard/internal/domain"
	"github.com/candbright/mc-dashboard/internal/pkg/errors"
	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	service domain.UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// RegisterRoutes 注册路由
func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/user")
	{
		users.GET("", h.GetUserInfo)
		users.POST("/logout", h.Logout)
	}
}

// Register godoc
// @Summary 用户注册
// @Description 用户注册接口
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body RegisterRequest true "注册信息"
// @Success 200 {object} domain.User
// @Failure 400 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /register [post]
func (h *UserHandler) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errors.ErrInvalidRequestBody.Wrap(err))
		return
	}

	user, err := h.service.Register(req.Phone, req.Email, req.Password)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// Login godoc
// @Summary 用户登录
// @Description 用户登录接口，支持手机号密码、邮箱密码、手机验证码登录
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body LoginRequest true "登录信息"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errors.ErrInvalidRequestBody.Wrap(err))
		return
	}

	token, err := h.service.Login(req.LoginType, req.Account, req.Credential)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{Token: token})
}

// GetUserInfo godoc
// @Summary 获取用户信息
// @Description 获取当前登录用户的信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} domain.User
// @Failure 401 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /user [get]
func (h *UserHandler) GetUserInfo(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	user, err := h.service.GetUserInfo(userID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// Logout godoc
// @Summary 用户登出
// @Description 用户登出接口，清除用户登录状态
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {string} string "登出成功"
// @Failure 401 {object} errors.AppError "未授权"
// @Failure 500 {object} errors.AppError "服务器错误"
// @Router /user/logout [post]
func (h *UserHandler) Logout(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	if err := h.service.Logout(userID); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required,len=11" example:"13800138000"`
	Email    string `json:"email" binding:"omitempty,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"123456"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	LoginType  string `json:"login_type" binding:"required,oneof=phone_password email_password phone_code" example:"phone_password"`
	Account    string `json:"account" binding:"required" example:"13800138000"`
	Credential string `json:"credential" binding:"required" example:"123456"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
