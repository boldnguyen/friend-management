package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFriendRepository_AddFriend(t *testing.T) {
	tcs := map[string]struct {
		userID1 int
		userID2 int
		expErr  error
	}{
		"success": {
			userID1: 1,
			userID2: 2,
		},
	}

	mockRepo := &MockFriendRepository{} // Use your mock implementation here

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			ctx := context.Background()

			// Mocking the behavior of AddFriend method
			mockRepo.On("AddFriend", ctx, tc.userID1, tc.userID2).Return(tc.expErr)

			// When
			err := mockRepo.AddFriend(ctx, tc.userID1, tc.userID2)

			// Then
			assert.Equal(t, tc.expErr, err, "Unexpected error")
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestFriendRepository_GetUserIDByEmail(t *testing.T) {
	tcs := map[string]struct {
		email  string
		expID  int
		expErr error
	}{
		"success": {
			email: "test@example.com",
			expID: 123,
		},
		// Add more test cases as needed
	}

	mockRepo := &MockFriendRepository{} // Use your mock implementation here

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			ctx := context.Background()

			// Mocking the behavior of GetUserIDByEmail method
			mockRepo.On("GetUserIDByEmail", ctx, tc.email).Return(tc.expID, tc.expErr)

			// When
			id, err := mockRepo.GetUserIDByEmail(ctx, tc.email)

			// Then
			assert.Equal(t, tc.expID, id, "Unexpected user ID")
			assert.Equal(t, tc.expErr, err, "Unexpected error")
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestFriendRepository_AreFriends(t *testing.T) {
	tcs := map[string]struct {
		userID1 int
		userID2 int
		expBool bool
		expErr  error
	}{
		"success": {
			userID1: 1,
			userID2: 2,
			expBool: true,
		},
	}

	mockRepo := &MockFriendRepository{} // Use your mock implementation here

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			ctx := context.Background()

			// Mocking the behavior of AreFriends method
			mockRepo.On("AreFriends", ctx, tc.userID1, tc.userID2).Return(tc.expBool, tc.expErr)

			// When
			isFriends, err := mockRepo.AreFriends(ctx, tc.userID1, tc.userID2)

			// Then
			assert.Equal(t, tc.expBool, isFriends, "Unexpected friendship status")
			assert.Equal(t, tc.expErr, err, "Unexpected error")
			mockRepo.AssertExpectations(t)
		})
	}
}
