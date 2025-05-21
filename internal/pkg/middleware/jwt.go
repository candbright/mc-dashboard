package middleware

import (
	"strings"
	"time"

	"github.com/candbright/mc-dashboard/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(errors.ErrMissingAuth.HTTPStatus, errors.ErrMissingAuth)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(errors.ErrInvalidAuth.HTTPStatus, errors.ErrInvalidAuth)
			c.Abort()
			return
		}

		claims, err := parseToken(parts[1])
		if err != nil {
			c.JSON(errors.ErrInvalidToken.HTTPStatus, errors.ErrInvalidToken)
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(errors.ErrInvalidUserID.HTTPStatus, errors.ErrInvalidUserID)
			c.Abort()
			return
		}

		phone, ok := claims["phone"].(string)
		if !ok {
			c.JSON(errors.ErrInvalidPhone.HTTPStatus, errors.ErrInvalidPhone)
			c.Abort()
			return
		}

		status, ok := claims["status"].(float64)
		if !ok {
			c.JSON(errors.ErrInvalidStatus.HTTPStatus, errors.ErrInvalidStatus)
			c.Abort()
			return
		}

		if status == 0 {
			c.JSON(errors.ErrUserForbidden.HTTPStatus, errors.ErrUserForbidden)
			c.Abort()
			return
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			c.JSON(errors.ErrInvalidExp.HTTPStatus, errors.ErrInvalidExp)
			c.Abort()
			return
		}

		if time.Now().Unix() > int64(exp) {
			c.JSON(errors.ErrTokenExpired.HTTPStatus, errors.ErrTokenExpired)
			c.Abort()
			return
		}

		c.Set("user_id", uint(userID))
		c.Set("phone", phone)
		c.Next()
	}
}

func GenerateToken(userID uint, phone string, status int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"phone":   phone,
		"status":  status,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// TODO: 从配置中获取密钥
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", errors.ErrGenerateToken.Wrap(err)
	}

	return tokenString, nil
}

// parseToken 解析JWT令牌
func parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// TODO: 从配置中获取密钥
		return []byte("your-secret-key"), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
