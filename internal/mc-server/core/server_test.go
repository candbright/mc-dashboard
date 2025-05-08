package core

import (
	"testing"
)

const testRootDir = "/opt/minecraft"

func testServer(t *testing.T) *Server {
	server, err := NewServer(ServerConfig{
		RootDir: testRootDir,
	})
	if err != nil {
		t.Fatalf("%+v", err)
	}
	return server
}
