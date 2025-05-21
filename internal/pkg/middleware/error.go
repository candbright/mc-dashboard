package middleware

import (
	"github.com/candbright/mc-dashboard/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ErrorHandler 统一处理错误响应
func ErrorHandler(log *logrus.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		// 如果已经有响应，则不再处理
		if ctx.Writer.Written() {
			return
		}

		// 获取最后一个错误
		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last().Err

			// 处理自定义错误
			appErr, ok := err.(*errors.AppError)
			if !ok {
				appErr = errors.ErrUnknownError.Wrap(err)
			}

			// 记录错误日志
			log.WithFields(logrus.Fields{
				"code": appErr.Code,
				"data": appErr.Data,
			}).Error(appErr.Error())

			// 返回错误响应
			ctx.JSON(appErr.HTTPStatus, appErr)
			return
		}
	}
}
