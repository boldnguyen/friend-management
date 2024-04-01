package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/boldnguyen/friend-management/internal/pkg/response"
	"github.com/boldnguyen/friend-management/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestNewHandler_CreateFriend tests the NewHandler function for creating a friend connection.
func TestNewHandler_CreateFriend(t *testing.T) {
	type mockService struct {
		expCall bool     // Whether the service method is expected to be called
		input   []string // Input data expected to be passed to the service method
		err     error    // Error expected to be returned by the service method
	}

	// Define test cases for different scenarios
	tcs := map[string]struct {
		input    []string    // Input friend data
		mockFn   mockService // Function to set up mock
		expCode  int         // expected HTTP response code
		expError string      // expected error message
	}{
		"success": {
			input:   []string{"test1@example.com", "test2@example.com"},
			mockFn:  mockService{expCall: true, input: []string{"test1@example.com", "test2@example.com"}, err: nil},
			expCode: http.StatusOK,
		},
		"already_friends": {
			input:    []string{"test1@example.com", "test2@example.com"},
			mockFn:   mockService{expCall: true, input: []string{"test1@example.com", "test2@example.com"}, err: errors.New("They are already friends")},
			expCode:  http.StatusInternalServerError,
			expError: response.ErrMsgAlreadyFriends,
		},
		"error": {
			input:    []string{"test1@example.com", "test2@example.com"},
			mockFn:   mockService{expCall: true, input: []string{"test1@example.com", "test2@example.com"}, err: errors.New("failed to create friend connection")},
			expCode:  http.StatusInternalServerError,
			expError: response.ErrMsgCreateFriend,
		},
	}

	// Iterate over each test case
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			mockFriendService := new(service.MockFriendService)
			if tc.mockFn.expCall {
				mockFriendService.On("CreateFriend", mock.Anything, tc.mockFn.input[0], tc.mockFn.input[1]).Return(tc.mockFn.err)
			}
			friendHandler := NewHandler(mockFriendService)

			// Marshal input data to JSON
			body, _ := json.Marshal(map[string][]string{"friends": tc.input})
			// Create HTTP request with the JSON payload
			req, err := http.NewRequest("POST", "/friend", bytes.NewBuffer(body))
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			// When
			// Call the handler function to handle the HTTP request
			friendHandler.ServeHTTP(rr, req)

			// Then
			// Assert the HTTP response code
			require.Equal(t, tc.expCode, rr.Code)
			// If an error message is expected, assert its presence in the response body
			if tc.expError != "" {
				require.True(t, strings.Contains(rr.Body.String(), tc.expError))
			}
			// Assert that the expected calls to the mock service were made
			mockFriendService.AssertExpectations(t)
		})
	}
}
