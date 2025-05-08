package dw

import (
	"github.com/pkg/errors"
	"os"
	"reflect"
)

type Config struct {
	Path      string
	Marshal   func(v any) ([]byte, error)
	Unmarshal func(data []byte, v any) error
}

type DataWriter[T any] struct {
	cfg       Config
	Data      T
	DataBytes []byte
}

func New[T any](cfg Config) (*DataWriter[T], error) {
	manager := &DataWriter[T]{
		cfg: cfg,
	}
	val := reflect.ValueOf(manager.Data)
	if val.Kind() == reflect.Map {
		manager.Data = reflect.MakeMap(val.Type()).Interface().(T)
	}
	err := manager.Read()
	if err != nil {
		return nil, err
	}
	return manager, nil
}

func (manager *DataWriter[T]) Read() error {
	fileBytes, err := os.ReadFile(manager.cfg.Path)
	if err != nil {
		return errors.WithStack(err)
	}
	err = manager.cfg.Unmarshal(fileBytes, &manager.Data)
	if err != nil {
		return errors.WithStack(err)
	}
	manager.DataBytes = fileBytes
	return nil
}

func (manager *DataWriter[T]) Write() error {
	marshalBytes, err := manager.cfg.Marshal(manager.Data)
	if err != nil {
		return errors.WithStack(err)
	}
	err = os.WriteFile(manager.cfg.Path, marshalBytes, 0666)
	if err != nil {
		return errors.WithStack(err)
	}
	manager.DataBytes = marshalBytes
	return nil
}
