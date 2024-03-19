package service

import (
	"context"
	"testing"

	"github.com/boldnguyen/friend-management/internal/repository"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

// TestFriendService_AddFriend is a unit test function for testing the AddFriend method in the FriendService.
func TestFriendService_AddFriend(t *testing.T) {
	type mockFriendRepository struct {
		expCall bool
		emails  []string
		userID1 int
		userID2 int
		err     error
	}
	tcs := map[string]struct {
		emails   []string
		expErr   error
		mockFunc mockFriendRepository
	}{
		"success": {
			emails: []string{"test1@example.com", "test2@example.com"},
			mockFunc: mockFriendRepository{
				expCall: true,
				emails:  []string{"test1@example.com", "test2@example.com"},
				userID1: 1,
				userID2: 2,
				err:     nil,
			},
			expErr: nil,
		},
		"invalid_number_of_emails": {
			emails: []string{"test1@example.com"},
			expErr: errors.New("Exactly two email addresses are required to add a friend connection"),
		},
	}

	// Iterating over each test case to execute the test
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			ctx := context.Background()

			// Creating a mock instance of FriendRepository
			mockRepo := new(repository.MockFriendRepository)
			if tc.mockFunc.expCall {
				// Set up expectations for GetUserIDByEmail
				mockRepo.On("GetUserIDByEmail", ctx, tc.mockFunc.emails[0]).Return(0, tc.mockFunc.err)
				mockRepo.On("GetUserIDByEmail", ctx, tc.mockFunc.emails[1]).Return(1, nil)

				// Expect AreFriends with user IDs obtained from the mock repository
				mockRepo.On("AreFriends", ctx, 0, 1).Return(false, nil)

				// Expect AddFriend with the user IDs obtained from the mock repository
				mockRepo.On("AddFriend", ctx, 0, 1).Return(nil)
			}

			// Creating an instance of FriendService with the mock repository
			friendService := NewFriendService(mockRepo)

			// When
			err := friendService.AddFriend(ctx, tc.emails)

			// Then
			// Verifying that the expected repository method calls are made
			mockRepo.AssertExpectations(t)

			// Verifying the returned error matches the expected error
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
