// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package openapi

import (
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/oapi-codegen/runtime"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Error defines model for Error.
type Error struct {
	// Message Human-readable error message
	Message string `json:"message"`
}

// BadRequestResponse defines model for BadRequestResponse.
type BadRequestResponse = Error

// InternalServerErrorResponse defines model for InternalServerErrorResponse.
type InternalServerErrorResponse = Error

// NotFoundResponse defines model for NotFoundResponse.
type NotFoundResponse = Error

// RequestTimeoutResponse defines model for RequestTimeoutResponse.
type RequestTimeoutResponse = Error

// ArcIngestJSONBody defines parameters for ArcIngest.
type ArcIngestJSONBody struct {
	// BlockHeight Block height where the transaction was included
	BlockHeight uint32 `json:"blockHeight"`

	// MerklePath Merkle path in hexadecimal format
	MerklePath string `json:"merklePath"`

	// Txid Transaction ID in hexadecimal format
	Txid string `json:"txid"`
}

// GetLookupServiceProviderDocumentationParams defines parameters for GetLookupServiceProviderDocumentation.
type GetLookupServiceProviderDocumentationParams struct {
	// LookupService The name of the lookup service provider to retrieve documentation for
	LookupService string `form:"lookupService" json:"lookupService"`
}

// GetTopicManagerDocumentationParams defines parameters for GetTopicManagerDocumentation.
type GetTopicManagerDocumentationParams struct {
	// TopicManager The name of the topic manager to retrieve documentation for
	TopicManager string `form:"topicManager" json:"topicManager"`
}

// LookupQuestionJSONBody defines parameters for LookupQuestion.
type LookupQuestionJSONBody struct {
	// Query Query parameters specific to the service
	Query map[string]interface{} `json:"query"`

	// Service Service name to query
	Service string `json:"service"`
}

// RequestForeignGASPNodeJSONBody defines parameters for RequestForeignGASPNode.
type RequestForeignGASPNodeJSONBody struct {
	// GraphID The graph ID in the format of "txID.outputIndex"
	GraphID string `json:"graphID"`

	// OutputIndex The output index
	OutputIndex uint32 `json:"outputIndex"`

	// TxID The transaction ID
	TxID string `json:"txID"`
}

// RequestForeignGASPNodeParams defines parameters for RequestForeignGASPNode.
type RequestForeignGASPNodeParams struct {
	XBSVTopic string `json:"X-BSV-Topic"`
}

// RequestSyncResponseJSONBody defines parameters for RequestSyncResponse.
type RequestSyncResponseJSONBody struct {
	// Since Timestamp or sequence number from which to start synchronization
	Since uint32 `json:"since"`

	// Version The version number of the GASP protocol
	Version int `json:"version"`
}

// RequestSyncResponseParams defines parameters for RequestSyncResponse.
type RequestSyncResponseParams struct {
	// XBSVTopic Topic identifier for the sync response request
	XBSVTopic string `json:"X-BSV-Topic"`
}

// SubmitTransactionParams defines parameters for SubmitTransaction.
type SubmitTransactionParams struct {
	XTopics []string `json:"x-topics"`
}

// ArcIngestJSONRequestBody defines body for ArcIngest for application/json ContentType.
type ArcIngestJSONRequestBody ArcIngestJSONBody

// LookupQuestionJSONRequestBody defines body for LookupQuestion for application/json ContentType.
type LookupQuestionJSONRequestBody LookupQuestionJSONBody

// RequestForeignGASPNodeJSONRequestBody defines body for RequestForeignGASPNode for application/json ContentType.
type RequestForeignGASPNodeJSONRequestBody RequestForeignGASPNodeJSONBody

// RequestSyncResponseJSONRequestBody defines body for RequestSyncResponse for application/json ContentType.
type RequestSyncResponseJSONRequestBody RequestSyncResponseJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /api/v1/admin/startGASPSync)
	StartGASPSync(c *fiber.Ctx) error

	// (POST /api/v1/admin/syncAdvertisements)
	AdvertisementsSync(c *fiber.Ctx) error

	// (POST /api/v1/arc-ingest)
	ArcIngest(c *fiber.Ctx) error

	// (GET /api/v1/getDocumentationForLookupServiceProvider)
	GetLookupServiceProviderDocumentation(c *fiber.Ctx, params GetLookupServiceProviderDocumentationParams) error

	// (GET /api/v1/getDocumentationForTopicManager)
	GetTopicManagerDocumentation(c *fiber.Ctx, params GetTopicManagerDocumentationParams) error

	// (GET /api/v1/listLookupServiceProviders)
	ListLookupServiceProviders(c *fiber.Ctx) error

	// (GET /api/v1/listTopicManagers)
	ListTopicManagers(c *fiber.Ctx) error

	// (POST /api/v1/lookup)
	LookupQuestion(c *fiber.Ctx) error

	// (POST /api/v1/requestForeignGASPNode)
	RequestForeignGASPNode(c *fiber.Ctx, params RequestForeignGASPNodeParams) error

	// (POST /api/v1/requestSyncResponse)
	RequestSyncResponse(c *fiber.Ctx, params RequestSyncResponseParams) error

	// (POST /api/v1/submit)
	SubmitTransaction(c *fiber.Ctx, params SubmitTransactionParams) error
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	handler           ServerInterface
	globalMiddleware  []fiber.Handler
	handlerMiddleware []fiber.Handler
}

