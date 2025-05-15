package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/stretchr/testify/require"
)

func TestGetTopicManagerDocumentation_Success(t *testing.T) {
	// Given
	service := app.NewTopicManagerDocumentationService(&testabilities.MockTopicManagerDocumentationProvider{ShouldFail: false})

	// When
	documentation, err := service.GetDocumentation(context.Background(), "test-topic-manager")

	// Then
	require.NoError(t, err)
	require.Equal(t, "# Topic Manager Documentation\nThis is a test markdown document.", documentation)
}

func TestGetTopicManagerDocumentation_EmptyTopicManagerName(t *testing.T) {
	// Given
	service := app.NewTopicManagerDocumentationService(&testabilities.MockTopicManagerDocumentationProvider{ShouldFail: false})

	// When
	documentation, err := service.GetDocumentation(context.Background(), "")

	// Then
	require.Error(t, err)
	require.Empty(t, documentation)
	var target app.Error
	require.True(t, errors.As(err, &target))
	require.Equal(t, app.ErrorTypeIncorrectInput, target.ErrorType())
	require.Equal(t, "topic manager name cannot be empty", target.Error())
}

func TestGetTopicManagerDocumentation_ProviderError(t *testing.T) {
	// Given
	service := app.NewTopicManagerDocumentationService(&testabilities.MockTopicManagerDocumentationProvider{ShouldFail: true})

	// When
	documentation, err := service.GetDocumentation(context.Background(), "test-topic-manager")

	// Then
	require.Error(t, err)
	require.Empty(t, documentation)
	var target app.Error
	require.True(t, errors.As(err, &target))
	require.Equal(t, app.ErrorTypeProviderFailure, target.ErrorType())
	require.Equal(t, "unable to retrieve documentation for topic manager", target.Error())
}

func TestNewTopicManagerDocumentationService_NilProvider(t *testing.T) {
	// When, Then
	require.Panics(t, func() {
		app.NewTopicManagerDocumentationService(nil)
	})
}
