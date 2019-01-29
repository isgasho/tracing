package system

// System ...
type System interface {
	Start() error
	Close() error
}
