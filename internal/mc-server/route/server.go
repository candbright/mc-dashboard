package route

import (
	"net/http"
	"strconv"

	"github.com/candbright/go-server/internal/mc-server/core"
	"github.com/candbright/go-server/internal/mc-server/model"
	"github.com/candbright/go-server/pkg/rest"
	"github.com/gin-gonic/gin"
)

func init() {
	registerRoute(func(e *gin.Engine) {
		e.POST("/server/info/list", rest.H(listCurrentServerInfo))
		e.POST("/server/:id/info/get", rest.H(getCurrentServerInfo))
		e.POST("/server/:id/download_start", rest.H(startDownloadServer))
		e.POST("/server/:id/start", rest.H(startServer))
		e.POST("/server/:id/stop", rest.H(stopServer))
	})
}

func listCurrentServerInfo(c *gin.Context) error {
	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}

	sizeNum, err := strconv.Atoi(size)
	if err != nil || sizeNum < 1 {
		sizeNum = 10
	}

	servers := manager.GetServers()
	infos := make([]model.ServerInfo, 0)
	for _, server := range servers {
		info, err := transServerInfo(server)
		if err != nil {
			continue
		}
		infos = append(infos, info)
	}

	// Calculate pagination
	total := len(infos)
	start := (pageNum - 1) * sizeNum
	end := start + sizeNum
	if start >= total {
		start = total
	}
	if end > total {
		end = total
	}

	// Return paginated result
	return rest.Json(gin.H{
		"total": total,
		"items": infos[start:end],
	})
}

func getCurrentServerInfo(c *gin.Context) error {
	id := c.Param("id")
	server, err := manager.GetServer(id)
	if err != nil {
		return rest.ErrorWithStatus(http.StatusNotFound, err)
	}
	info, err := transServerInfo(server)
	if err != nil {
		return err
	}
	return rest.Json(info)
}

func startDownloadServer(c *gin.Context) error {
	id := c.Param("id")
	return manager.DownloadServer(id, "")
}

func startServer(c *gin.Context) error {
	id := c.Param("id")
	server, err := manager.GetServer(id)
	if err != nil {
		return rest.ErrorWithStatus(http.StatusNotFound, err)
	}
	err = server.Start()
	if err != nil {
		return err
	}
	return nil
}

func stopServer(c *gin.Context) error {
	id := c.Param("id")
	server, err := manager.GetServer(id)
	if err != nil {
		return rest.ErrorWithStatus(http.StatusNotFound, err)
	}
	err = server.Stop()
	if err != nil {
		return err
	}
	return nil
}

func transServerInfo(server *core.Server) (model.ServerInfo, error) {
	info := model.ServerInfo{
		ID:      server.GetID(),
		Name:    server.GetServerName(),
		Version: server.GetVersion(),
	}

	exist := server.ServerExist()
	info.Exist = exist

	downloading := server.Downloading()
	info.Downloading = downloading

	active := server.Active()
	info.Active = active

	serverProperties, _ := server.ServerProperties()
	if serverProperties != nil {
		info.ServerProperties = serverProperties.GetAll()
	}

	allowList, _ := server.GetAllowList()
	if allowList != nil {
		info.AllowList = allowList
	}
	return info, nil
}
