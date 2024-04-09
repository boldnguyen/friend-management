package repository

import (
	"context"

	"github.com/boldnguyen/friend-management/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockFriendRepository is a mock implementation of FriendRepository for testing purposes.
type MockRepo struct {
	mock.Mock
}

// GetUserByEmail mocks the GetUserByEmail method.
func (m *MockRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*models.User), args.Error(1)
}

// AddFriend mocks the AddFriend method.
func (m *MockRepo) AddFriend(ctx context.Context, userID1, userID2 int) error {
	args := m.Called(ctx, userID1, userID2)
	return args.Error(0)
}

// CheckFriends mocks the CheckFriends method.
func (m *MockRepo) CheckFriends(ctx context.Context, userID1, userID2 int) (bool, error) {
	args := m.Called(ctx, userID1, userID2)
	return args.Bool(0), args.Error(1)
}

// GetFriendsList mocks the GetFriendsList method.
func (m *MockRepo) GetFriendsList(ctx context.Context, userID int) ([]string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

// GetCommonFriends is a mock implementation of retrieving common friends.
func (r *MockRepo) GetCommonFriends(ctx context.Context, userID1, userID2 int) ([]string, error) {
	return []string{}, nil
}
