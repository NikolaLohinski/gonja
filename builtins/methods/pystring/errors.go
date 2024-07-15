package pystring

import "fmt"

var (
	ErrValue     = fmt.Errorf("ValueError")
	ErrIndex     = fmt.Errorf("IndexError")
	ErrInternal  = fmt.Errorf("InternalError")
	ErrKey       = fmt.Errorf("KeyError")
	ErrArguments = fmt.Errorf("ArgumentError")
)
