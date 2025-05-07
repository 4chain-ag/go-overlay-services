package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/stretchr/testify/require"
)

func TestGetProviderDocumentation_Success(t *testing.T) {
	// Given
	service := app.NewLookupProviderDocumentationService(&testabilities.MockProviderDocumentationProvider{ShouldFail: false})

	// When
	documentation, err := service.GetDocumentation(context.Background(), "test-service")

	// Then
	require.NoError(t, err)
	require.Equal(t, "# Test Documentation\nThis is a test markdown document.", documentation)
}

func TestGetProviderDocumentation_EmptyLookupService(t *testing.T) {
	// Given
	service := app.NewLookupProviderDocumentationService(&testabilities.MockProviderDocumentationProvider{ShouldFail: false})

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

func TestGetProviderDocumentation_ProviderError(t *testing.T) {
	// Given
	service := app.NewLookupProviderDocumentationService(&testabilities.MockProviderDocumentationProvider{ShouldFail: true})

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

func TestNewLookupProviderDocumentationService_NilProvider(t *testing.T) {
	// When, Then
	require.Panics(t, func() {
		app.NewLookupProviderDocumentationService(nil)
	})
}
