package dto

// HandlerResponse is a custom test response returned by any handler under test
// and development. This is only a temporary solution and will be removed later.
type HandlerResponse struct {
	Message string
}

// HandlerResponseOK describes a valid response type returned by the HTTP Overlay API.
var HandlerResponseOK = HandlerResponse{Message: "OK :-)"}

// HandlerResponseNonOK describes invalid response type returned by the HTTP Overlay API.
var HandlerResponseNonOK = HandlerResponse{Message: "Non OK :-("}
