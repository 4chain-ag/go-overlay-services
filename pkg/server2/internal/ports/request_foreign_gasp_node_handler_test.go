package ports_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestRequestForeignGASPNodeHandler_InvalidCases(t *testing.T) {

	tests := map[string]struct {
		payload interface{}

		headers map[string]string

		expectations testabilities.RequestForeignGASPNodeProviderMockExpectations

		expectedStatusCode int
	}{

		"missing topic header": {

			payload: map[string]interface{}{

				"graphID": "0000000000000000000000000000000000000000000000000000000000000000.1",

				"txID": "0000000000000000000000000000000000000000000000000000000000000000",

				"outputIndex": 1,
			},

			headers: map[string]string{

				fiber.HeaderContentType: fiber.MIMEApplicationJSON,
			},

			expectedStatusCode: fiber.StatusBadRequest,
		},

		"invalid JSON": {

			payload: "INVALID_JSON",

			headers: map[string]string{

				fiber.HeaderContentType: fiber.MIMEApplicationJSON,

				ports.XBSVTopicHeader: "test-topic",
			},

			expectedStatusCode: fiber.StatusBadRequest,
		},

		"invalid txID": {

			payload: map[string]interface{}{

				"graphID": "0000000000000000000000000000000000000000000000000000000000000000.1",

				"txID": "INVALID_TXID",

				"outputIndex": 1,
			},

			headers: map[string]string{

				fiber.HeaderContentType: fiber.MIMEApplicationJSON,

				ports.XBSVTopicHeader: "test-topic",
			},

			expectedStatusCode: fiber.StatusBadRequest,
		},

		"invalid graphID": {

			payload: map[string]interface{}{

				"graphID": "INVALID_GRAPH_ID",

				"txID": "0000000000000000000000000000000000000000000000000000000000000000",

				"outputIndex": 1,
			},

			headers: map[string]string{

				fiber.HeaderContentType: fiber.MIMEApplicationJSON,

				ports.XBSVTopicHeader: "test-topic",
			},

			expectedStatusCode: fiber.StatusBadRequest,
		},

		"provider error": {

			payload: map[string]interface{}{

				"graphID": "0000000000000000000000000000000000000000000000000000000000000000.1",

				"txID": "0000000000000000000000000000000000000000000000000000000000000000",

				"outputIndex": 1,
			},

			headers: map[string]string{

				fiber.HeaderContentType: fiber.MIMEApplicationJSON,

				ports.XBSVTopicHeader: "test-topic",
			},

			expectations: testabilities.RequestForeignGASPNodeProviderMockExpectations{

				ProvideForeignGASPNodeCall: true,

				Error: errors.New("provider error"),
			},

			expectedStatusCode: fiber.StatusInternalServerError,
		},
	}

	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {

			// given:

			stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithRequestForeignGASPNodeProvider(

				testabilities.NewRequestForeignGASPNodeProviderMock(t, tc.expectations),
			))

			fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

			// when:

			res, _ := fixture.Client().
				R().
				SetHeaders(tc.headers).
				SetBody(tc.payload).
				Post("/api/v1/requestForeignGASPNode")

			// then:

			require.Equal(t, tc.expectedStatusCode, res.StatusCode())

			if tc.expectations.ProvideForeignGASPNodeCall {

				stub.AssertProvidersState()

			}

		})

	}

}

func TestRequestForeignGASPNodeHandler_ValidCase(t *testing.T) {

	// given:

	expectedNode := &core.GASPNode{}

	expectations := testabilities.RequestForeignGASPNodeProviderMockExpectations{

		ProvideForeignGASPNodeCall: true,

		Node: expectedNode,
	}

	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithRequestForeignGASPNodeProvider(

		testabilities.NewRequestForeignGASPNodeProviderMock(t, expectations),
	))

	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

	// when:

	var actualNode core.GASPNode

	res, _ := fixture.Client().
		R().
		SetHeader(fiber.HeaderContentType, fiber.MIMEApplicationJSON).
		SetHeader(ports.XBSVTopicHeader, "test-topic").
		SetBody(map[string]interface{}{

			"graphID": "0000000000000000000000000000000000000000000000000000000000000000.1",

			"txID": "0000000000000000000000000000000000000000000000000000000000000000",

			"outputIndex": 1,
		}).
		SetResult(&actualNode).
		Post("/api/v1/requestForeignGASPNode")

	// then:

	require.Equal(t, http.StatusOK, res.StatusCode())

	require.Equal(t, *expectedNode, actualNode)

	stub.AssertProvidersState()

}
