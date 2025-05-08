package config

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

const YAML ParseType = "YAML"
const XML ParseType = "XML"
const JSON ParseType = "JSON"

type ParseType string

type Config struct {
	originData []byte
	parseType  ParseType
	data       map[interface{}]interface{}
}

func Parse(data []byte, parseType ParseType) (*Config, error) {
	config := &Config{
		originData: data,
		parseType:  parseType,
	}
	err := config.parse()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func ParseFile(fileName string, parseType ParseType) (*Config, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	config := &Config{
		originData: data,
		parseType:  parseType,
	}
	err = config.parse()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) Reload() error {
	return c.parse()
}

func (c *Config) parse() error {
	var err error
	switch c.parseType {
	case YAML:
		err = yaml.Unmarshal(c.originData, &c.data)
		if err != nil {
			return err
		}
		break
	case XML:
		err = xml.Unmarshal(c.originData, &c.data)
		if err != nil {
			return err
		}
		break
	case JSON:
		err = json.Unmarshal(c.originData, &c.data)
		if err != nil {
			return err
		}
		break
	default:
		return errors.New("unsupported parse type")
	}
	return nil
}

func (c *Config) Get(key string) string {
	if c.data == nil {
		fmt.Println("config data is nil, please parse first")
		return ""
	}
	split := strings.Split(key, ".")
	value := get(split, c.data)
	if value == nil {
		fmt.Println("value is nil, please input the right key")
		return ""
	}
	if val, ok := value.(string); ok {
		return val
	}
	if valMap, ok := value.(map[interface{}]interface{}); ok {
		env := envParser(valMap)
		if env != nil {
			return env.Get()
		}
		archV := archParser(valMap)
		if archV != nil {
			return archV.Get()
		}
		osV := osParser(valMap)
		if osV != nil {
			return osV.Get()
		}
	}
	if val, ok := value.(interface{}); ok {
		return fmt.Sprint(val)
	}
	return ""
}

func (c *Config) GetInt(key string) int {
	return int(c.GetInt64(key))
}

func (c *Config) GetInt64(key string) int64 {
	if c.data == nil {
		fmt.Println("config data is nil, please parse first")
		return -1
	}
	split := strings.Split(key, ".")
	value := get(split, c.data)
	if value == nil {
		fmt.Println("value is nil, please input the right key")
		return -1
	}
	if val, ok := value.(int); ok {
		return int64(val)
	}
	if val, ok := value.(int64); ok {
		return val
	}
	if valMap, ok := value.(map[interface{}]interface{}); ok {
		env := envParser(valMap)
		if env != nil {
			return env.GetInt64()
		}
		archV := archParser(valMap)
		if archV != nil {
			return archV.GetInt64()
		}
		osV := osParser(valMap)
		if osV != nil {
			return osV.GetInt64()
		}
	}
	fmt.Println("value type is unsupported")
	return -1
}

func (c *Config) GetBool(key string) bool {
	if c.data == nil {
		fmt.Println("config data is nil, please parse first")
		return false
	}
	split := strings.Split(key, ".")
	value := get(split, c.data)
	if value == nil {
		fmt.Println("value is nil, please input the right key")
		return false
	}
	if val, ok := value.(bool); ok {
		return val
	}
	if valMap, ok := value.(map[interface{}]interface{}); ok {
		env := envParser(valMap)
		if env != nil {
			return env.GetBool()
		}
		archV := archParser(valMap)
		if archV != nil {
			return archV.GetBool()
		}
		osV := osParser(valMap)
		if osV != nil {
			return osV.GetBool()
		}
	}
	fmt.Println("value type is unsupported")
	return false
}

func get(keys []string, tree map[interface{}]interface{}) interface{} {
	if keys == nil {
		return nil
	}
	if len(keys) == 1 {
		value, ok := tree[keys[0]]
		if ok {
			return value
		} else {
			return nil
		}
	}
	value, ok := tree[keys[0]]
	if valueMap, valueOk := value.(map[interface{}]interface{}); valueOk && ok {
		return get(keys[1:], valueMap)
	} else {
		return nil
	}
}
