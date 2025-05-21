package usecase

import (
	"fmt"
	"strconv"

	"github.com/candbright/mc-dashboard/internal/app/infra/minecraft"
	"github.com/candbright/mc-dashboard/internal/domain"
	"github.com/sirupsen/logrus"
)

type ServerService struct {
	log       *logrus.Logger
	mcManager *minecraft.Manager
	repo      domain.ServerRepository
}

// NewService 创建服务器服务实例
func NewServerService(log *logrus.Logger, mcManager *minecraft.Manager, repo domain.ServerRepository) domain.ServerService {
	return &ServerService{
		log:       log,
		mcManager: mcManager,
		repo:      repo,
	}
}

// ListServers 获取服务器列表
func (s *ServerService) ListServerInfos(page, size int, order, orderBy string) ([]domain.ServerInfo, int64, error) {
	servers, total, err := s.repo.List(page, size, order, orderBy)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get servers from repository: %w", err)
	}
	serverInfos := make([]domain.ServerInfo, 0, len(servers))
	for _, server := range servers {
		serverInfo, err := s.mcManager.GetServer(server.IDString())
		if err != nil {
			serverInfos = append(serverInfos, domain.ServerInfo{
				Server: server,
				Exist:  false,
			})
		} else {
			serverInfos = append(serverInfos, *s.parseServerInfo(&server, serverInfo))
		}
	}

	return serverInfos, total, nil
}

// GetServer 获取服务器信息
func (s *ServerService) GetServerInfo(id uint) (*domain.ServerInfo, error) {
	server, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get server from repository: %w", err)
	}
	serverInfo, err := s.mcManager.GetServer(server.IDString())
	if err != nil {
		return nil, fmt.Errorf("failed to get server: %w", err)
	}
	info := s.parseServerInfo(server, serverInfo)
	return info, nil
}

// CreateServer 创建服务器
func (s *ServerService) CreateServer(server *domain.Server) error {
	if server.Name == "" {
		return fmt.Errorf("server name is required")
	}

	if server.WorldName == "" {
		server.WorldName = "Dedicated Server"
	}

	if server.Version == "" {
		var err error
		server.Version, err = s.mcManager.LatestVersion()
		if err != nil {
			return fmt.Errorf("failed to get latest version: %w", err)
		}
	}

	//添加到数据库
	err := s.repo.Create(server)
	if err != nil {
		return fmt.Errorf("failed to create server in repository: %w", err)
	}

	if server.ID == 0 {
		// 重新查询获取自增ID
		created, err := s.repo.FindByName(server.Name)
		if err != nil {
			return fmt.Errorf("failed to get created server from repository: %w", err)
		}
		server.ID = created.ID
	}

	_, err = s.mcManager.StartCreateServer(minecraft.CreateServerConfig{
		ID:        server.IDString(),
		WorldName: server.WorldName,
		Version:   server.Version,
	})
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	return nil
}

// UpdateServer 更新服务器
func (s *ServerService) UpdateServer(server *domain.Server) error {
	oldServer, err := s.repo.FindByID(server.ID)
	if err != nil {
		return fmt.Errorf("failed to get server from repository: %w", err)
	}

	if server.WorldName == "" {
		server.WorldName = oldServer.WorldName
	}
	if oldServer.WorldName != server.WorldName {
		return fmt.Errorf("world name is not allowed to be changed")
	}

	if server.Version == "" {
		server.Version = oldServer.Version
	}
	if oldServer.Version != server.Version {
		return fmt.Errorf("version is not allowed to be changed")
	}

	err = s.repo.Update(server)
	if err != nil {
		return fmt.Errorf("failed to update server in repository: %w", err)
	}

	return nil
}

// StartServer 启动服务器
func (s *ServerService) StartServer(id uint) error {
	server, err := s.mcManager.GetServer(strconv.FormatUint(uint64(id), 10))
	if err != nil {
		return fmt.Errorf("failed to get server: %w", err)
	}

	if server.Active() {
		return fmt.Errorf("server is already running")
	}

	return server.Start()
}

// StopServer 停止服务器
func (s *ServerService) StopServer(id uint) error {
	server, err := s.mcManager.GetServer(strconv.FormatUint(uint64(id), 10))
	if err != nil {
		return fmt.Errorf("failed to get server: %w", err)
	}

	if !server.Active() {
		return fmt.Errorf("server is not running")
	}

	return server.Stop()
}

