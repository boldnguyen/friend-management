package service

import (
	"context"

	"github.com/boldnguyen/friend-management/internal/repository"
)

// FriendService provides methods for managing friend connections.
type FriendService interface {
	CreateFriend(ctx context.Context, email1, email2 string) error
	GetFriendsList(ctx context.Context, email string) ([]string, error)
}

// friendService implements the FriendService interface.
type friendService struct {
	repo repository.FriendRepository
}

// NewFriendService creates a new FriendService instance.
func NewFriendService(repo repository.FriendRepository) FriendService {
	return &friendService{repo: repo}
}
