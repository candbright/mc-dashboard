package route

import (
	"net/http"

	"github.com/candbright/go-server/internal/mc-server/utils"
	"github.com/candbright/go-server/pkg/rest"
	"github.com/gin-gonic/gin"
)

func init() {
	registerRoute(func(e *gin.Engine) {
		e.POST("/server/:id/server_properties/get", rest.H(getServerProperties))
		e.POST("/server/:id/server_properties/set", rest.H(setServerProperties))
	})
}

type setServerPropertiesReq map[string]string

func getServerProperties(c *gin.Context) error {
	id := c.Param("id")
	server, err := manager.GetServer(id)
	if err != nil {
		return rest.ErrorWithStatus(http.StatusNotFound, err)
	}
	serverProperties, err := server.ServerProperties()
	if err != nil {
		return rest.ErrorWithStatus(http.StatusNotFound, err)
	}
	return rest.Json(serverProperties.GetAll())
}

func setServerProperties(c *gin.Context) error {
	req := new(setServerPropertiesReq)
	if err := c.ShouldBindJSON(req); err != nil {
		return err
	}
	id := c.Param("id")
	server, err := manager.GetServer(id)
	if err != nil {
		return rest.ErrorWithStatus(http.StatusNotFound, err)
	}
	serverProperties, err := server.ServerProperties()
	if err != nil {
		return err
	}

	// 获取当前配置
	currentConfig := serverProperties.GetAll()

	// 将请求中的键转换为小驼峰格式并更新配置
	for key, value := range *req {
		if value != "" { // 只更新非空值
			camelKey := utils.ToCamelCase(key)
			currentConfig[camelKey] = value
		}
	}

	err = serverProperties.SetAll(currentConfig)
	if err != nil {
		return err
	}
	err = server.Reload()
	if err != nil {
		return err
	}
	return nil
}
