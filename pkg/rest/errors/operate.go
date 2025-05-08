package errors

import (
	"fmt"
	"strings"
)

type ExistError struct {
	Type string
	Id   string
}

func (error ExistError) Error() string {
	return fmt.Sprintf("%s %s exist", error.Type, error.Id)
}

type NotExistError struct {
	Type string
	Id   string
}

func (error NotExistError) Error() string {
	return fmt.Sprintf("%s %s doesn't exist", error.Type, error.Id)
}

type UpdateErr struct {
	Type   string
	Id     string
	Fields []string
}

func (err *UpdateErr) Error() string {
	if err.Fields == nil {
		return ""
	}
	return fmt.Sprintf("update %s %s failed, including %s", err.Type, err.Id, strings.Join(err.Fields, ", "))
}

func (err *UpdateErr) Add(field string) {
	if err.Fields == nil {
		err.Fields = make([]string, 1)
		err.Fields[0] = field
	} else {
		err.Fields = append(err.Fields, field)
	}
}

type DeleteErr struct {
	Type   string
	Id     string
	Result string
}

func (err DeleteErr) Error() string {
	return fmt.Sprintf("remove %s %s failed, result is: %s", err.Type, err.Id, err.Result)
}
