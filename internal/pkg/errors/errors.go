package errors

import (
	"fmt"
	"net/http"
)

// AppError 应用错误
type AppError struct {
	Code       int    `json:"code"`              // 错误码
	Message    string `json:"message,omitempty"` // 错误信息
	Err        error  `json:"-"`                 // 原始错误
	HTTPStatus int    `json:"-"`                 // HTTP状态码
	Data       any    `json:"data,omitempty"`    // 返回数据
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// New 创建新的应用错误
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap 包装已有错误
func (e *AppError) Wrap(err error) *AppError {
	e.Err = err
	return e
}

// WithStatus 添加错误详情
func (e *AppError) WithStatus(statusCode int) *AppError {
	e.HTTPStatus = statusCode
	return e
}

// WithData 添加返回数据
func (e *AppError) WithData(data any) *AppError {
	e.Data = data
	return e
}

// SuccessMessage 返回消息
func SuccessMessage(message string) *AppError {
	return &AppError{
		Code:    0,
		Message: message,
	}
}

// Success 成功并返回数据
func Success(data any) *AppError {
	return &AppError{
		Code: 0,
		Data: data,
	}
}

// 预定义错误
var (
	ErrUnknownError        = New(1, "未处理的错误").WithStatus(http.StatusInternalServerError)
	ErrInvalidRequestBody  = New(2, "请求体验证失败").WithStatus(http.StatusBadRequest)
	ErrInvalidRequestParam = New(2, "请求参数验证失败").WithStatus(http.StatusBadRequest)
	ErrParseRequestParam   = New(2, "请求参数类型转换失败").WithStatus(http.StatusInternalServerError)
	ErrUploadFile          = New(3, "文件上传失败").WithStatus(http.StatusBadRequest)
	ErrInvalidToken        = New(5, "无效的token").WithStatus(http.StatusUnauthorized)
	ErrMissingAuth         = New(6, "未提供认证信息").WithStatus(http.StatusUnauthorized)
	ErrInvalidAuth         = New(7, "认证格式错误").WithStatus(http.StatusUnauthorized)
	ErrInvalidUserID       = New(8, "无效的用户ID").WithStatus(http.StatusUnauthorized)
	ErrInvalidPhone        = New(8, "无效的手机号").WithStatus(http.StatusUnauthorized)
	ErrInvalidStatus       = New(8, "无效的用户状态").WithStatus(http.StatusUnauthorized)
	ErrInvalidExp          = New(8, "无效的token过期时间").WithStatus(http.StatusUnauthorized)
	ErrUserForbidden       = New(8, "用户被禁止").WithStatus(http.StatusUnauthorized)
	ErrTokenExpired        = New(8, "token已过期").WithStatus(http.StatusUnauthorized)
	ErrTimeout             = New(9, "请求超时").WithStatus(http.StatusGatewayTimeout)

	ErrPhoneRegistered      = New(2, "手机号已注册").WithStatus(http.StatusBadRequest)
	ErrEmailRegistered      = New(3, "邮箱已注册").WithStatus(http.StatusBadRequest)
	ErrGeneratePassword     = New(3, "密码加密失败").WithStatus(http.StatusInternalServerError)
	ErrCreateUser           = New(4, "创建用户失败").WithStatus(http.StatusInternalServerError)
	ErrUpdateUser           = New(5, "更新用户失败").WithStatus(http.StatusInternalServerError)
	ErrUserNotFound         = New(4, "用户不存在").WithStatus(http.StatusNotFound)
	ErrWrongPassword        = New(4, "密码错误").WithStatus(http.StatusUnauthorized)
	ErrUnsupportedLoginType = New(12, "不支持的登录方式").WithStatus(http.StatusBadRequest)
	ErrGenerateToken        = New(13, "生成token失败").WithStatus(http.StatusInternalServerError)

	ErrInvalidEmail          = New(9, "无效的邮箱").WithStatus(http.StatusBadRequest)
	ErrInvalidPassword       = New(10, "无效的密码").WithStatus(http.StatusBadRequest)
	ErrInvalidCode           = New(11, "无效的验证码").WithStatus(http.StatusBadRequest)
	ErrUnsupportedFileFormat = New(13, "不支持的文件格式").WithStatus(http.StatusBadRequest)
	ErrSaveAlreadyExists     = New(14, "存档已存在").WithStatus(http.StatusConflict)
	ErrSaveFile              = New(15, "保存存档文件失败").WithStatus(http.StatusInternalServerError)
)
