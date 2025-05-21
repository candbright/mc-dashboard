package middleware

import (
	"context"
	"time"

	"github.com/candbright/mc-dashboard/internal/pkg/errors"
	"github.com/gin-gonic/gin"
)

// Timeout 超时中间件
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建一个带超时的上下文
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// 将新的上下文设置到请求中
		c.Request = c.Request.WithContext(ctx)

		// 创建一个通道来接收处理结果
		done := make(chan struct{}, 1)
		go func() {
			c.Next()
			done <- struct{}{}
		}()

		// 等待处理完成或超时
		select {
		case <-done:
			// 处理完成
			return
		case <-ctx.Done():
			// 超时
			c.Abort()
			c.JSON(errors.ErrTimeout.HTTPStatus, errors.ErrTimeout)
			return
		}
	}
}
