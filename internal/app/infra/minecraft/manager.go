package minecraft

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"
	"time"

	"github.com/candbright/mc-dashboard/internal/domain"
	"github.com/candbright/mc-dashboard/internal/pkg/utils"
	"github.com/candbright/mc-dashboard/pkg/async"
	"github.com/sirupsen/logrus"
)

var (
	DefaultWorldName = "Dedicated Server"
)

type ManagerConfig struct {
	RootDir      string
	LoadInterval time.Duration
	CacheTTL     time.Duration
}

type Manager struct {
	rootDir               string
	serversDir            string
	versionsDir           string
	savesDir              string
	backupsDir            string
	servers               *sync.Map // key: server-id, value: *Server，存储所有服务器实例
	saves                 *sync.Map // key: filename, value: SaveInfo，存储所有存档信息
	createServerExecutors *sync.Map // key: server-id, value: *async.AsyncTaskExecutor，存储每个服务器的创建任务执行器
	loadInterval          time.Duration
	cacheTTL              time.Duration
	lastLoad              time.Time
	mu                    sync.RWMutex
}

// CreateServerConfig 创建服务器的配置
type CreateServerConfig struct {
	ID        string // 可选，不指定时自动生成
	WorldName string // 可选，不指定时使用默认名称
	Version   string // 可选，不指定时使用最新版本
}

func NewManager(cfg ManagerConfig) *Manager {
	if cfg.LoadInterval == 0 {
		cfg.LoadInterval = time.Minute
	}
	if cfg.CacheTTL == 0 {
		cfg.CacheTTL = 5 * time.Minute
	}

	mgr := &Manager{
		rootDir:               cfg.RootDir,
		serversDir:            path.Join(cfg.RootDir, "servers"),
		versionsDir:           path.Join(cfg.RootDir, "versions"),
		savesDir:              path.Join(cfg.RootDir, "saves"),
		backupsDir:            path.Join(cfg.RootDir, "backups"),
		servers:               &sync.Map{},
		saves:                 &sync.Map{},
		createServerExecutors: &sync.Map{},
		loadInterval:          cfg.LoadInterval,
		cacheTTL:              cfg.CacheTTL,
		lastLoad:              time.Now().Add(-cfg.CacheTTL), // 设置为一个已经过期的时间，确保第一次加载会执行
	}

	// 初始化时加载一次服务器列表
	if err := mgr.ScanServers(); err != nil {
		logrus.WithError(err).Error("Failed to load servers during initialization")
	}

	// 初始化时扫描一次存档目录
	if err := mgr.ScanSaves(); err != nil {
		logrus.WithError(err).Error("Failed to scan saves during initialization")
	}

	// 启动定期加载
	mgr.StartScanServers()

	// 启动存档扫描
	mgr.StartScanSaves()

	return mgr
}

func (manager *Manager) StartScanServers() {
	go func() {
		ticker := time.NewTicker(manager.loadInterval)
		defer ticker.Stop()

		for range ticker.C {
			if err := manager.ScanServers(); err != nil {
				logrus.WithError(err).Error("Failed to load servers")
			}
		}
	}()
}

