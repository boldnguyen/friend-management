package repository

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockFriendRepository is a mock implementation of the FriendRepository interface for testing purposes.
type MockFriendRepository struct {
	mock.Mock
}

// AddFriend mocks the AddFriend method of the FriendRepository interface.
func (m *MockFriendRepository) AddFriend(ctx context.Context, userID1, userID2 int) error {
	args := m.Called(ctx, userID1, userID2)
	return args.Error(0)
}

// GetUserIDByEmail mocks the GetUserIDByEmail method of the FriendRepository interface.
func (m *MockFriendRepository) GetUserIDByEmail(ctx context.Context, email string) (int, error) {
	args := m.Called(ctx, email)
	return args.Int(0), args.Error(1)
}

// AreFriends mocks the AreFriends method of the FriendRepository interface.
func (m *MockFriendRepository) AreFriends(ctx context.Context, userID1, userID2 int) (bool, error) {
	args := m.Called(ctx, userID1, userID2)
	return args.Bool(0), args.Error(1)
}
