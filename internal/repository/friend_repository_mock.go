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

// GetCommonFriends mocks the GetCommonFriends method.
func (m *MockRepo) GetCommonFriends(ctx context.Context, userID1, userID2 int) ([]string, error) {
	args := m.Called(ctx, userID1, userID2)
	return args.Get(0).([]string), args.Error(1)
}

// SubscribeUpdates mocks the SubscribeUpdates method.
func (m *MockRepo) SubscribeUpdates(ctx context.Context, requestor, target string) error {
	args := m.Called(ctx, requestor, target)
	return args.Error(0)
}

// CheckSubscription mocks the CheckSubscription method.
func (m *MockRepo) CheckSubscription(ctx context.Context, requestor, target string) (bool, error) {
	args := m.Called(ctx, requestor, target)
	return args.Bool(0), args.Error(1)
}

// DeleteSubscription mocks the DeleteSubscription method.
func (m *MockRepo) DeleteSubscription(ctx context.Context, requestorID, targetID int) error {
	args := m.Called(ctx, requestorID, targetID)
	return args.Error(0)
}

// BlockUser mocks the BlockUser method.
func (m *MockRepo) BlockUser(ctx context.Context, requestorID, targetID int) error {
	args := m.Called(ctx, requestorID, targetID)
	return args.Error(0)
}
