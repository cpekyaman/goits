package commons

import (
	"strings"
)

type Layer uint8

var layers = [...]string{"api", "svc", "rep", "dom"}

const (
	LayerApi Layer = iota
	LayerService
	LayerRepository
	LayerDomain
)

func (this Layer) String() string {
	return layers[this]
}

type Subsystem uint8

var subsystems = [...]string{"vld", "bind", "db", "scr"}

const (
	SysValidation Subsystem = iota
	SysBinding
	SysDb
	SysSecurity
)

func (this Subsystem) String() string {
	return subsystems[this]
}

type ErrorType uint8

var errorTypes = [...]string{"validation", "client", "notfound", "db", "internal"}

const (
	ErrValidation ErrorType = iota
	ErrClient
	ErrNotFound
	ErrDb
	ErrInternal
)

func (this ErrorType) String() string {
	return errorTypes[this]
}

// AppError
type AppError struct {
	ErrorType ErrorType `json:"errorType"`
	Message   string    `json:"message"`
	Cause     string    `json:"cause"`
}

func DetermineErrorType(err error) ErrorType {
	msg := err.Error()

	if strings.HasPrefix(msg, "validation:") {
		return ErrValidation
	} else if strings.HasPrefix(msg, "sql:") || strings.HasPrefix(msg, "driver:") {
		if strings.Contains(msg, "no rows") {
			return ErrNotFound
		} else {
			return ErrDb
		}
	} else if strings.HasPrefix(msg, "binding:") {
		return ErrClient
	} else {
		return ErrInternal
	}
}
