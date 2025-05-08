package route

import (
	"time"

	"github.com/candbright/go-server/internal/mc-server/core"
	"github.com/candbright/go-server/pkg/config"
)

var manager *core.ServerManager

func Init() {
	manager = core.NewServersManager(
		core.ServerManagerConfig{
			RootDir:      config.Global.Get("mc.path"),
			LoadInterval: 1 * time.Minute,
		},
	)
}