// DeleteServer 删除服务器
func (s *ServerService) DeleteServer(id uint) error {
	server, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to get server: %w", err)
	}

	_, err = s.mcManager.GetServer(server.IDString())
	if err == nil {
		innerErr := s.mcManager.DeleteServer(strconv.FormatUint(uint64(id), 10))
		if innerErr != nil {
			logrus.Errorf("failed to delete server: %w", innerErr)
		}
	}

	err = s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete server from repository: %w", err)
	}

	return nil
}

// GetServerProperties 获取服务器属性
func (s *ServerService) GetServerProperties(id uint) (map[string]string, error) {
	server, err := s.mcManager.GetServer(strconv.FormatUint(uint64(id), 10))
	if err != nil {
		return nil, fmt.Errorf("failed to get server: %w", err)
	}

	properties, err := server.GetServerProperties()
	if err != nil {
		return nil, fmt.Errorf("failed to get server properties: %w", err)
	}

	return properties.GetAll(), nil
}

// UpdateServerProperties 更新服务器属性
func (s *ServerService) UpdateServerProperties(id uint, properties map[string]string) error {
	server, err := s.mcManager.GetServer(strconv.FormatUint(uint64(id), 10))
	if err != nil {
		return fmt.Errorf("failed to get server: %w", err)
	}

	serverProperties, err := server.GetServerProperties()
	if err != nil {
		return fmt.Errorf("failed to get server properties: %w", err)
	}

	// 获取当前配置
	currentConfig := serverProperties.GetAll()

	for k, v := range properties {
		if _, ok := currentConfig[k]; ok {
			currentConfig[k] = v
		}
	}

	err = serverProperties.SetAll(currentConfig)
	if err != nil {
		return fmt.Errorf("failed to set server properties: %w", err)
	}

	return nil
}

// GetAllowList 获取白名单列表
func (s *ServerService) GetAllowList(id uint) (domain.AllowList, error) {
	server, err := s.mcManager.GetServer(strconv.FormatUint(uint64(id), 10))
	if err != nil {
		return nil, fmt.Errorf("failed to get server: %w", err)
	}

	allowList, err := server.GetAllowList()
	if err != nil {
		return nil, fmt.Errorf("failed to get allow list: %w", err)
	}

	return allowList, nil
}

// AddAllowListUser 添加白名单用户
func (s *ServerService) AddAllowListUser(id uint, username string) error {
	server, err := s.mcManager.GetServer(strconv.FormatUint(uint64(id), 10))
	if err != nil {
		return fmt.Errorf("failed to get server: %w", err)
	}
	err = server.AllowListAdd(username)
	if err != nil {
		return fmt.Errorf("failed to add user to allow list: %w", err)
	}

	return nil
}

// DeleteAllowListUser 删除白名单用户
func (s *ServerService) DeleteAllowListUser(id uint, username string) error {
	server, err := s.mcManager.GetServer(strconv.FormatUint(uint64(id), 10))
	if err != nil {
		return fmt.Errorf("failed to get server: %w", err)
	}

	err = server.AllowListDelete(username)
	if err != nil {
		return fmt.Errorf("failed to remove user from allow list: %w", err)
	}

	return nil
}

func (s *ServerService) parseServerInfo(server *domain.Server, runningServer *minecraft.Server) *domain.ServerInfo {
	info := &domain.ServerInfo{
		Server: *server,
		Exist:  runningServer.Exist(),
		Active: runningServer.Active(),
	}
	createStatus, err := s.mcManager.GetServerCreateStatus(runningServer.GetID())
	if err == nil {
		info.CreateStatus = createStatus
	}

	serverProperties, _ := runningServer.GetServerProperties()
	if serverProperties != nil {
		info.ServerProperties = serverProperties.GetAll()
	}

	allowList, _ := runningServer.GetAllowList()
	if allowList != nil {
		info.AllowList = allowList
	}

	return info
}

// GetConsoleLog 获取服务器控制台日志
func (s *ServerService) GetConsoleLog(id uint, line int) (string, error) {
	// 获取服务器实例
	server, err := s.mcManager.GetServer(strconv.FormatUint(uint64(id), 10))
	if err != nil {
		return "", fmt.Errorf("failed to get server: %w", err)
	}

	// 获取日志内容
	logContent, err := server.ScanLog(line)
	if err != nil {
		return "", fmt.Errorf("failed to scan log: %w", err)
	}

	return logContent, nil
}
