package core

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/candbright/go-log/log"
	"github.com/candbright/go-server/pkg/downloader"
)

type ServerManagerConfig struct {
	RootDir      string
	LoadInterval time.Duration
	CacheTTL     time.Duration
}

type ServerManager struct {
	rootDir      string
	servers      *sync.Map //key: id
	downloaders  *sync.Map //key: version
	loadInterval time.Duration
	cacheTTL     time.Duration
	lastLoad     time.Time
	mu           sync.RWMutex
}

func NewServersManager(cfg ServerManagerConfig) *ServerManager {
	if cfg.LoadInterval == 0 {
		cfg.LoadInterval = time.Minute
	}
	if cfg.CacheTTL == 0 {
		cfg.CacheTTL = 5 * time.Minute
	}

	ssm := &ServerManager{
		rootDir:      cfg.RootDir,
		servers:      &sync.Map{},
		downloaders:  &sync.Map{},
		loadInterval: cfg.LoadInterval,
		cacheTTL:     cfg.CacheTTL,
		lastLoad:     time.Now().Add(-cfg.CacheTTL), // 设置为一个已经过期的时间，确保第一次加载会执行
	}

	// 初始化时加载一次
	if err := ssm.LoadServers(); err != nil {
		log.WithError(err).Error("Failed to load servers during initialization")
	}

	// 启动定期加载
	ssm.StartLoadServers()
	return ssm
}

func (manager *ServerManager) StartLoadServers() {
	go func() {
		ticker := time.NewTicker(manager.loadInterval)
		defer ticker.Stop()

		for range ticker.C {
			if err := manager.LoadServers(); err != nil {
				log.WithError(err).Error("Failed to load servers")
			}
		}
	}()
}

func (manager *ServerManager) VersionsDir() string {
	return path.Join(manager.rootDir, "versions")
}

func (manager *ServerManager) LatestVersion() (string, error) {
	//TODO
	/*resp, err := m.client.R().
		Get("https://www.minecraft.net/en-us/download/server/bedrock")
	if err != nil {
		return "", errors.WithStack(err)
	}
	return string(resp.Body()), nil*/
	return "1.21.62.01", nil
}

func (manager *ServerManager) ZipFileName(version string) string {
	return fmt.Sprintf("bedrock-server-%s.zip", version)
}

func (manager *ServerManager) ZipFile(version string) string {
	return path.Join(manager.VersionsDir(), manager.ZipFileName(version))
}

func (manager *ServerManager) ZipExist(version string) bool {
	//zip文件是否存在
	return Exists(manager.ZipFile(version))
}

func (manager *ServerManager) DownloadLatestVersion() error {
	latestVersion, err := manager.LatestVersion()
	if err != nil {
		return err
	}
	return manager.DownloadVersion(latestVersion)
}

func (manager *ServerManager) DownloadVersion(version string) error {
	//检测是否存在当前版本的zip文件
	existZ := manager.ZipExist(version)
	if existZ {
		return nil
	}
	//不存在当前版本的zip文件，开始下载
	if _, ok := manager.downloaders.Load(version); !ok {
		manager.downloaders.Store(version, downloader.NewDownloader())
	}
	value, _ := manager.downloaders.Load(version)
	d := value.(*downloader.Downloader)
	status := d.GetCurrentStatus()
	if status.IsDownloading {
		return nil
	} else {
		var downloadUrl string
		switch runtime.GOOS {
		case "linux":
			downloadUrl = fmt.Sprintf("https://www.minecraft.net/bedrockdedicatedserver/bin-linux/bedrock-server-%s.zip", version)
		case "windows":
			downloadUrl = fmt.Sprintf("https://www.minecraft.net/bedrockdedicatedserver/bin-win/bedrock-server-%s.zip", version)
		}
		err := os.MkdirAll(manager.VersionsDir(), os.ModePerm)
		if err != nil {
			return err
		}
		d.Download(downloadUrl, manager.ZipFile(version))
		return nil
	}
}

