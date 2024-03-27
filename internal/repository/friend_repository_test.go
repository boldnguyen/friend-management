package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/boldnguyen/friend-management/internal/models"
	"github.com/boldnguyen/friend-management/internal/pkg/response"
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
		"already_friends": {
			userID1: 3,
			userID2: 4,
			expErr:  errors.New("they are already friends"),
		},
		"invalid_user_id": {
			userID1: -1,
			userID2: 5,
			expErr:  errors.New("invalid user ID"),
		},
	}

	mockRepo := &MockRepo{}

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

func TestFriendRepository_GetUserByEmail(t *testing.T) {
	tcs := map[string]struct {
		email  string
		expID  int
		expErr error
	}{
		"Valid email": {
			email:  "test@example.com",
			expID:  1,
			expErr: nil,
		},
		"Nonexistent email": {
			email:  "nonexistent@example.com",
			expID:  0,
			expErr: errors.New(response.ErrMsgUserNotFound),
		},
	}

	mockRepo := &MockRepo{}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			ctx := context.Background()

			// Mocking the behavior of GetUserByEmail method
			mockRepo.On("GetUserByEmail", ctx, tc.email).Return(&models.User{ID: tc.expID}, tc.expErr)

			// When
			id, err := mockRepo.GetUserByEmail(ctx, tc.email)

			// Then
			assert.Equal(t, tc.expID, id.ID, "Unexpected user ID")
			assert.Equal(t, tc.expErr, err, "Unexpected error")
			mockRepo.AssertExpectations(t)
		})
	}
}
func TestFriendRepository_CheckFriends(t *testing.T) {
	tcs := map[string]struct {
		userID1 int
		userID2 int
		expBool bool
		expErr  error
	}{
		"success_friends_direct": {
			userID1: 1,
			userID2: 2,
			expBool: true,
		},
		"success_friends_reverse": {
			userID1: 2,
			userID2: 1,
			expBool: true,
		},
		"success_not_friends": {
			userID1: 3,
			userID2: 4,
			expBool: false,
		},
		"success_not_friends_reverse": {
			userID1: 4,
			userID2: 3,
			expBool: false,
		},
		"error_database": {
			userID1: 5,
			userID2: 6,
			expBool: false,
			expErr:  errors.New("database error"),
		},
	}

	mockRepo := &MockRepo{}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			ctx := context.Background()

			// Mocking the behavior of CheckFriends method
			mockRepo.On("CheckFriends", ctx, tc.userID1, tc.userID2).Return(tc.expBool, tc.expErr)

			// When
			isFriends, err := mockRepo.CheckFriends(ctx, tc.userID1, tc.userID2)

			// Then
			assert.Equal(t, tc.expBool, isFriends, "Unexpected friendship status")
			assert.Equal(t, tc.expErr, err, "Unexpected error")
			mockRepo.AssertExpectations(t)
		})
	}
}
