package config

import (
	"fmt"
	"testing"
)

var testData = `
application:
  name: 'test'
  port:
    default: 10088
    env: ${CORE_PORT}
  dir:
    x86: '/root'
    arm: '/home'
  osdir:
    linux:
      env: ${CORE_OSDIR_LINUX}
      default: '/root'
    windows: 'D:/home'
release: true
db:
  default: 'tcp:127.0.0.1:3306'
  env: ${CORE_DB}`

func TestAppConfig_Get(t *testing.T) {
	cfg, err := Parse([]byte(testData), YAML)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(cfg.Get("application.name"))
	fmt.Println(cfg.Get("application.port"))
	fmt.Println(cfg.Get("application.dir"))
	fmt.Println(cfg.Get("application.osdir"))
	fmt.Println(cfg.Get("db"))
	fmt.Println(cfg.Get("release"))
}
