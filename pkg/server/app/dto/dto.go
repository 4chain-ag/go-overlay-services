package dto

type HandlerResponse struct {
	Message string
}

var HandlerResponseOK = HandlerResponse{Message: "OK :-)"}
var HandlerResponseNonOK = HandlerResponse{Message: "Non OK :-("}
