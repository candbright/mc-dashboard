package core

import "testing"

func TestWindow_SendMessage(t *testing.T) {
	window := Window{title: "E:\\minecraft\\server-1\\1.21.62.01\\bedrock_server.exe"}
	err := window.ExecuteCommand("allowlist list")
	if err != nil {
		t.Error(err)
	}
}
