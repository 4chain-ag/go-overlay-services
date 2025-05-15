package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/stretchr/testify/require"
)

func TestGetLookupServiceProviderDocumentation_Success(t *testing.T) {
	// Given
	service := app.NewLookupServiceProviderDocumentationService(&testabilities.MockLookupServiceProviderDocumentationProvider{ShouldFail: false})

	// When
	documentation, err := service.GetDocumentation(context.Background(), "test-service")

	// Then
	require.NoError(t, err)
	require.Equal(t, "# Test Documentation\nThis is a test markdown document.", documentation)
}

func TestGetLookupServiceProviderDocumentation_EmptyLookupServiceName(t *testing.T) {
	// Given
	service := app.NewLookupServiceProviderDocumentationService(&testabilities.MockLookupServiceProviderDocumentationProvider{ShouldFail: false})

	// When
	documentation, err := service.GetDocumentation(context.Background(), "")

	// Then
	require.Error(t, err)
	require.Empty(t, documentation)
	var target app.Error
	require.True(t, errors.As(err, &target))
	require.Equal(t, app.ErrorTypeIncorrectInput, target.ErrorType())
	require.Equal(t, "lookup service name cannot be empty", target.Error())
}

func TestGetLookupServiceProviderDocumentation_ProviderError(t *testing.T) {
	// Given
	service := app.NewLookupServiceProviderDocumentationService(&testabilities.MockLookupServiceProviderDocumentationProvider{ShouldFail: true})

	// When
	documentation, err := service.GetDocumentation(context.Background(), "test-service")

	// Then
	require.Error(t, err)
	require.Empty(t, documentation)
	var target app.Error
	require.True(t, errors.As(err, &target))
	require.Equal(t, app.ErrorTypeProviderFailure, target.ErrorType())
	require.Equal(t, "unable to retrieve documentation for lookup service provider", target.Error())
}

func TestNewLookupServiceProviderDocumentationService_NilProvider(t *testing.T) {
	// When, Then
	require.Panics(t, func() {
		app.NewLookupServiceProviderDocumentationService(nil)
	})
}
