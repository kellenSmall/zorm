package errs

import (
	"errors"
	"fmt"
)

var (
	ErrInputNil             = errors.New("orm: 不支持 nil")
	ErrPointerOnly          = errors.New("orm: 只支持一级指针作为输入，例如 *User")
	ErrNoRows               = errors.New("orm: 未找到数据")
	ErrTooManyReturnColumns = errors.New("orm: 过多列")
)

func NewErrInvalidTagContent(tag string) error {
	return fmt.Errorf("orm: 错误的标签设置: %s", tag)
}

// NewErrUnsupportedSelectable 返回一个不支持该 selectable 的错误信息
// 即 exp 不能作为 SELECT xxx 的一部分
func NewErrUnsupportedSelectable(exp any) error {
	return fmt.Errorf("orm: 不支持的目标列 %v", exp)
}

// NewErrUnknownColumn 返回代表未知列的错误
// 一般意味着你使用了错误的列名
// 注意和 NewErrUnknownField 区别
func NewErrUnknownColumn(col string) error {
	return fmt.Errorf("orm: 未知列 %s", col)
}

func NewErrUnsupportedAssignableType(exp any) error {
	return fmt.Errorf("orm: 不支持的 Assignable 表达式 %v", exp)
}

// NewErrUnsupportedExpressionType 返回一个不支持该 expression 错误信息
func NewErrUnsupportedExpressionType(exp any) error {
	return fmt.Errorf("orm: 不支持的表达式 %v", exp)
}