// StartGASPSync operation middleware
func (siw *ServerInterfaceWrapper) StartGASPSync(c *fiber.Ctx) error {

	c.Context().SetUserValue(BearerAuthScopes, []string{"admin"})

	for _, m := range siw.handlerMiddleware {
		if err := m(c); err != nil {
			return err
		}
	}
	return siw.handler.StartGASPSync(c)
}

// AdvertisementsSync operation middleware
func (siw *ServerInterfaceWrapper) AdvertisementsSync(c *fiber.Ctx) error {

	c.Context().SetUserValue(BearerAuthScopes, []string{"admin"})

	for _, m := range siw.handlerMiddleware {
		if err := m(c); err != nil {
			return err
		}
	}
	return siw.handler.AdvertisementsSync(c)
}

// ArcIngest operation middleware
func (siw *ServerInterfaceWrapper) ArcIngest(c *fiber.Ctx) error {

	c.Context().SetUserValue(BearerAuthScopes, []string{"user"})

	for _, m := range siw.handlerMiddleware {
		if err := m(c); err != nil {
			return err
		}
	}
	return siw.handler.ArcIngest(c)
}

// GetLookupServiceProviderDocumentation operation middleware
func (siw *ServerInterfaceWrapper) GetLookupServiceProviderDocumentation(c *fiber.Ctx) error {

	var err error

	c.Context().SetUserValue(BearerAuthScopes, []string{"user"})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetLookupServiceProviderDocumentationParams

	var query url.Values
	query, err = url.ParseQuery(string(c.Request().URI().QueryString()))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid format for query string")
	}

	// ------------- Required query parameter "lookupService" -------------

	if paramValue := c.Query("lookupService"); paramValue != "" {

	} else {
		return fiber.NewError(fiber.StatusBadRequest, "A valid lookupService must be provided to retrieve documentation.")
	}

	err = runtime.BindQueryParameter("form", true, true, "lookupService", query, &params.LookupService)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid format for parameter lookupService")
	}

	for _, m := range siw.handlerMiddleware {
		if err := m(c); err != nil {
			return err
		}
	}
	return siw.handler.GetLookupServiceProviderDocumentation(c, params)
}

// GetTopicManagerDocumentation operation middleware
func (siw *ServerInterfaceWrapper) GetTopicManagerDocumentation(c *fiber.Ctx) error {

	var err error

	c.Context().SetUserValue(BearerAuthScopes, []string{"user"})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetTopicManagerDocumentationParams

	var query url.Values
	query, err = url.ParseQuery(string(c.Request().URI().QueryString()))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid format for query string")
	}

	// ------------- Required query parameter "topicManager" -------------

	if paramValue := c.Query("topicManager"); paramValue != "" {

	} else {
		return fiber.NewError(fiber.StatusBadRequest, "A valid topicManager must be provided to retrieve documentation.")
	}

	err = runtime.BindQueryParameter("form", true, true, "topicManager", query, &params.TopicManager)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid format for parameter topicManager")
	}

	for _, m := range siw.handlerMiddleware {
		if err := m(c); err != nil {
			return err
		}
	}
	return siw.handler.GetTopicManagerDocumentation(c, params)
}

// ListLookupServiceProviders operation middleware
func (siw *ServerInterfaceWrapper) ListLookupServiceProviders(c *fiber.Ctx) error {

	c.Context().SetUserValue(BearerAuthScopes, []string{"user"})

	for _, m := range siw.handlerMiddleware {
		if err := m(c); err != nil {
			return err
		}
	}
	return siw.handler.ListLookupServiceProviders(c)
}

// ListTopicManagers operation middleware
func (siw *ServerInterfaceWrapper) ListTopicManagers(c *fiber.Ctx) error {

	c.Context().SetUserValue(BearerAuthScopes, []string{"user"})

	for _, m := range siw.handlerMiddleware {
		if err := m(c); err != nil {
			return err
		}
	}
	return siw.handler.ListTopicManagers(c)
}

// LookupQuestion operation middleware
func (siw *ServerInterfaceWrapper) LookupQuestion(c *fiber.Ctx) error {

	c.Context().SetUserValue(BearerAuthScopes, []string{"user"})

	for _, m := range siw.handlerMiddleware {
		if err := m(c); err != nil {
			return err
		}
	}
	return siw.handler.LookupQuestion(c)
}

