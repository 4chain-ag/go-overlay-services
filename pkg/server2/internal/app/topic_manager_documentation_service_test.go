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
	expectations := testabilities.TopicManagerDocumentationProviderMockExpectations{
		DocumentationCall: true,
		Documentation:     "# Test Documentation\nThis is a test markdown document.",
	}
	mock := testabilities.NewTopicManagerDocumentationProviderMock(t, expectations)
	service := app.NewTopicManagerDocumentationService(mock)

	// When
	documentation, err := service.GetDocumentation(context.Background(), "test-topic-manager")

	// Then
	require.NoError(t, err)
	require.Equal(t, expectations.Documentation, documentation)
	mock.AssertCalled()
}

func TestGetTopicManagerDocumentation_EmptyTopicManagerName(t *testing.T) {
	// Given
	expectations := testabilities.TopicManagerDocumentationProviderMockExpectations{
		DocumentationCall: false,
		Error:             errors.New("topic manager name cannot be empty"),
	}
	mock := testabilities.NewTopicManagerDocumentationProviderMock(t, expectations)
	service := app.NewTopicManagerDocumentationService(mock)
	expectedError := app.NewEmptyTopicManagerNameError()
	// When
	documentation, err := service.GetDocumentation(context.Background(), "")

	// Then
	require.Empty(t, documentation)

	var actualErr app.Error
	require.True(t, errors.As(err, &actualErr))
	require.Equal(t, expectedError, actualErr)

	mock.AssertCalled()
}

func TestGetTopicManagerDocumentation_ProviderError(t *testing.T) {
	// Given
	expectations := testabilities.TopicManagerDocumentationProviderMockExpectations{
		DocumentationCall: true,
		Error:             errors.New("topic manager name cannot be empty"),
	}
	mock := testabilities.NewTopicManagerDocumentationProviderMock(t, expectations)
	service := app.NewTopicManagerDocumentationService(mock)
	expectedError := app.NewTopicManagerDocumentationError(expectations.Error)
	// When
	documentation, err := service.GetDocumentation(context.Background(), "test-topic-manager")

	// Then
	require.Empty(t, documentation)

	var actualErr app.Error
	require.True(t, errors.As(err, &actualErr))
	require.Equal(t, expectedError.ErrorType(), actualErr.ErrorType())
	require.Equal(t, expectedError.Error(), actualErr.Error())

	mock.AssertCalled()
}
