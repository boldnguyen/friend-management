package service

import (
	"context"
	"log"

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
		// If the user doesn't exist, create the user and obtain the user ID
		log.Printf("Creating user with email %s", emails[0])
		userID1, err = s.repo.CreateUser(ctx, emails[0])
		if err != nil {
			log.Printf("Error creating user with email %s: %v", emails[0], err)
			return errors.Wrap(err, "failed to create user")
		}
	}

	userID2, err := s.repo.GetUserIDByEmail(ctx, emails[1])
	if err != nil {
		// If the user doesn't exist, create the user and obtain the user ID
		log.Printf("Creating user with email %s", emails[1])
		userID2, err = s.repo.CreateUser(ctx, emails[1])
		if err != nil {
			log.Printf("Error creating user with email %s: %v", emails[1], err)
			return errors.Wrap(err, "failed to create user")
		}
	}

	// Add friend connection using user IDs
	err = s.repo.AddFriend(ctx, userID1, userID2)
	if err != nil {
		log.Printf("Error adding friend connection: %v", err)
		return errors.Wrap(err, "failed to add friend")
	}

	return nil
}
