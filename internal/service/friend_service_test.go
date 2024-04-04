package service

import (
	"context"
	"errors"
	"testing"

	"github.com/boldnguyen/friend-management/internal/models"
	"github.com/boldnguyen/friend-management/internal/pkg/response"
	"github.com/boldnguyen/friend-management/internal/repository"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestCreateFriend tests the CreateFriend method of the FriendService.
func TestCreateFriend(t *testing.T) {
	// mockRepoService is a struct to define the expected behavior of the mock repository.
	type mockRepoService struct {
		expGetUserByEmail map[string]*models.User // Expected GetUserByEmail return values
		expCheckFriends   bool                    // Expected CheckFriends return value
		expAddFriendErr   error                   // Expected AddFriend error
	}
	// Define test cases for different scenarios
	tcs := map[string]struct {
		email1   string          // First user's email
		email2   string          // Second user's email
		mockFn   mockRepoService // Function to set up mock
		expError string          // expected error message
	}{
		"success": {
			email1: "test1@example.com",
			email2: "test2@example.com",
			mockFn: mockRepoService{
				expGetUserByEmail: map[string]*models.User{
					"test1@example.com": {ID: 1},
					"test2@example.com": {ID: 2},
				},
				expCheckFriends: false,
				expAddFriendErr: nil,
			},
			expError: "",
		},
		"error_user_not_found": {
			email1: "test1@example.com",
			email2: "test2@example.com",
			mockFn: mockRepoService{
				expGetUserByEmail: map[string]*models.User{
					"test1@example.com": {ID: 1},
					"test2@example.com": {ID: 2},
				},
				expCheckFriends: false,
				expAddFriendErr: errors.New("user not found"),
			},
			expError: response.ErrMsgUserNotFound,
		},
		"error_already_friends": {
			email1: "test1@example.com",
			email2: "test2@example.com",
			mockFn: mockRepoService{
				expGetUserByEmail: map[string]*models.User{
					"test1@example.com": {ID: 1},
					"test2@example.com": {ID: 2},
				},
				expCheckFriends: true, // Users are already friends
				expAddFriendErr: nil,
			},
			expError: response.ErrMsgAlreadyFriends,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			mockRepo := new(repository.MockRepo)
			friendService := NewFriendService(mockRepo)

			// Set up mock expectations for GetUserByEmail
			mockRepo.On("GetUserByEmail", mock.Anything, tc.email1).Return(tc.mockFn.expGetUserByEmail[tc.email1], nil).Once()
			mockRepo.On("GetUserByEmail", mock.Anything, tc.email2).Return(tc.mockFn.expGetUserByEmail[tc.email2], nil).Once()

			// Set up mock expectations for CheckFriends
			mockRepo.On("CheckFriends", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(tc.mockFn.expCheckFriends, nil).Once()

			// Expect AddFriend only if not already friends and user found
			if !tc.mockFn.expCheckFriends && tc.mockFn.expGetUserByEmail[tc.email1] != nil && tc.mockFn.expGetUserByEmail[tc.email2] != nil {
				mockRepo.On("AddFriend", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(tc.mockFn.expAddFriendErr).Once()
			}

			// When
			err := friendService.CreateFriend(context.Background(), tc.email1, tc.email2)

			// Then
			if tc.expError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expError)
			} else {
				require.NoError(t, err)
			}

			// Assert that the expected calls to the mock repository were made
			mockRepo.AssertExpectations(t)
		})
	}

}

// TestGetFriendsList tests the GetFriendsList method of the FriendService.
func TestGetFriendsList(t *testing.T) {
	// mockRepoService is a struct to define the expected behavior of the mock repository.
	type mockRepoService struct {
		expGetUserByEmail map[string]*models.User // Expected GetUserByEmail return values
		expGetFriendsList []string                // Expected GetFriendsList return value
		expErr            error                   // Expected error
	}
	// Define test cases for different scenarios
	tcs := map[string]struct {
		email    string          // User's email
		mockFn   mockRepoService // Function to set up mock
		expError string          // expected error message
		expList  []string        // Expected list of friends
	}{
		"success": {
			email: "test@example.com",
			mockFn: mockRepoService{
				expGetUserByEmail: map[string]*models.User{
					"test@example.com": {ID: 1},
				},
				expGetFriendsList: []string{"friend1@example.com", "friend2@example.com"},
				expErr:            nil,
			},
			expError: "",
			expList:  []string{"friend1@example.com", "friend2@example.com"},
		},
		"user_not_found": {
			email: "test@example.com",
			mockFn: mockRepoService{
				expGetUserByEmail: map[string]*models.User{
					"test@example.com": nil, // User not found
				},
				expGetFriendsList: nil,
				expErr:            errors.New("user not found"),
			},
			expError: response.ErrMsgUserNotFound,
			expList:  nil,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			mockRepo := new(repository.MockRepo)
			friendService := NewFriendService(mockRepo)

			// Set up mock expectations for GetUserByEmail
			mockRepo.On("GetUserByEmail", mock.Anything, tc.email).Return(tc.mockFn.expGetUserByEmail[tc.email], nil).Once()

			if tc.mockFn.expGetUserByEmail[tc.email] != nil {
				// Set up mock expectations for GetFriendsList if user found
				mockRepo.On("GetFriendsList", mock.Anything, mock.AnythingOfType("int")).Return(tc.mockFn.expGetFriendsList, tc.mockFn.expErr).Once()
			}

			// When
			list, err := friendService.GetFriendsList(context.Background(), tc.email)

			// Then
			if tc.expError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expError)
				require.Nil(t, list)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expList, list)
			}

			// Assert that the expected calls to the mock repository were made
			mockRepo.AssertExpectations(t)
		})
	}

}
