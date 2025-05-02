package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/stretchr/testify/require"
)

type mockTopicManagerDocumentationProvider struct {
	documentation string
	err           error
}

func (m *mockTopicManagerDocumentationProvider) GetDocumentationForTopicManager(topicManager string) (string, error) {
	return m.documentation, m.err
}

func TestTopicManagerDocumentationService_GetDocumentation_Success(t *testing.T) {
	// Given
	expectedDocumentation := "# Test Documentation\nThis is a test markdown document."
	mockProvider := &mockTopicManagerDocumentationProvider{
		documentation: expectedDocumentation,
		err:           nil,
	}
	service := app.NewTopicManagerDocumentationService(mockProvider)

	// When
	documentation, err := service.GetDocumentation(context.Background(), "testTopicManager")

	// Then
	require.NoError(t, err)
	require.Equal(t, expectedDocumentation, documentation)
}

func TestTopicManagerDocumentationService_GetDocumentation_EmptyTopicManagerName(t *testing.T) {
	// Given
	mockProvider := &mockTopicManagerDocumentationProvider{
		documentation: "# Test Documentation\nThis is a test markdown document.",
		err:           nil,
	}
	service := app.NewTopicManagerDocumentationService(mockProvider)

	// When
	_, err := service.GetDocumentation(context.Background(), "")

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, app.ErrEmptyTopicManagerName)
}

func TestTopicManagerDocumentationService_GetDocumentation_ProviderError(t *testing.T) {
	// Given
	mockError := errors.New("provider error")
	mockProvider := &mockTopicManagerDocumentationProvider{
		documentation: "",
		err:           mockError,
	}
	service := app.NewTopicManagerDocumentationService(mockProvider)

	// When
	_, err := service.GetDocumentation(context.Background(), "testTopicManager")

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, app.ErrTopicManagerNotFound)
}

func TestNewTopicManagerDocumentationService_NilProvider(t *testing.T) {
	// Given, When, Then
	require.Panics(t, func() {
		app.NewTopicManagerDocumentationService(nil)
	})
}
