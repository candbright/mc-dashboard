package core

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/candbright/go-log/log"
	"github.com/candbright/go-server/internal/mc-server/core/model"
	"github.com/candbright/go-server/pkg/config"
	"github.com/candbright/go-server/pkg/dw"
	"github.com/pkg/errors"
)

type ServerConfig struct {
	ID      string
	RootDir string
}

type Server struct {
	id               string
	name             string
	version          string
	rootDir          string
	process          *Process
	backup           bool
	serverProperties *ServerProperties
}

func NewServer(cfg ServerConfig) (*Server, error) {
	server := &Server{
		id:      cfg.ID,
		rootDir: cfg.RootDir,
	}
	server.process = NewProcess(ProcessConfig{
		RootDir: server.WorkDir(),
	})
	return server, nil
}

func (server *Server) BackupDir() string {
	return path.Join(config.Global.Get("mc.path"), "backup", fmt.Sprintf("backup-"+server.rootDir))
}

func (server *Server) VersionFilePath() string {
	return path.Join(server.rootDir, "version")
}

func (server *Server) DownloadingFilePath() string {
	return path.Join(server.rootDir, "downloading")
}

func (server *Server) WorkDir() string {
	return path.Join(server.rootDir, server.GetVersion())
}

func (server *Server) WorldsDir() string {
	return path.Join(server.WorkDir(), "worlds")
}

func (server *Server) WorldDataDir() string {
	serverName := server.serverProperties.Data["server-name"]
	return path.Join(server.WorkDir(), "worlds", serverName)
}

func (server *Server) AllowListFilePath() string {
	return path.Join(server.WorkDir(), "allowlist.json")
}

func (server *Server) GetServerName() string {
	if server.name != "" {
		return server.name
	}
	bytes, err := os.ReadFile(path.Join(server.rootDir, "servername"))
	if err != nil {
		server.name = "unnamed server"
	}
	if len(bytes) > 128 {
		bytes = bytes[:128]
	}
	server.name = string(bytes)
	return server.name
}

func (server *Server) GetID() string {
	return server.id
}

func (server *Server) GetVersion() string {
	if server.version != "" {
		return server.version
	}
	bytes, err := os.ReadFile(server.VersionFilePath())
	if err != nil {
		server.version = ""
		return server.version
	}
	server.version = string(bytes)
	return server.version
}

func (server *Server) ServerExist() bool {
	//是否正在下载
	downloading := server.Downloading()
	if downloading {
		return false
	}
	//版本文件是否存在
	versionExists := Exists(server.VersionFilePath())
	if !versionExists {
		return false
	}
	//服务器目录是否存在
	return Exists(server.WorkDir())
}

func (server *Server) Downloading() bool {
	return Exists(server.DownloadingFilePath())
}

func (server *Server) Active() bool {
	exist := server.ServerExist()
	if !exist {
		return false
	}
	return server.process.Active()
}

func (server *Server) ServerProperties() (*ServerProperties, error) {
	exist := server.ServerExist()
	if !exist {
		return nil, errors.New("server not exist")
	}
	if server.serverProperties == nil {
		server.serverProperties = NewServerProperties(ServerPropertiesConfig{
			Version: server.version,
			RootDir: server.WorkDir(),
		})
	}
	return server.serverProperties, nil
}

func (server *Server) Start() error {
	return server.process.Start()
}

func (server *Server) Stop() error {
	return server.process.Stop()
}

func (server *Server) Reload() error {
	needStart := false
	if server.Active() {
		err := server.process.Stop()
		if err != nil {
			return err
		}
		needStart = true
	}
	server.process = NewProcess(ProcessConfig{
		RootDir: server.WorkDir(),
	})
	if needStart {
		err := server.process.Start()
		if err != nil {
			return err
		}
	}
	server.serverProperties = NewServerProperties(ServerPropertiesConfig{
		Version: server.version,
		RootDir: server.WorkDir(),
	})

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
	existS := server.ServerExist()
	if existS {
		err := os.RemoveAll(server.WorkDir())
		if err != nil {
			return err
		}
	}
	return nil
}

func (server *Server) startBackup() {
	if !server.backup {
		server.backup = true
		go func() {
			for server.backup {
				server.backupTick()
				time.Sleep(time.Hour * 1)
			}
		}()
		go func() {
			time.Sleep(time.Hour * 1)
			for server.backup {
				server.backupClearTick()
				time.Sleep(time.Hour * 24)
			}
		}()
	}
}

func (server *Server) stopBackup() {
	server.backup = false
}

func (server *Server) backupTick() {
	if !server.Active() {
		return
	}
	//zip data
	sourceDir := server.WorldsDir()
	backupDir := server.BackupDir()
	backupFile := fmt.Sprintf("%s/backup-%s.zip", backupDir, time.Now().Format("20060102-150405"))
	err := os.MkdirAll(backupDir, os.ModePerm)
	if err != nil {
		log.WithError(err).Error("make backup dir failed")
		return
	}
	err = Zip(backupFile, sourceDir)
	if err != nil {
		log.WithError(err).Error("zip failed")
		return
	}
	log.Infof("backup has been saved to: %s\n", backupFile)
}

func (server *Server) backupClearTick() {
	if !server.Active() {
		return
	}
	// read dir
	files, err := os.ReadDir(server.BackupDir())
	if err != nil {
		log.WithError(err).Error("read backup dir failed")
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
			log.WithError(e).Error("remove previous backup file failed")
		}
	}
}

func (server *Server) AllowListAdd(username string) error {
	return server.process.ExecCmd("allowlist", "add", username)
}

func (server *Server) AllowListDelete(username string) error {
	return server.process.ExecCmd("allowlist", "remove", username)
}

func (server *Server) AllowListOn() error {
	return server.process.ExecCmd("allowlist", "on")
}

func (server *Server) AllowListOff() error {
	return server.process.ExecCmd("allowlist", "off")
}

func (server *Server) GetAllowList() (model.AllowList, error) {
	w, err := dw.Default[model.AllowList](server.AllowListFilePath())
	if err != nil {
		return nil, err
	}
	return w.Data, nil
}

// ApplySave 将存档应用到服务器
func (server *Server) ApplySave(savePath string) error {
	// 检查服务器是否存在
	if !server.ServerExist() {
		return errors.New("服务器不存在")
	}

	// 检查存档文件是否存在
	if !Exists(savePath) {
		return errors.New("存档文件不存在")
	}

	// 获取世界数据目录
	worldDataDir := server.WorldDataDir()

	// 如果世界目录存在，先删除它
	if Exists(worldDataDir) {
		if err := os.RemoveAll(worldDataDir); err != nil {
			return errors.WithStack(err)
		}
	}

	// 创建新的世界目录
	if err := os.MkdirAll(worldDataDir, 0755); err != nil {
		return errors.WithStack(err)
	}

	// 解压存档到世界目录
	if err := Unzip(savePath, worldDataDir); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
