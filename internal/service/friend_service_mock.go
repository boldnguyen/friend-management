package service

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockFriendService is a mock implementation of FriendService interface.
type MockFriendService struct {
	mock.Mock
}

// CreateFriend mocks the CreateFriend method of the FriendService interface.
func (m *MockFriendService) CreateFriend(ctx context.Context, email1, email2 string) error {
	args := m.Called(ctx, email1, email2)
	return args.Error(0)
}

// GetFriendsList mocks the GetFriendsList method of the FriendService interface.
func (m *MockFriendService) GetFriendsList(ctx context.Context, email string) ([]string, error) {
	args := m.Called(ctx, email)
	return args.Get(0).([]string), args.Error(1)
}

// GetCommonFriends mocks the GetCommonFriends method of the FriendService interface.
func (m *MockFriendService) GetCommonFriends(ctx context.Context, email1, email2 string) ([]string, error) {
	args := m.Called(ctx, email1, email2)
	return args.Get(0).([]string), args.Error(1)
}

// SubscribeUpdates mocks the SubscribeUpdates method of the FriendService interface.
func (m *MockFriendService) SubscribeUpdates(ctx context.Context, requestor, target string) error {
	args := m.Called(ctx, requestor, target)
	return args.Error(0)
}

// BlockUpdates mocks the BlockUpdates method of the FriendService interface.
func (m *MockFriendService) BlockUpdates(ctx context.Context, requestor, target string) error {
	args := m.Called(ctx, requestor, target)
	return args.Error(0)
}

// GetEligibleRecipients mocks the GetEligibleRecipients method of the FriendService interface.
func (m *MockFriendService) GetEligibleRecipients(ctx context.Context, senderEmail, text string) ([]string, error) {
	args := m.Called(ctx, senderEmail, text)
	return args.Get(0).([]string), args.Error(1)
}
