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

// TestFriendListHandler_GetFriendsList tests the FriendListHandler function for retrieving the friends list.
func TestFriendListHandler_GetFriendsList(t *testing.T) {
	type mockService struct {
		expCall bool     // Whether the service method is expected to be called
		input   string   // Input data expected to be passed to the service method
		output  []string // Output data to be returned by the service method
		err     error    // Error expected to be returned by the service method
	}

	// Define test cases for different scenarios
	tcs := map[string]struct {
		input    string      // Input email data
		mockFn   mockService // Function to set up mock
		expCode  int         // expected HTTP response code
		expError string      // expected error message
	}{
		"success": {
			input:   "test@example.com",
			mockFn:  mockService{expCall: true, input: "test@example.com", output: []string{"friend1@example.com", "friend2@example.com"}, err: nil},
			expCode: http.StatusOK,
		},
		"no_friends": {
			input:    "test@example.com",
			mockFn:   mockService{expCall: true, input: "test@example.com", output: []string{}, err: nil},
			expCode:  http.StatusOK,
			expError: "",
		},
		"error": {
			input:    "test@example.com",
			mockFn:   mockService{expCall: true, input: "test@example.com", output: nil, err: errors.New("failed to get friends list")},
			expCode:  http.StatusInternalServerError,
			expError: response.ErrMsgGetFriendsList,
		},
	}

	// Iterate over each test case
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			mockFriendService := new(service.MockFriendService)
			if tc.mockFn.expCall {
				mockFriendService.On("GetFriendsList", mock.Anything, tc.mockFn.input).Return(tc.mockFn.output, tc.mockFn.err)
			}
			friendListHandler := FriendListHandler(mockFriendService)

			// Marshal input data to JSON
			body, _ := json.Marshal(map[string]string{"email": tc.input})
			// Create HTTP request with the JSON payload
			req, err := http.NewRequest("POST", "/friend/list", bytes.NewBuffer(body))
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			// When
			// Call the handler function to handle the HTTP request
			friendListHandler.ServeHTTP(rr, req)

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

// TestCommonFriendsHandler_GetCommonFriends tests the CommonFriendsHandler function for retrieving common friends list.
func TestCommonFriendsHandler_GetCommonFriends(t *testing.T) {
	type mockService struct {
		expCall bool     // Whether the service method is expected to be called
		input   []string // Input data expected to be passed to the service method
		output  []string // Output data to be returned by the service method
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
			mockFn:  mockService{expCall: true, input: []string{"test1@example.com", "test2@example.com"}, output: []string{"common@example.com"}, err: nil},
			expCode: http.StatusOK,
		},
		"no_common_friends": {
			input:    []string{"test1@example.com", "test2@example.com"},
			mockFn:   mockService{expCall: true, input: []string{"test1@example.com", "test2@example.com"}, output: []string{}, err: nil},
			expCode:  http.StatusOK,
			expError: "",
		},
		"error": {
			input:    []string{"test1@example.com", "test2@example.com"},
			mockFn:   mockService{expCall: true, input: []string{"test1@example.com", "test2@example.com"}, output: nil, err: errors.New("failed to retrieve common friends list")},
			expCode:  http.StatusInternalServerError,
			expError: response.ErrMsgGetCommonFriends,
		},
	}

	// Iterate over each test case
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			mockFriendService := new(service.MockFriendService)
			if tc.mockFn.expCall {
				mockFriendService.On("GetCommonFriends", mock.Anything, tc.mockFn.input[0], tc.mockFn.input[1]).Return(tc.mockFn.output, tc.mockFn.err)
			}
			commonFriendsHandler := CommonFriendsHandler(mockFriendService)

			// Marshal input data to JSON
			body, _ := json.Marshal(map[string][]string{"friends": tc.input})
			// Create HTTP request with the JSON payload
			req, err := http.NewRequest("POST", "/friend/common", bytes.NewBuffer(body))
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			// When
			// Call the handler function to handle the HTTP request
			commonFriendsHandler.ServeHTTP(rr, req)

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

// TestSubscribeHandler_SubscribeUpdates tests the SubscribeHandler function for subscribing to updates.
func TestSubscribeHandler_SubscribeUpdates(t *testing.T) {
	type mockService struct {
		expCall bool             // Whether the service method is expected to be called
		input   SubscribeRequest // Input data expected to be passed to the service method
		err     error            // Error expected to be returned by the service method
	}

	// Define test cases for different scenarios
	tcs := map[string]struct {
		req      SubscribeRequest // Input request data
		mockFn   mockService      // Function to set up mock
		expCode  int              // expected HTTP response code
		expError string           // expected error message
	}{
		"success": {
			req:     SubscribeRequest{Requestor: "subscriber@example.com", Target: "target@example.com"},
			mockFn:  mockService{expCall: true, input: SubscribeRequest{Requestor: "subscriber@example.com", Target: "target@example.com"}, err: nil},
			expCode: http.StatusOK,
		},
		"invalid_json": {
			req:      SubscribeRequest{},                                              // Invalid JSON data
			mockFn:   mockService{expCall: true, input: SubscribeRequest{}, err: nil}, // Mock function should be called
			expCode:  http.StatusBadRequest,
			expError: response.ErrMsgInvalidRequest,
		},

		// Add more test cases for other scenarios
	}

	// Iterate over each test case
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			mockFriendService := new(service.MockFriendService)
			if tc.mockFn.expCall {
				mockFriendService.On("SubscribeUpdates", mock.Anything, tc.mockFn.input.Requestor, tc.mockFn.input.Target).Return(tc.mockFn.err)
			}
			subscribeHandler := SubscribeHandler(mockFriendService)

			// Marshal input data to JSON
			body, _ := json.Marshal(tc.req)
			// Create HTTP request with the JSON payload
			req, err := http.NewRequest("POST", "/subscribe", bytes.NewBuffer(body))
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			// When
			// Call the handler function to handle the HTTP request
			subscribeHandler.ServeHTTP(rr, req)

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
