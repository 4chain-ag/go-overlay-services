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
	expectations := testabilities.LookupServiceDocumentationProviderMockExpectations{
		DocumentationCall: true,
		Documentation:     "# Test Documentation\nThis is a test markdown document.",
	}
	mock := testabilities.NewLookupServiceDocumentationProviderMock(t, expectations)
	service := app.NewLookupDocumentationService(mock)

	// When
	documentation, err := service.GetDocumentation(context.Background(), "test-service")

	// Then
	require.NoError(t, err)
	require.Equal(t, expectations.Documentation, documentation)
	mock.AssertCalled()
}

func TestGetLookupServiceProviderDocumentation_EmptyLookupServiceName(t *testing.T) {
	// Given
	expectations := testabilities.LookupServiceDocumentationProviderMockExpectations{
		DocumentationCall: false,
		Error:             errors.New("lookup service name cannot be empty"),
	}
	mock := testabilities.NewLookupServiceDocumentationProviderMock(t, expectations)
	service := app.NewLookupDocumentationService(mock)
	expectedError := app.NewEmptyLookupServiceNameError()
	// When
	documentation, err := service.GetDocumentation(context.Background(), "")

	// Then
	require.Empty(t, documentation)

	var actualErr app.Error
	require.True(t, errors.As(err, &actualErr))
	require.Equal(t, expectedError, actualErr)

	mock.AssertCalled()
}

func TestGetLookupServiceProviderDocumentation_ProviderError(t *testing.T) {
	// Given
	expectations := testabilities.LookupServiceDocumentationProviderMockExpectations{
		DocumentationCall: true,
		Error:             errors.New("lookup service name cannot be empty"),
	}
	mock := testabilities.NewLookupServiceDocumentationProviderMock(t, expectations)
	service := app.NewLookupDocumentationService(mock)
	expectedError := app.NewLookupServiceProviderDocumentationError(expectations.Error)
	// When
	documentation, err := service.GetDocumentation(context.Background(), "test-service")

	// Then
	require.Empty(t, documentation)

	var actualErr app.Error
	require.True(t, errors.As(err, &actualErr))
	require.Equal(t, expectedError, actualErr)

	mock.AssertCalled()
}