// RequestForeignGASPNode operation middleware
func (siw *ServerInterfaceWrapper) RequestForeignGASPNode(c *fiber.Ctx) error {

	var err error

	c.Context().SetUserValue(BearerAuthScopes, []string{"user"})

	// Parameter object where we will unmarshal all parameters from the context
	var params RequestForeignGASPNodeParams

	headers := c.GetReqHeaders()

	// ------------- Required header parameter "X-BSV-Topic" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-BSV-Topic")]; found {
		var XBSVTopic string

		err = runtime.BindStyledParameterWithOptions("simple", "X-BSV-Topic", valueList[0], &XBSVTopic, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: true})
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "One or more topics are in an invalid format. Empty string values are not allowed.")
		}

		params.XBSVTopic = XBSVTopic

	} else {
		return fiber.NewError(fiber.StatusBadRequest, "The submitted request does not include required header: X-BSV-Topic.")
	}

	for _, m := range siw.handlerMiddleware {
		if err := m(c); err != nil {
			return err
		}
	}
	return siw.handler.RequestForeignGASPNode(c, params)
}

// RequestSyncResponse operation middleware
func (siw *ServerInterfaceWrapper) RequestSyncResponse(c *fiber.Ctx) error {

	var err error

	c.Context().SetUserValue(BearerAuthScopes, []string{"user"})

	// Parameter object where we will unmarshal all parameters from the context
	var params RequestSyncResponseParams

	headers := c.GetReqHeaders()

	// ------------- Required header parameter "X-BSV-Topic" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-BSV-Topic")]; found {
		var XBSVTopic string

		err = runtime.BindStyledParameterWithOptions("simple", "X-BSV-Topic", valueList[0], &XBSVTopic, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: true})
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "One or more topics are in an invalid format. Empty string values are not allowed.")
		}

		params.XBSVTopic = XBSVTopic

	} else {
		return fiber.NewError(fiber.StatusBadRequest, "The submitted request does not include required header: X-BSV-Topic.")
	}

	for _, m := range siw.handlerMiddleware {
		if err := m(c); err != nil {
			return err
		}
	}
	return siw.handler.RequestSyncResponse(c, params)
}

// SubmitTransaction operation middleware
func (siw *ServerInterfaceWrapper) SubmitTransaction(c *fiber.Ctx) error {

	var err error

	c.Context().SetUserValue(BearerAuthScopes, []string{"user"})

	// Parameter object where we will unmarshal all parameters from the context
	var params SubmitTransactionParams

	headers := c.GetReqHeaders()

	// ------------- Required header parameter "x-topics" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("x-topics")]; found {
		var XTopics []string

		err = runtime.BindStyledParameterWithOptions("simple", "x-topics", valueList[0], &XTopics, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: true, Required: true})
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "One or more topics are in an invalid format. Empty string values are not allowed.")
		}

		params.XTopics = XTopics

	} else {
		return fiber.NewError(fiber.StatusBadRequest, "The submitted request does not include required header: x-topics.")
	}

	for _, m := range siw.handlerMiddleware {
		if err := m(c); err != nil {
			return err
		}
	}
	return siw.handler.SubmitTransaction(c, params)
}

// FiberServerOptions provides options for the Fiber server.
type FiberServerOptions struct {
	BaseURL           string
	GlobalMiddleware  []fiber.Handler
	HandlerMiddleware []fiber.Handler
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router fiber.Router, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, FiberServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router fiber.Router, si ServerInterface, options FiberServerOptions) {
	wrapper := ServerInterfaceWrapper{
		handler:           si,
		globalMiddleware:  options.GlobalMiddleware,
		handlerMiddleware: options.HandlerMiddleware,
	}

	for _, m := range options.GlobalMiddleware {
		router.Use(m)
	}

	router.Post(options.BaseURL+"/api/v1/admin/startGASPSync", wrapper.StartGASPSync)

	router.Post(options.BaseURL+"/api/v1/admin/syncAdvertisements", wrapper.AdvertisementsSync)

	router.Post(options.BaseURL+"/api/v1/arc-ingest", wrapper.ArcIngest)

	router.Get(options.BaseURL+"/api/v1/getDocumentationForLookupServiceProvider", wrapper.GetLookupServiceProviderDocumentation)

	router.Get(options.BaseURL+"/api/v1/getDocumentationForTopicManager", wrapper.GetTopicManagerDocumentation)

	router.Get(options.BaseURL+"/api/v1/listLookupServiceProviders", wrapper.ListLookupServiceProviders)

	router.Get(options.BaseURL+"/api/v1/listTopicManagers", wrapper.ListTopicManagers)

	router.Post(options.BaseURL+"/api/v1/lookup", wrapper.LookupQuestion)

	router.Post(options.BaseURL+"/api/v1/requestForeignGASPNode", wrapper.RequestForeignGASPNode)

	router.Post(options.BaseURL+"/api/v1/requestSyncResponse", wrapper.RequestSyncResponse)

	router.Post(options.BaseURL+"/api/v1/submit", wrapper.SubmitTransaction)

}
