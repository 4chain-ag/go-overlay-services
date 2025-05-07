package app

// ErrorType represents a generic category of error used as descriptor
// to clarify the nature of a failure that occurred in application-layer dependencies.
type ErrorType struct {
	s string
}

// IsZero returns true if the error is in its zero-value state.
func (e ErrorType) IsZero() bool { return e == ErrorType{} }

// String returns the internal error type string.
func (e ErrorType) String() string { return e.s }

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

// Error defines a generic application-layer error that should be translated
// into the specific response format returned to the requester.
// An error includes the source error message and a type describing the particular
// category of the failure. The type should be used during the translation process
// in the error-handling implementation.
// The source error message may contain internal details, so it is not recommended
// to include it in the final response due to the potential risk of leaking
// sensitive data to the requester.
type Error struct {
	err       string
	service   string
	errorType ErrorType
}

func (e Error) Service() string { return e.service }

// IsZero returns true if the error is in its zero-value state.
func (e Error) IsZero() bool { return e == Error{} }

// Error returns the source error message, which may contain internal details.
func (e Error) Error() string { return e.err }

// ErrorType returns the category of the error, which should be used
// during error handling and response format translation.
func (e Error) ErrorType() ErrorType { return e.errorType }

// NewIncorrectInputError creates a new Error with ErrorTypeIncorrectInput.
func NewIncorrectInputError(service, err string) Error {
	return Error{
		err:       err,
		service:   service,
		errorType: ErrorTypeIncorrectInput,
	}
}

// NewProviderFailureError creates a new Error with ErrorTypeProviderFailure.
func NewProviderFailureError(service, err string) Error {
	return Error{
		err:       err,
		service:   service,
		errorType: ErrorTypeProviderFailure,
	}
}

// NewOperationTimeoutError creates a new Error with ErrorTypeOperationTimeout.
func NewOperationTimeoutError(service, err string) Error {
	return Error{
		err:       err,
		service:   service,
		errorType: ErrorTypeOperationTimeout,
	}
}