func (manager *Manager) LatestVersion() (string, error) {
	// 根据操作系统选择不同的下载页面
	var downloadURL string
	switch runtime.GOOS {
	case "linux":
		downloadURL = "https://www.minecraft.net/en-us/download/server/bedrock"
	case "windows":
		downloadURL = "https://www.minecraft.net/en-us/download/server/bedrock"
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// 最大重试次数
	maxRetries := 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		// 创建请求
		req, err := utils.CreateWindowsHTTPRequest("GET", downloadURL, nil)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d/%d failed to create request: %w", i+1, maxRetries, err)
			logrus.WithError(err).Warnf("Failed to create request, retrying... (%d/%d)", i+1, maxRetries)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		// 发送请求
		resp, err := utils.GetHTTPClient().Do(req)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d/%d failed: %w", i+1, maxRetries, err)
			logrus.WithError(err).Warnf("Failed to get download page, retrying... (%d/%d)", i+1, maxRetries)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		// 检查响应状态
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			lastErr = fmt.Errorf("attempt %d/%d failed: status code %d", i+1, maxRetries, resp.StatusCode)
			logrus.WithField("status_code", resp.StatusCode).Warnf("Failed to get download page, retrying... (%d/%d)", i+1, maxRetries)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		// 读取响应内容
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("attempt %d/%d failed: %w", i+1, maxRetries, err)
			logrus.WithError(err).Warnf("Failed to read response body, retrying... (%d/%d)", i+1, maxRetries)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		// 使用正则表达式匹配版本号
		re := regexp.MustCompile(`bedrock-server-(\d+\.\d+\.\d+\.\d+)\.zip`)
		matches := re.FindStringSubmatch(string(body))
		if len(matches) < 2 {
			lastErr = fmt.Errorf("attempt %d/%d failed: version number not found", i+1, maxRetries)
			logrus.Warnf("Failed to find version number, retrying... (%d/%d)", i+1, maxRetries)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		version := matches[1]
		logrus.WithField("version", version).Info("Found latest version")
		return version, nil
	}

	return "", fmt.Errorf("failed to get latest version after %d attempts: %w", maxRetries, lastErr)
}

func (manager *Manager) ZipFileName(version string) string {
	return fmt.Sprintf("bedrock-server-%s.zip", version)
}

func (manager *Manager) ZipFile(version string) string {
	return path.Join(manager.versionsDir, manager.ZipFileName(version))
}

func (manager *Manager) ZipExist(version string) bool {
	return utils.FileExist(manager.ZipFile(version))
}

func (manager *Manager) DownloadVersionFile(version string) error {
	//检测是否存在当前版本的zip文件
	exist := manager.ZipExist(version)
	if exist {
		return nil
	}

	//不存在当前版本的zip文件，开始下载
	var downloadUrl string
	switch runtime.GOOS {
	case "linux":
		downloadUrl = fmt.Sprintf("https://www.minecraft.net/bedrockdedicatedserver/bin-linux/bedrock-server-%s.zip", version)
	case "windows":
		downloadUrl = fmt.Sprintf("https://www.minecraft.net/bedrockdedicatedserver/bin-win/bedrock-server-%s.zip", version)
	}

	err := os.MkdirAll(manager.versionsDir, os.ModePerm)
	if err != nil {
		return err
	}

	return utils.Download(downloadUrl, manager.ZipFile(version))
}

// ScanServers 加载服务器列表
func (manager *Manager) ScanServers() error {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	// 检查缓存是否过期
	if time.Since(manager.lastLoad) < manager.cacheTTL {
		return nil
	}

	_ = os.MkdirAll(manager.serversDir, os.ModePerm)
	dir, err := os.ReadDir(manager.serversDir)
	if err != nil {
		return fmt.Errorf("failed to read server directory: %w", err)
	}

	// 创建新的服务器映射
	newServers := &sync.Map{}

	// 保存现有服务器的状态
	existingServers := make(map[string]*Server)
	manager.servers.Range(func(key, value interface{}) bool {
		server := value.(*Server)
		existingServers[key.(string)] = server
		return true
	})

	for _, file := range dir {
		if !file.IsDir() {
			continue
		}
		serverId := file.Name()

		// 检查是否存在现有服务器实例
		if existingServer, exists := existingServers[serverId]; exists {
			// 如果存在，保留原有实例
			newServers.Store(serverId, existingServer)
		} else {
			// 如果不存在，创建新实例
			server := NewServer(ServerConfig{
				ID:         serverId,
				RootDir:    path.Join(manager.serversDir, serverId),
				BackupsDir: path.Join(manager.backupsDir, serverId),
			})
			newServers.Store(serverId, server)
		}
	}

	// 原子性更新服务器列表
	manager.servers = newServers
	manager.lastLoad = time.Now()

	return nil
}

func (manager *Manager) GetServer(id string) (*Server, error) {
	// 检查缓存是否过期
	if time.Since(manager.lastLoad) >= manager.cacheTTL {
		if err := manager.ScanServers(); err != nil {
			return nil, err
		}
		manager.lastLoad = time.Now() // 更新最后加载时间
	}

	if server, ok := manager.servers.Load(id); ok {
		return server.(*Server), nil
	}
	return nil, fmt.Errorf("server %s not found", id)
}

func (manager *Manager) GetServers() map[string]*Server {
	// 检查缓存是否过期
	if time.Since(manager.lastLoad) >= manager.cacheTTL {
		if err := manager.ScanServers(); err != nil {
			logrus.WithError(err).Error("Failed to refresh server list")
		}
		manager.lastLoad = time.Now() // 更新最后加载时间
	}

	servers := make(map[string]*Server)
	manager.servers.Range(func(key, value interface{}) bool {
		servers[key.(string)] = value.(*Server)
		return true
	})
	return servers
}

// GetServerCreateStatus 获取服务器创建状态
func (manager *Manager) GetServerCreateStatus(id string) (async.TaskStatus, error) {
	// 检查服务器是否存在
	if _, ok := manager.servers.Load(id); !ok {
		return async.TaskStatus{}, fmt.Errorf("server %s not found", id)
	}

	// 获取执行器
	value, ok := manager.createServerExecutors.Load(id)
	if !ok {
		return async.TaskStatus{}, fmt.Errorf("no create task found for server %s", id)
	}

	executor := value.(*async.AsyncTaskExecutor)
	return executor.GetCurrentStatus(), nil
}

// CreateServer 创建新服务器
func (manager *Manager) StartCreateServer(cfg CreateServerConfig) (*Server, error) {
	// 生成或验证服务器ID
	serverID := cfg.ID
	if serverID == "" {
		// 生成唯一的服务器ID
		for i := 1; i <= 1000; i++ { // 最多尝试1000次
			serverID = fmt.Sprintf("%d", time.Now().UnixNano())
			if _, exists := manager.servers.Load(serverID); !exists {
				break
			}
			time.Sleep(time.Millisecond) // 等待1毫秒后重试
		}
		if serverID == "" {
			return nil, fmt.Errorf("failed to generate unique server ID")
		}
	} else {
		// 验证ID是否已存在
		if _, exists := manager.servers.Load(serverID); exists {
			return nil, fmt.Errorf("server ID %s already exists", serverID)
		}
	}
	// 服务器目录
	serverDir := path.Join(manager.serversDir, serverID)

	// 设置服务器版本
	version := cfg.Version
	if version == "" {
		var err error
		version, err = manager.LatestVersion()
		if err != nil {
			return nil, fmt.Errorf("failed to get latest version: %w", err)
		}
	}

	// 服务器实例
	server := NewServer(ServerConfig{
		ID:         serverID,
		Name:       cfg.WorldName,
		RootDir:    serverDir,
		Version:    version,
		BackupsDir: path.Join(manager.backupsDir, serverID),
	})

	//服务器添加到管理器
	manager.servers.Store(serverID, server)

	//创建任务执行器
	executor := async.NewAsyncTaskExecutor()
	manager.createServerExecutors.Store(serverID, executor)

	executor.AddTask(async.NewSimpleTask("创建服务器目录", func() error {
		if err := os.MkdirAll(serverDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create server directory: %w", err)
		}

		return nil
	}))

	executor.AddTask(async.NewSimpleTask("下载版本文件", func() error {
		if err := manager.DownloadVersionFile(version); err != nil {
			return fmt.Errorf("failed to download version file: %w", err)
		}
		return nil
	}))

	executor.AddTask(async.NewSimpleTask("解压服务器文件", func() error {
		// 确保工作目录存在
		if err := os.MkdirAll(server.rootDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create server directory: %w", err)
		}

		// 解压服务器文件
		if err := utils.Unzip(manager.ZipFile(version), server.rootDir); err != nil {
			return fmt.Errorf("failed to unzip server files: %w", err)
		}

		return nil
	}))

	executor.AddTask(async.NewSimpleTask("配置服务器", func() error {
		if cfg.WorldName != DefaultWorldName {
			// 设置服务器名称
			serverProperties, err := server.GetServerProperties()
			if err != nil {
				return fmt.Errorf("failed to get server properties: %w", err)
			}
			if err := serverProperties.SetServerName(cfg.WorldName); err != nil {
				return fmt.Errorf("failed to set server name: %w", err)
			}
		}

		return nil
	}))
	executor.SetErrorCallback(func(task async.Task, err error) {
		logrus.WithError(err).WithField("server_id", serverID).Error("Failed to create server")
		// 删除服务器目录
		os.RemoveAll(serverDir)
		// 从管理器中移除服务器
		manager.servers.Delete(serverID)
		// 删除任务执行器
		manager.createServerExecutors.Delete(serverID)
	})

	// 开始执行任务
	executor.Start()

	return server, nil
}

// DeleteServer 删除服务器
func (manager *Manager) DeleteServer(id string) error {
	// 1. 获取服务器实例
	server, err := manager.GetServer(id)
	if err != nil {
		return fmt.Errorf("failed to get server: %w", err)
	}

	// 2. 如果服务器正在运行，先停止服务器
	if server.Active() {
		if err := server.Stop(); err != nil {
			return fmt.Errorf("failed to stop server: %w", err)
		}
	}

	// 3. 删除服务器目录
	serverDir := path.Join(manager.serversDir, id)
	if err := os.RemoveAll(serverDir); err != nil {
		return fmt.Errorf("failed to remove server directory: %w", err)
	}

	// 4. 从管理器中移除服务器
	manager.servers.Delete(id)

	// 5. 删除任务执行器
	manager.createServerExecutors.Delete(id)

	// 6. 重新加载服务器列表
	if err := manager.ScanServers(); err != nil {
		return fmt.Errorf("failed to reload server list: %w", err)
	}

	return nil
}

// GetSaves 获取所有存档信息
func (manager *Manager) GetSaves() []domain.Save {
	// 检查缓存是否过期
	if time.Since(manager.lastLoad) >= manager.cacheTTL {
		if err := manager.ScanSaves(); err != nil {
			logrus.WithError(err).Error("Failed to refresh save list")
		}
		manager.lastLoad = time.Now() // 更新最后加载时间
	}

	var saves []domain.Save
	manager.saves.Range(func(key, value interface{}) bool {
		saves = append(saves, value.(domain.Save))
		return true
	})
	return saves
}

// GetSave 获取指定存档信息
func (manager *Manager) GetSave(name string) (domain.Save, error) {
	// 检查缓存是否过期
	if time.Since(manager.lastLoad) >= manager.cacheTTL {
		if err := manager.ScanSaves(); err != nil {
			logrus.WithError(err).Error("Failed to refresh save list")
		}
		manager.lastLoad = time.Now() // 更新最后加载时间
	}

	if value, ok := manager.saves.Load(name); ok {
		return value.(domain.Save), nil
	}
	return domain.Save{}, fmt.Errorf("save %s not found", name)
}

// StartScanSaves 启动定期扫描存档的任务
func (manager *Manager) StartScanSaves() {
	go func() {
		ticker := time.NewTicker(manager.loadInterval)
		defer ticker.Stop()

		for range ticker.C {
			if err := manager.ScanSaves(); err != nil {
				logrus.WithError(err).Error("Failed to scan saves")
			}
		}
	}()
}

// ScanSaves 扫描存档目录
func (manager *Manager) ScanSaves() error {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	// 确保存档目录存在
	savesDir := manager.savesDir
	if err := os.MkdirAll(savesDir, 0755); err != nil {
		return fmt.Errorf("failed to create save directory: %w", err)
	}

	// 创建新的存档映射
	newSaves := &sync.Map{}

	err := filepath.Walk(savesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 处理 .mcworld 和 .zip 文件
		if !info.IsDir() {
			ext := filepath.Ext(path)
			if ext == ".mcworld" || ext == ".zip" {
				saveInfo := domain.Save{
					Name:         info.Name(),
					Size:         info.Size(),
					LastModified: info.ModTime(),
				}
				newSaves.Store(info.Name(), saveInfo)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to scan saves directory: %w", err)
	}

	// 原子性更新存档列表
	manager.saves = newSaves
	return nil
}

func (manager *Manager) AddSave(file multipart.File, header *multipart.FileHeader) (os.FileInfo, error) {
	//创建存档目录
	if err := os.MkdirAll(manager.savesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create saves directory: %w", err)
	}
	// 检查文件扩展名
	ext := filepath.Ext(header.Filename)
	if ext != ".mcworld" && ext != ".zip" {
		return nil, errors.New("unsupported file format")
	}
	// 如果是 mcworld 文件，将后缀改为 zip
	filename := header.Filename
	if ext == ".mcworld" {
		filename = filename[:len(filename)-len(ext)] + ".zip"
	}

	// 创建目标文件
	filepath := filepath.Join(manager.savesDir, filename)
	out, err := os.Create(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %w", filepath, err)
	}
	defer out.Close()

	// 将上传的文件内容复制到目标文件
	if _, err := io.Copy(out, file); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// 获取文件信息
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	err = manager.ScanSaves()
	if err != nil {
		return nil, fmt.Errorf("failed to scan saves directory: %w", err)
	}

	return fileInfo, nil
}

func (manager *Manager) ApplySave(saveName string, serverID string) error {
	// 获取服务器实例
	server, err := manager.GetServer(serverID)
	if err != nil {
		return fmt.Errorf("failed to get server: %w", err)
	}

	// 检查存档是否存在
	if _, err := manager.GetSave(saveName); err != nil {
		return fmt.Errorf("failed to get save: %w", err)
	}

	// 构建存档文件路径
	savePath := filepath.Join(manager.savesDir, saveName)

	// 获取世界数据目录
	worldDataDir := server.WorldDataDir()

	// 如果世界目录存在，先删除它
	if utils.FileExist(worldDataDir) {
		if err := os.RemoveAll(worldDataDir); err != nil {
			return fmt.Errorf("failed to remove world directory: %w", err)
		}
	}

	// 创建新的世界目录
	if err := os.MkdirAll(worldDataDir, 0755); err != nil {
		return fmt.Errorf("failed to create world directory: %w", err)
	}

	// 解压存档到世界目录
	if err := utils.Unzip(savePath, worldDataDir); err != nil {
		return fmt.Errorf("failed to unzip save file: %w", err)
	}

	return nil
}

func (manager *Manager) DeleteSave(name string) error {
	// 检查存档是否存在
	if _, ok := manager.saves.Load(name); !ok {
		return fmt.Errorf("save %s not found", name)
	}

	// 构建存档文件路径
	filepath := filepath.Join(manager.savesDir, name)

	// 删除文件
	if err := os.Remove(filepath); err != nil {
		return fmt.Errorf("failed to delete save file: %w", err)
	}

	// 从存档列表中移除
	manager.saves.Delete(name)

	// 重新扫描存档目录以更新列表
	if err := manager.ScanSaves(); err != nil {
		return fmt.Errorf("failed to scan saves directory: %w", err)
	}

	return nil
}
