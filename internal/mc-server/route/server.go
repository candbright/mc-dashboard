package route

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
		e.POST("/server/saves/upload", rest.H(uploadServerFile))
		e.POST("/server/saves/list", rest.H(listSaves))
		e.POST("/server/saves/delete", rest.H(deleteSave))
		e.POST("/server/saves/apply", rest.H(applySave))
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

func uploadServerFile(c *gin.Context) error {
	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		return rest.ErrorWithStatus(http.StatusBadRequest, fmt.Errorf("无法获取上传文件: %v", err))
	}
	defer file.Close()

	// 检查文件扩展名
	ext := filepath.Ext(header.Filename)
	if ext != ".mcworld" && ext != ".zip" {
		return rest.ErrorWithStatus(http.StatusBadRequest, fmt.Errorf("只支持上传 MCWORLD 或 ZIP 格式的存档文件"))
	}

	// 创建上传目录
	uploadDir := manager.UploadDir()
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return rest.ErrorWithStatus(http.StatusInternalServerError, fmt.Errorf("创建上传目录失败: %v", err))
	}

	// 如果是 mcworld 文件，将后缀改为 zip
	filename := header.Filename
	if ext == ".mcworld" {
		filename = filename[:len(filename)-len(ext)] + ".zip"
	}

	// 创建目标文件
	filepath := filepath.Join(uploadDir, filename)
	out, err := os.Create(filepath)
	if err != nil {
		return rest.ErrorWithStatus(http.StatusInternalServerError, fmt.Errorf("创建文件失败: %v", err))
	}
	defer out.Close()

	// 将上传的文件内容复制到目标文件
	if _, err := io.Copy(out, file); err != nil {
		return rest.ErrorWithStatus(http.StatusInternalServerError, fmt.Errorf("保存文件失败: %v", err))
	}

	// 触发一次存档目录扫描
	manager.ScanSaves()

	return rest.Json(gin.H{
		"message":  "文件上传成功",
		"filename": filename,
		"size":     header.Size,
	})
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

// listSaves 获取存档列表
func listSaves(c *gin.Context) error {
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

	// 获取所有存档信息
	saves := manager.GetSaves()

	// 转换为切片以便分页
	saveList := make([]core.SaveInfo, 0, len(saves))
	for _, save := range saves {
		saveList = append(saveList, save)
	}

	// 计算分页
	total := len(saveList)
	start := (pageNum - 1) * sizeNum
	end := start + sizeNum
	if start >= total {
		start = total
	}
	if end > total {
		end = total
	}

	// 返回分页结果
	return rest.Json(gin.H{
		"total": total,
		"items": saveList[start:end],
	})
}

// deleteSave 删除存档
func deleteSave(c *gin.Context) error {
	var req struct {
		Filename string `json:"filename" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		return rest.ErrorWithStatus(http.StatusBadRequest, fmt.Errorf("无效的请求参数: %v", err))
	}

	// 获取存档信息
	saves := manager.GetSaves()
	saveInfo, exists := saves[req.Filename]
	if !exists {
		return rest.ErrorWithStatus(http.StatusNotFound, fmt.Errorf("存档文件不存在"))
	}

	// 删除文件
	if err := os.Remove(saveInfo.Path); err != nil {
		return rest.ErrorWithStatus(http.StatusInternalServerError, fmt.Errorf("删除文件失败: %v", err))
	}

	// 触发一次存档扫描
	manager.ScanSaves()

	return rest.Json(gin.H{
		"message": "存档删除成功",
	})
}

// applySave 应用存档到指定服务器
func applySave(c *gin.Context) error {
	var req struct {
		ServerID string `json:"server_id" binding:"required"`
		Filename string `json:"filename" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		return rest.ErrorWithStatus(http.StatusBadRequest, fmt.Errorf("无效的请求参数: %v", err))
	}

	// 获取服务器信息
	server, err := manager.GetServer(req.ServerID)
	if err != nil {
		return rest.ErrorWithStatus(http.StatusNotFound, fmt.Errorf("服务器不存在"))
	}

	// 获取存档信息
	saves := manager.GetSaves()
	saveInfo, exists := saves[req.Filename]
	if !exists {
		return rest.ErrorWithStatus(http.StatusNotFound, fmt.Errorf("存档文件不存在"))
	}

	// 应用存档
	if err := server.ApplySave(saveInfo.Path); err != nil {
		return rest.ErrorWithStatus(http.StatusInternalServerError, fmt.Errorf("应用存档失败: %v", err))
	}

	return rest.Json(gin.H{
		"message": "存档应用成功",
	})
}
