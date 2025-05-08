package model

type ServerInfo struct {
	ID               string      `json:"id"`
	Name             string      `json:"name"`
	Version          string      `json:"version"`
	Exist            bool        `json:"exist"`
	Downloading      bool        `json:"downloading"`
	Active           bool        `json:"active"`
	ServerProperties interface{} `json:"server_properties"`
	AllowList        interface{} `json:"allow_list"`
}
