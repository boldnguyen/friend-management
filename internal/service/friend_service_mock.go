package service

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockFriendService is a mock implementation of the FriendService interface for testing purposes.
type MockFriendService struct {
	mock.Mock
}

// AddFriend mocks the AddFriend method of the FriendService interface.
func (m *MockFriendService) AddFriend(ctx context.Context, emails []string) error {
	args := m.Called(ctx, emails)
	return args.Error(0)
}

// GetUserIDByEmail mocks the GetUserIDByEmail method of the FriendService interface.
func (m *MockFriendService) GetUserIDByEmail(ctx context.Context, email string) (int, error) {
	args := m.Called(ctx, email)
	return args.Int(0), args.Error(1)
}

// AreFriends mocks the AreFriends method of the FriendService interface.
func (m *MockFriendService) AreFriends(ctx context.Context, userID1, userID2 int) (bool, error) {
	args := m.Called(ctx, userID1, userID2)
	return args.Bool(0), args.Error(1)
}
