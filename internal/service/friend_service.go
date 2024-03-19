// service.go
package service

import (
	"context"

	"github.com/boldnguyen/friend-management/internal/repository"
	"github.com/pkg/errors"
)

// FriendService defines the business logic for managing friend connections.
type FriendService interface {
	// AddFriend creates a friend connection between two users using their email addresses.
	AddFriend(ctx context.Context, emails []string) error
}

// friendService is the concrete implementation of FriendService.
type friendService struct {
	repo repository.FriendRepository
}

// NewFriendService creates a new instance of FriendService with the provided repository.
func NewFriendService(repo repository.FriendRepository) FriendService {
	return &friendService{repo: repo}
}

// AddFriend creates a friend connection between two users using their email addresses.
func (s *friendService) AddFriend(ctx context.Context, emails []string) error {
	// Validate the number of emails
	if len(emails) != 2 {
		return errors.New("Exactly two email addresses are required to add a friend connection")
	}

	// Get user IDs from emails
	userID1, err := s.repo.GetUserIDByEmail(ctx, emails[0])
	if err != nil {
		return errors.Wrap(err, "failed to get user ID for email ")
	}

	userID2, err := s.repo.GetUserIDByEmail(ctx, emails[1])
	if err != nil {
		return errors.Wrap(err, "failed to get user ID for email ")
	}

	// Check if the users are already friends
	alreadyFriends, err := s.repo.AreFriends(ctx, userID1, userID2)
	if err != nil {
		return errors.Wrap(err, "failed to check if users are already friends")
	}
	if alreadyFriends {
		return errors.New("They are already friends")
	}

	// Add friend connection using user IDs
	err = s.repo.AddFriend(ctx, userID1, userID2)
	if err != nil {
		return errors.Wrap(err, "failed to add friend")
	}

	return nil
}