func (manager *ServerManager) LoadServers() error {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	// 检查缓存是否过期
	if time.Since(manager.lastLoad) < manager.cacheTTL {
		return nil
	}

	_ = os.MkdirAll(manager.rootDir, os.ModePerm)
	dir, err := os.ReadDir(manager.rootDir)
	if err != nil {
		return fmt.Errorf("failed to read server directory: %w", err)
	}

	// 创建新的服务器映射
	newServers := &sync.Map{}

	for _, file := range dir {
		if !file.IsDir() {
			continue
		}

		prefix := "server-"
		if !strings.HasPrefix(file.Name(), prefix) {
			continue
		}

		idStr := file.Name()[len(prefix):]
		server, err := NewServer(ServerConfig{
			ID:      idStr,
			RootDir: path.Join(manager.rootDir, file.Name()),
		})
		if err != nil {
			log.WithError(err).WithField("server_id", idStr).Error("Failed to create server")
			continue
		}

		newServers.Store(idStr, server)
	}

	// 原子性更新服务器列表
	manager.servers = newServers
	manager.lastLoad = time.Now()

	return nil
}

func (manager *ServerManager) GetServer(id string) (*Server, error) {
	// 检查缓存是否过期
	if time.Since(manager.lastLoad) >= manager.cacheTTL {
		if err := manager.LoadServers(); err != nil {
			return nil, err
		}
		manager.lastLoad = time.Now() // 更新最后加载时间
	}

	if server, ok := manager.servers.Load(id); ok {
		return server.(*Server), nil
	}
	return nil, fmt.Errorf("server [%s] not found", id)
}

func (manager *ServerManager) GetServers() map[string]*Server {
	// 检查缓存是否过期
	if time.Since(manager.lastLoad) >= manager.cacheTTL {
		if err := manager.LoadServers(); err != nil {
			log.WithError(err).Error("Failed to refresh server list")
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

func (manager *ServerManager) DownloadServer(id string, version string) error {
	server, ok := manager.servers.Load(id)
	if !ok {
		return fmt.Errorf("server %s not found", id)
	}

	if version == "" {
		//获取当前最新版本
		latestVersion, err := manager.LatestVersion()
		if err != nil {
			return err
		}
		version = latestVersion
		s := server.(*Server)
		s.version = latestVersion
	}

	//下载版本压缩包
	err := manager.DownloadVersion(version)
	if err != nil {
		return err
	}

	//判断当前任务是否正在执行
	s := server.(*Server)
	downloading := s.Downloading()
	if downloading {
		return nil
	}
	_, err = os.Create(s.DownloadingFilePath())
	if err != nil {
		return err
	}
	defer func() {
		//删除lock文件
		_ = os.Remove(s.DownloadingFilePath())
	}()
	//解压zip文件
	err = os.MkdirAll(s.WorkDir(), os.ModePerm)
	if err != nil {
		return err
	}
	err = Unzip(manager.ZipFile(version), s.WorkDir())
	if err != nil {
		return err
	}
	//reload
	err = s.Reload()
	if err != nil {
		return err
	}
	err = os.WriteFile(s.VersionFilePath(), []byte(version), 0666)
	if err != nil {
		return err
	}

	return nil
}

/*
func (manager *ServerManager) Upgrade() error {
	//1. 获取最新版本
	oldServer := manager.current
	newVersion, err := manager.LatestVersion()
	if err != nil {
		return err
	}
	//2. 若最新版本和当前版本不同，则下载最新版本
	if newVersion == oldServer.version {
		return nil
	}
	newServer, err := NewServer(ServerConfig{
		GetVersion: newVersion,
		RootDir: oldServer.rootDir,
		Session: manager.session,
	})
	if err != nil {
		return err
	}
	err = newServer.Download()
	if err != nil {
		return err
	}
	manager.current = newServer
	//3. 复制旧版本数据文件到新版本
	err = manager.session.Run("cp", "-r",
		path.Join(oldServer.WorkDir(), "world"),
		path.Join(manager.current.WorkDir()+"/"))
	if err != nil {
		return err
	}
	err = manager.session.Run("cp",
		path.Join(oldServer.WorkDir(), oldServer.allowList.FileName()),
		path.Join(manager.current.WorkDir(), manager.current.allowList.FileName()))
	if err != nil {
		return err
	}
	err = manager.session.Run("cp",
		path.Join(oldServer.WorkDir(), oldServer.serverProperties.FileName()),
		path.Join(manager.current.WorkDir(), manager.current.serverProperties.FileName()))
	if err != nil {
		return err
	}
	return nil
}
*/
