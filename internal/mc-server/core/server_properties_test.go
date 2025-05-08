package core

import (
	"testing"
)

func TestNewServerProperties_GetServerName(t *testing.T) {
	p := NewServerProperties(ServerPropertiesConfig{
		Version: "1.20.62.02",
		RootDir: "./example",
	})
	t.Log(p.GetServerName())
}
