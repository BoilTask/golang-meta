package subsystem

import "meta/object"

type Interface interface {
	object.Interface
	Init() error
	Start() error
	Stop() error
}
