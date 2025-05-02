package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/stretchr/testify/require"
)

type mockLookupServiceDocumentationProvider struct {
	documentation string
	err           error
}

func (m *mockLookupServiceDocumentationProvider) GetDocumentationForLookupServiceProvider(lookupService string) (string, error) {
	return m.documentation, m.err
}

func TestLookupServiceDocumentationService_GetDocumentation_Success(t *testing.T) {
	// Given
	expectedDocumentation := "# Test Documentation\nThis is a test markdown document."
	mockProvider := &mockLookupServiceDocumentationProvider{
		documentation: expectedDocumentation,
		err:           nil,
	}
	service := app.NewLookupServiceDocumentationService(mockProvider)

	// When
	documentation, err := service.GetDocumentation(context.Background(), "testProvider")

	// Then
	require.NoError(t, err)
	require.Equal(t, expectedDocumentation, documentation)
}

func TestLookupServiceDocumentationService_GetDocumentation_EmptyLookupServiceName(t *testing.T) {
	// Given
	mockProvider := &mockLookupServiceDocumentationProvider{
		documentation: "# Test Documentation\nThis is a test markdown document.",
		err:           nil,
	}
	service := app.NewLookupServiceDocumentationService(mockProvider)

	// When
	_, err := service.GetDocumentation(context.Background(), "")

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, app.ErrEmptyLookupServiceName)
}

func TestLookupServiceDocumentationService_GetDocumentation_ProviderError(t *testing.T) {
	// Given
	mockError := errors.New("provider error")
	mockProvider := &mockLookupServiceDocumentationProvider{
		documentation: "",
		err:           mockError,
	}
	service := app.NewLookupServiceDocumentationService(mockProvider)

	// When
	_, err := service.GetDocumentation(context.Background(), "testProvider")

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, app.ErrLookupServiceProviderNotFound)
}

func TestNewLookupServiceDocumentationService_NilProvider(t *testing.T) {
	// Given, When, Then
	require.Panics(t, func() {
		app.NewLookupServiceDocumentationService(nil)
	})
}
