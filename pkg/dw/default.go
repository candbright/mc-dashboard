package dw

import (
	"encoding/json"
	"encoding/xml"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"
)

func Default[T any](path string) (*DataWriter[T], error) {
	return Json[T](path)
}

func Json[T any](path string) (*DataWriter[T], error) {
	cfg := Config{
		Path:      path,
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
	}
	return New[T](cfg)
}

func Xml[T any](path string) (*DataWriter[T], error) {
	cfg := Config{
		Path:      path,
		Marshal:   xml.Marshal,
		Unmarshal: xml.Unmarshal,
	}
	return New[T](cfg)
}

func Yaml[T any](path string) (*DataWriter[T], error) {
	cfg := Config{
		Path:      path,
		Marshal:   yaml.Marshal,
		Unmarshal: yaml.Unmarshal,
	}
	return New[T](cfg)
}

func Toml[T any](path string) (*DataWriter[T], error) {
	cfg := Config{
		Path:      path,
		Marshal:   toml.Marshal,
		Unmarshal: toml.Unmarshal,
	}
	return New[T](cfg)
}
