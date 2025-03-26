package server_test

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

// SetupMockResty sets up a Resty client with a mocked HTTP server
func SetupMockResty(t *testing.T, method string, mockURL string, statusCode int, body string) *resty.Client {
	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())
	t.Cleanup(func() { httpmock.DeactivateAndReset() })

	httpmock.RegisterResponder(method, mockURL,
		httpmock.NewStringResponder(statusCode, body),
	)

	return client
}

func TestRestyWithPost_ShouldReturnCreated(t *testing.T) {
	// Given:
	mockURL := "https://api.example.com/create"
	expectedResponse := `{"status": "created"}`
	client := SetupMockResty(t, http.MethodPost, mockURL, 201, expectedResponse)

	// When:
	resp, err := client.R().Post(mockURL)

	// Then:
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode())
	require.JSONEq(t, expectedResponse, resp.String())
}

func TestRestyWithGet_ShouldReturnOK(t *testing.T) {
	// Given:
	mockURL := "https://api.example.com/data"
	expectedResponse := `{"key": "value"}`
	client := SetupMockResty(t, http.MethodGet, mockURL, 200, expectedResponse)

	// When:
	resp, err := client.R().Get(mockURL)

	// Then:
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode())
	require.JSONEq(t, expectedResponse, resp.String())
}

func TestRestyWithPut_ShouldReturnAccepted(t *testing.T) {
	// Given:
	mockURL := "https://api.example.com/update"
	expectedResponse := `{"status": "updated"}`
	client := SetupMockResty(t, http.MethodPut, mockURL, 202, expectedResponse)

	// When:
	resp, err := client.R().Put(mockURL)

	// Then:
	require.NoError(t, err)
	require.Equal(t, 202, resp.StatusCode())
	require.JSONEq(t, expectedResponse, resp.String())
}

func TestRestyWithPatch_ShouldReturnOK(t *testing.T) {
	// Given:
	mockURL := "https://api.example.com/patch"
	expectedResponse := `{"status": "patched"}`
	client := SetupMockResty(t, http.MethodPatch, mockURL, 200, expectedResponse)

	// When:
	resp, err := client.R().Patch(mockURL)

	// Then:
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode())
	require.JSONEq(t, expectedResponse, resp.String())
}

func TestRestyWithDelete_ShouldReturnNoContent(t *testing.T) {
	// Given:
	mockURL := "https://api.example.com/delete"
	client := SetupMockResty(t, http.MethodDelete, mockURL, 204, "")

	// When:
	resp, err := client.R().Delete(mockURL)

	// Then:
	require.NoError(t, err)
	require.Equal(t, 204, resp.StatusCode())
	require.Empty(t, resp.String())
}
