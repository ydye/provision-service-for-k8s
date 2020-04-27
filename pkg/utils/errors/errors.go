package errors

import (
	"fmt"
)

// ProvisionErrorType describes a high-level category of a given error
type ProvisionErrorType string

// ProvisionError contains information about provision errors
type ProvisionError interface {
	// Error implements golang error interface
	Error() string

	// Type returns the type of ProvisionError
	Type() ProvisionErrorType

	// AddPrefix adds a prefix to error message.
	// Example:
	// If err := DoSomething(myObject); err != nil {
	//   return err.AddPrefix("can't do somthine with %v: ", myObject)
	// }
	AddPrefix(msg string, args ...interface{}) ProvisionError
}

type provisionErrorImpl struct {
	errorType ProvisionErrorType
	msg string
}

const (
	InternalError ProvisionErrorType = "internalError"
	ConfigurationError ProvisionErrorType = "configurationError"
	ApiCallError ProvisionErrorType = "apiCallError"
)

func NewProvisionError(errorType ProvisionErrorType, msg string, args ...interface{}) ProvisionError {
	return provisionErrorImpl{
		errorType: errorType,
		msg:       msg,
	}
}

func ToProvisionError(defaultType ProvisionErrorType, err error) ProvisionError {
	if e, ok := err.(ProvisionError); ok {
		return e
	}
	return NewProvisionError(defaultType, "%v", err)
}

func (e provisionErrorImpl) Error() string {
	return e.msg
}

func (e provisionErrorImpl) Type() ProvisionErrorType {
	return e.errorType
}

func (e provisionErrorImpl) AddPrefix(msg string, args ...interface{}) ProvisionError {
	e.msg = fmt.Sprintf(msg, args...) + e.msg
	return e
}



