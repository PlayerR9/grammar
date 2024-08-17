package assert

import (
	"fmt"
	"strconv"
	"strings"
)

// Assert panics iff the condition is false. The panic is not a string
// but an error of type *ErrAssertionFailed.
//
// Parameters:
//   - cond: the condition to check.
//   - msg: the message to show if the condition is not met.
//
// Example:
//
//	foo := "foo"
//	Assert(foo == "bar", "foo is not bar") // panics: "assertion failed: foo is not bar"
func Assert(cond bool, msg string) {
	if cond {
		return
	}

	panic(NewErrAssertionFailed(msg))
}

// AssertF same as Assert but with a format string and arguments that are in
// accordance with fmt.Printf.
//
// Parameters:
//   - cond: the condition to check.
//   - format: the format string to show if the condition is not met.
//   - args: the arguments to pass to the format string.
//
// Example:
//
//	foo := "foo"
//	bar := "bar"
//	AssertF(foo == bar, "%s is not %s", foo, bar) // panics: "assertion failed: foo is not bar"
func AssertF(cond bool, format string, args ...any) {
	if cond {
		return
	}

	panic(NewErrAssertionFailed(fmt.Sprintf(format, args...)).Error())
}

// AssertErr is the same as Assert but for errors. Best used for ensuring that a function
// does not return an unexpected error.
//
// Parameters:
//   - err: the error to check.
//   - format: the format describing the function's signature.
//   - args: the arguments passed to the function.
//
// Example:
//
//	foo := "foo"
//	err := my_function(foo, "bar")
//	AssertErr(err, "my_function(%s, %s)", foo, "bar")
//	// panics: "assertion failed: function my_function(foo, bar) returned the error: <err>"
func AssertErr(err error, format string, args ...any) {
	if err == nil {
		return
	}

	var builder strings.Builder

	builder.WriteString("function ")
	fmt.Fprintf(&builder, format, args...)
	builder.WriteString(" returned the error: ")
	builder.WriteString(err.Error())

	panic(NewErrAssertionFailed(builder.String()).Error())
}

// AssertOk is the same as Assert but for booleans. Best used for ensuring that a function that
// are supposed to return the boolean `true` does not return `false`.
//
// Parameters:
//   - cond: the result of the function.
//   - format: the format describing the function's signature.
//   - args: the arguments passed to the function.
//
// Example:
//
//	foo := "foo"
//	ok := my_function(foo, "bar")
//	AssertOk(ok, "my_function(%s, %s)", foo, "bar")
//	// panics: "assertion failed: function my_function(foo, bar) returned false while true was expected"
func AssertOk(cond bool, format string, args ...any) {
	if cond {
		return
	}

	var builder strings.Builder

	builder.WriteString("function ")
	fmt.Fprintf(&builder, format, args...)
	builder.WriteString(" returned false while true was expected")

	panic(NewErrAssertionFailed(builder.String()).Error())
}

// AssertDeref tries to dereference an element and panics if it is nil.
//
// Parameters:
//   - elem: the element to dereference.
//   - param_name: the name of the parameter.
//
// Returns:
//   - T: the dereferenced element.
func AssertDeref[T any](elem *T, param_name string) T {
	if elem != nil {
		return *elem
	}

	var builder strings.Builder

	builder.WriteString("Parameter (")
	builder.WriteString(strconv.Quote(param_name))
	builder.WriteString(") must not be nil")

	panic(NewErrAssertionFailed(builder.String()).Error())
}

// AssertNotNil panics if the element is nil.
//
// Parameters:
//   - elem: the element to check.
//   - param_name: the name of the parameter.
func AssertNotNil(elem any, param_name string) {
	if elem != nil {
		return
	}

	var builder strings.Builder

	builder.WriteString("Parameter (")
	builder.WriteString(strconv.Quote(param_name))
	builder.WriteString(") must not be nil")

	panic(NewErrAssertionFailed(builder.String()).Error())
}

// AssertTypeOf panics if the element is not of the expected type.
//
// Parameters:
//   - elem: the element to check.
//   - var_name: the name of the variable.
//   - allow_nil: if the element can be nil.
func AssertTypeOf[T any](elem any, var_name string, allow_nil bool) {
	if elem == nil {
		if !allow_nil {
			var builder strings.Builder

			builder.WriteString("expected ")
			builder.WriteString(strconv.Quote(var_name))
			builder.WriteString(" to be of type ")
			builder.WriteString(fmt.Sprintf("%T", *new(T)))
			builder.WriteString(", got nil instead")

			panic(NewErrAssertionFailed(builder.String()).Error())
		}

		return
	}

	_, ok := elem.(T)
	if !ok {
		var builder strings.Builder

		builder.WriteString("expected ")
		builder.WriteString(strconv.Quote(var_name))
		builder.WriteString(" to be of type ")
		builder.WriteString(fmt.Sprintf("%T", *new(T)))
		builder.WriteString(", got ")
		builder.WriteString(fmt.Sprintf("%T", elem))
		builder.WriteString(" instead")

		panic(NewErrAssertionFailed(builder.String()).Error())
	}
}

// AssertConv tries to convert an element to the expected type and panics if it is not possible.
//
// Parameters:
//   - elem: the element to check.
//   - var_name: the name of the variable.
//
// Returns:
//   - T: the converted element.
func AssertConv[T any](elem any, var_name string) T {
	if elem == nil {
		var builder strings.Builder

		builder.WriteString("expected ")
		builder.WriteString(strconv.Quote(var_name))
		builder.WriteString(" to be of type ")
		builder.WriteString(fmt.Sprintf("%T", *new(T)))
		builder.WriteString(", got nil instead")

		panic(NewErrAssertionFailed(builder.String()).Error())
	}

	res, ok := elem.(T)
	if !ok {
		var builder strings.Builder

		builder.WriteString("expected ")
		builder.WriteString(strconv.Quote(var_name))
		builder.WriteString(" to be of type ")
		builder.WriteString(fmt.Sprintf("%T", *new(T)))
		builder.WriteString(", got ")
		builder.WriteString(fmt.Sprintf("%T", elem))
		builder.WriteString(" instead")

		panic(NewErrAssertionFailed(builder.String()).Error())
	}

	return res
}
