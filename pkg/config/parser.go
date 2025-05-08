package config

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type Parser interface {
	parser
	GetBool() bool
	GetInt() int
	GetInt64() int64
}

type parser interface {
	Get() string
}

type BaseParser struct {
	parser
}

func (p *BaseParser) GetBool() bool {
	if strings.ToLower(p.Get()) == "true" {
		return true
	}
	return false
}

func (p *BaseParser) GetInt() int {
	return int(p.GetInt64())
}

func (p *BaseParser) GetInt64() int64 {
	intVal, err := strconv.Atoi(p.Get())
	if err != nil {
		return 0
	}
	return int64(intVal)
}

type EnvParser struct {
	Default interface{} `yaml:"default"`
	Env     interface{} `yaml:"env"`
}

func (p *EnvParser) Get() string {
	env := os.ExpandEnv(fmt.Sprint(p.Env))
	if env != "" {
		return env
	}
	return fmt.Sprint(p.Default)
}

type ArchParser struct {
	X86 interface{} `yaml:"x86"`
	Arm interface{} `yaml:"arm"`
}

func (p *ArchParser) Get() string {
	if strings.Contains(runtime.GOARCH, "arm") {
		return fmt.Sprint(p.Arm)
	}
	return fmt.Sprint(p.X86)
}

type OsParser struct {
	Linux   interface{} `yaml:"linux"`
	Windows interface{} `yaml:"windows"`
}

func (value *OsParser) Get() string {
	if strings.Contains(runtime.GOOS, "linux") {
		return fmt.Sprint(value.Linux)
	}
	return fmt.Sprint(value.Windows)
}

func envParser(valMap map[interface{}]interface{}) Parser {
	if valMap == nil {
		return nil
	}
	Val1, ok1 := valMap["default"]
	Val2, ok2 := valMap["env"]
	if !ok1 || !ok2 {
		return nil
	}
	return &BaseParser{&EnvParser{Default: Val1, Env: Val2}}
}

func archParser(valMap map[interface{}]interface{}) Parser {
	if valMap == nil {
		return nil
	}
	Val1, ok1 := valMap["x86"]
	Val2, ok2 := valMap["arm"]
	if !ok1 || !ok2 {
		return nil
	}
	valE1, ok := Val1.(map[interface{}]interface{})
	if p := envParser(valE1); ok && p != nil {
		return p
	}
	valE2, ok := Val2.(map[interface{}]interface{})
	if p := envParser(valE2); ok && p != nil {
		return p
	}
	return &BaseParser{&ArchParser{X86: Val1, Arm: Val2}}
}

func osParser(valMap map[interface{}]interface{}) Parser {
	Val1, ok1 := valMap["linux"]
	Val2, ok2 := valMap["windows"]
	if !ok1 || !ok2 {
		return nil
	}
	valE1, ok := Val1.(map[interface{}]interface{})
	if p := envParser(valE1); ok && p != nil {
		return p
	}
	valE2, ok := Val2.(map[interface{}]interface{})
	if p := envParser(valE2); ok && p != nil {
		return p
	}
	return &BaseParser{&OsParser{Linux: Val1, Windows: Val2}}
}
