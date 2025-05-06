package app

// ErrorType represents a category of error as a string.
type ErrorType struct {
	s string
}

// Predefined error types for categorizing different kinds of errors.
var (
	// ErrorTypeProviderFailure indicates an error caused by a provider/service dependency.
	ErrorTypeProviderFailure = ErrorType{"provider-failure"}

	// ErrorTypeAuthorization indicates an error due to unauthorized access or permissions.
	ErrorTypeAuthorization = ErrorType{"authorization"}

	// ErrorTypeIncorrectInput indicates an error due to malformed or invalid user input.
	ErrorTypeIncorrectInput = ErrorType{"incorrect-input"}

	// ErrorTypeUnknown is used when the error cause is not clearly identifiable.
	ErrorTypeUnknown = ErrorType{"unknown"}

	// ErrorTypeOperationTimeout indicates the operation timed out.
	ErrorTypeOperationTimeout = ErrorType{"operation-timeout"}
)

// Error represents a custom application error with type and slug information.
type Error struct {
	err       string
	errorType ErrorType
}

// Error implements the standard error interface.
func (e Error) Error() string { return e.err }

// ErrorType returns the type of the error.
func (e Error) ErrorType() ErrorType { return e.errorType }

// NewIncorrectInputError creates a new Error with ErrorTypeIncorrectInput.
func NewIncorrectInputError(err string) Error {
	return Error{
		err:       err,
		errorType: ErrorTypeIncorrectInput,
	}
}

// NewProviderFailureError creates a new Error with ErrorTypeProviderFailure.
func NewProviderFailureError(err string) Error {
	return Error{
		err:       err,
		errorType: ErrorTypeProviderFailure,
	}
}

// NewOperationTimeoutError creates a new Error with ErrorTypeOperationTimeout.
func NewOperationTimeoutError(err string) Error {
	return Error{
		err:       err,
		errorType: ErrorTypeOperationTimeout,
	}
}
