package errs

import "errors"

var (
	ErrInputNil    = errors.New("orm: 不支持 nil")
	ErrPointerOnly = errors.New("orm: 只支持一级指针作为输入，例如 *User")
)
