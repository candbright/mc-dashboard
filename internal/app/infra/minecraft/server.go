package minecraft

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/candbright/mc-dashboard/internal/app/infra/minecraft/process"
	"github.com/candbright/mc-dashboard/internal/domain"
	"github.com/candbright/mc-dashboard/internal/pkg/utils"
	"github.com/candbright/mc-dashboard/pkg/dw"
	"github.com/sirupsen/logrus"
)

type ServerConfig struct {
	ID         string
	Name       string
	Version    string
	RootDir    string
	BackupsDir string
}

type Server struct {
	id               string // 服务器ID
	name             string // 服务器名称
	version          string // 服务器版本
	rootDir          string // 服务器根目录
	backupsDir       string // 备份目录
	worldsDir        string // 世界目录
	logsDir          string // 日志目录
	consoleLogFile   string // 控制台日志文件
	allowListFile    string // 白名单文件
	versionFile      string // 版本记录文件
	process          process.Process
	backupStatus     bool
	serverProperties *ServerProperties
}

func NewServer(cfg ServerConfig) *Server {
	logsDir := path.Join(cfg.RootDir, "logs")
	err := os.MkdirAll(logsDir, os.ModePerm)
	if err != nil {
		logrus.WithError(err).Error("make logs dir failed")
	}
	consoleLogFile := path.Join(logsDir, "console.log")

	server := &Server{
		id:             cfg.ID,
		name:           cfg.Name,
		version:        cfg.Version,
		rootDir:        cfg.RootDir,
		backupsDir:     cfg.BackupsDir,
		logsDir:        logsDir,
		consoleLogFile: consoleLogFile,
		worldsDir:      path.Join(cfg.RootDir, "worlds"),
		allowListFile:  path.Join(cfg.RootDir, "allowlist.json"),
		versionFile:    path.Join(cfg.RootDir, "version.txt"),
		process: process.NewProcess(process.ProcessConfig{
			ID:      cfg.ID,
			RootDir: cfg.RootDir,
			LogFile: consoleLogFile,
		}),
	}
	return server
}

func (server *Server) WorldDataDir() string {
	return path.Join(server.worldsDir, server.GetServerName())
}

func (server *Server) GetServerName() string {
	if server.name != "" {
		return server.name
	}
	serverProperties, err := server.GetServerProperties()
	if err != nil {
		return server.name
	}
	return serverProperties.GetServerName()
}

func (server *Server) SetServerName(name string) error {
	serverProperties, err := server.GetServerProperties()
	if err != nil {
		return err
	}
	return serverProperties.SetServerName(name)
}

func (server *Server) GetID() string {
	return server.id
}

func (server *Server) GetVersion() string {
	if server.version != "" {
		return server.version
	}
	bytes, err := os.ReadFile(server.versionFile)
	if err != nil {
		return server.version
	}
	server.version = string(bytes)
	return server.version
}

func (server *Server) Exist() bool {
	return server.process.Exist()
}

func (server *Server) Active() bool {
	return server.process.Active()
}

func (server *Server) GetServerProperties() (*ServerProperties, error) {
	if server.serverProperties == nil {
		var err error
		server.serverProperties, err = NewServerProperties(ServerPropertiesConfig{
			RootDir: server.rootDir,
		})
		if err != nil {
			return nil, err
		}
	}
	return server.serverProperties, nil
}

func (server *Server) Start() error {
	err := server.process.Start()
	if err != nil {
		return err
	}
	server.startBackup()
	return nil
}

func (server *Server) Stop() error {
	err := server.process.Stop()
	if err != nil {
		return err
	}
	server.stopBackup()
	return nil
}

