package config

import "path/filepath"

var Global *Config

func InitFromFile(fileName string) error {
	return InitFromFileCostume(fileName, YAML)
}

func InitFromFileCostume(fileName string, parseType ParseType) error {
	absPath, err := filepath.Abs(fileName)
	if err != nil {
		return err
	}

	cfg, err := ParseFile(absPath, parseType)
	if err != nil {
		return err
	}
	Global = cfg
	return nil
}

func Init(data []byte) error {
	return InitCostume(data, YAML)
}

func InitCostume(data []byte, parseType ParseType) error {
	cfg, err := Parse(data, parseType)
	if err != nil {
		return err
	}
	Global = cfg
	return nil
}
