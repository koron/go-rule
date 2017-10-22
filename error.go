package rule

import (
	"fmt"
	"reflect"
)

// NativeFuncError represents runtime error in native function.
type NativeFuncError struct {
	Name string
	Err  error
}

func (err *NativeFuncError) Error() string {
	return fmt.Sprintf("native func %s() is failed: %s", err.Name, err.Err.Error())
}

// ArgsCountError represents mismatch of arguments.
type ArgsCountError struct {
	Expected int
	Given    int
}

func (err *ArgsCountError) Error() string {
	return fmt.Sprintf("args count mismatch, expected %d but given %d",
		err.Expected, err.Given)
}

// ArgTypeError reprents mismatch type of an argument.
type ArgTypeError struct {
	Pos      int
	Expected string
	Given    interface{} // Given value
}

func (err *ArgTypeError) Error() string {
	return fmt.Sprintf("type of arg#%d mismatch, expected %q but given %q",
		err.Pos, err.Expected, reflect.TypeOf(err.Given).String())
}
