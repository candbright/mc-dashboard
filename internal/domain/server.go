package domain

import (
	"strconv"
	"time"
)

type ServerInfo struct {
	Server
	Exist            bool        `json:"exist"`
	Active           bool        `json:"active"`
	CreateStatus     interface{} `json:"create_status"`
	ServerProperties interface{} `json:"server_properties"`
	AllowList        interface{} `json:"allow_list"`
}

// ServerService 服务器服务接口
type ServerService interface {
	ListServerInfos(page, size int, order, orderBy string) ([]ServerInfo, int64, error)
	GetServerInfo(id uint) (*ServerInfo, error)
	CreateServer(server *Server) error
	UpdateServer(server *Server) error
	StartServer(id uint) error
	StopServer(id uint) error
	DeleteServer(id uint) error
	GetServerProperties(id uint) (map[string]string, error)
	UpdateServerProperties(id uint, properties map[string]string) error
	GetAllowList(id uint) (AllowList, error)
	AddAllowListUser(id uint, username string) error
	DeleteAllowListUser(id uint, username string) error
	GetConsoleLog(id uint, line int) (string, error)
}

// Server 服务器信息
type Server struct {
	ID          uint      `gorm:"primarykey" json:"id"`                   // 主键
	Name        string    `gorm:"uniqueIndex;size:255" json:"name"`       // 服务器名称
	Description string    `gorm:"type:text;size:1000" json:"description"` // 服务器描述
	WorldName   string    `gorm:"type:text;size:255" json:"world_name"`   // 世界名称
	Version     string    `gorm:"type:text;size:15" json:"version"`       // 版本
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`       // 创建时间
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`       // 更新时间
}

func (Server) TableName() string {
	return "mc_dashboard_servers"
}

func (s *Server) IDString() string {
	return strconv.FormatUint(uint64(s.ID), 10)
}

// ServerRepository 服务器仓储接口
type ServerRepository interface {
	BaseRepository[Server]
	FindByName(name string) (*Server, error)
}
