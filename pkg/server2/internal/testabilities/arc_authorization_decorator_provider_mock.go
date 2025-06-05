package testabilities

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

// HandlerMockExpectations defines the expected behavior of the HandlerMock during a test.
type HandlerMockExpectations struct {
	Error      error
	HandleCall bool
}

// DefaultHandlerMockExpectations provides default expectations for HandlerMock,
// including no error and expecting the handle call.
var DefaultHandlerMockExpectations = HandlerMockExpectations{
	Error:      nil,
	HandleCall: true,
}

// HandlerMock is a mock implementation of the Handler interface,
// used for testing the behavior of components that depend on handlers.
type HandlerMock struct {
	t            *testing.T
	expectations HandlerMockExpectations
	called       bool
	calledCtx    *fiber.Ctx
}

// Handle simulates the handling of a request. It records the call and returns
// the predefined error if set.
func (h *HandlerMock) Handle(c *fiber.Ctx) error {
	h.t.Helper()
	h.called = true
	h.calledCtx = c
	return h.expectations.Error
}

// AssertCalled verifies that the Handle method was called if it was expected to be.
func (h *HandlerMock) AssertCalled() {
	h.t.Helper()
	require.Equal(h.t, h.expectations.HandleCall, h.called, "Discrepancy between expected and actual Handle call")
}

// CalledCtx returns the fiber.Ctx that was passed to the Handle method.
func (h *HandlerMock) CalledCtx() *fiber.Ctx {
	h.t.Helper()
	return h.calledCtx
}

// NewHandlerMock creates a new instance of HandlerMock with the given expectations.
func NewHandlerMock(t *testing.T, expectations HandlerMockExpectations) *HandlerMock {
	mock := &HandlerMock{
		t:            t,
		expectations: expectations,
	}
	return mock
}