func (server *Server) Delete() error {
	//如果服务器正在运行则先关闭
	active := server.process.Active()
	if active {
		err := server.process.Stop()
		if err != nil {
			return err
		}
	}
	//TODO: 备份

	//删除服务器目录
	existS := server.Exist()
	if existS {
		err := os.RemoveAll(server.rootDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func (server *Server) startBackup() {
	if !server.backupStatus {
		server.backupStatus = true
		go func() {
			for server.backupStatus {
				server.backupTick()
				time.Sleep(time.Hour * 1)
			}
		}()
		go func() {
			time.Sleep(time.Hour * 1)
			for server.backupStatus {
				server.backupClearTick()
				time.Sleep(time.Hour * 24)
			}
		}()
	}
}

func (server *Server) stopBackup() {
	server.backupStatus = false
}

func (server *Server) backupTick() {
	if !server.Active() {
		return
	}
	//zip data
	sourceDir := server.worldsDir
	backupsDir := server.backupsDir
	backupFile := path.Join(backupsDir, fmt.Sprintf("%s-%s.zip", server.name, time.Now().Format("20060102-150405")))
	err := os.MkdirAll(backupsDir, os.ModePerm)
	if err != nil {
		logrus.WithError(err).Error("make backup dir failed")
		return
	}
	err = utils.Zip(sourceDir, backupFile)
	if err != nil {
		logrus.WithError(err).Error("zip failed")
		return
	}
	logrus.Infof("backup has been saved to: %s\n", backupFile)
}

func (server *Server) backupClearTick() {
	if !server.Active() {
		return
	}
	// read dir
	files, err := os.ReadDir(server.backupsDir)
	if err != nil {
		logrus.WithError(err).Error("read backup dir failed")
		return
	}
	getFileTime := func(name string) string {
		return strings.Split(strings.Split(name, "-")[1], ".")[0]
	}
	// sort files by time
	sort.Slice(files, func(i, j int) bool {
		fi1, _ := os.Stat(getFileTime(files[i].Name()))
		fi2, _ := os.Stat(getFileTime(files[j].Name()))

		return fi1.ModTime().Before(fi2.ModTime())
	})

	// remove previous backup file
	for i := 0; i < len(files)-24; i++ {
		e := os.Remove(files[i].Name())
		if e != nil {
			logrus.WithError(e).Error("remove previous backup file failed")
		}
	}
}

func (server *Server) AllowListAdd(username string) error {
	if !server.Active() {
		return errors.New("服务器未启动")
	}
	return server.process.ExecCmd("allowlist", "add", username)
}

func (server *Server) AllowListDelete(username string) error {
	if !server.Active() {
		return errors.New("服务器未启动")
	}
	return server.process.ExecCmd("allowlist", "remove", username)
}

func (server *Server) AllowListOn() error {
	return server.process.ExecCmd("allowlist", "on")
}

func (server *Server) AllowListOff() error {
	return server.process.ExecCmd("allowlist", "off")
}

func (server *Server) GetAllowList() (domain.AllowList, error) {
	w, err := dw.Default[domain.AllowList](server.allowListFile)
	if err != nil {
		return nil, err
	}
	return w.Data, nil
}

// ApplySave 将存档应用到服务器
func (server *Server) ApplySave(savePath string) error {
	// 检查服务器是否存在
	if !server.Exist() {
		return errors.New("服务器不存在")
	}

	// 检查存档文件是否存在
	if !utils.FileExist(savePath) {
		return errors.New("存档文件不存在")
	}

	// 获取世界数据目录
	worldDataDir := server.WorldDataDir()

	// 如果世界目录存在，先删除它
	if utils.FileExist(worldDataDir) {
		if err := os.RemoveAll(worldDataDir); err != nil {
			return err
		}
	}

	// 创建新的世界目录
	if err := os.MkdirAll(worldDataDir, 0755); err != nil {
		return err
	}

	// 解压存档到世界目录
	if err := utils.Unzip(savePath, worldDataDir); err != nil {
		return err
	}

	return nil
}

// ScanLog 获取服务器日志
func (server *Server) ScanLog(line int) (string, error) {
	if server.process == nil {
		return "", fmt.Errorf("server process not found")
	}
	return server.process.ScanLog(line)
}
