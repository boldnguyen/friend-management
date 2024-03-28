package service

import (
	"context"
	"log"

	"github.com/boldnguyen/friend-management/internal/pkg/response"
	"github.com/pkg/errors"
)

// CreateFriend creates a friend connection between two users using their email addresses.
func (serv friendService) CreateFriend(ctx context.Context, email1, email2 string) error {
	// Get user IDs from emails
	user1, err := serv.repo.GetUserByEmail(ctx, email1)
	if err != nil {
		log.Printf("Failed to get user for email %s: %v", email1, err)
		return errors.Wrap(err, response.ErrMsgGetUserByEmail)
	}
	if user1 == nil {
		return errors.New(response.ErrMsgUserNotFound)
	}
	userID1 := user1.ID

	user2, err := serv.repo.GetUserByEmail(ctx, email2)
	if err != nil {
		log.Printf("Failed to get user for email %s: %v", email2, err)
		return errors.Wrap(err, response.ErrMsgGetUserByEmail)
	}
	if user2 == nil {
		return errors.New(response.ErrMsgUserNotFound)
	}
	userID2 := user2.ID

	// Check if the users are already friends
	alreadyFriends, err := serv.repo.CheckFriends(ctx, userID1, userID2)
	if err != nil {
		return errors.Wrap(err, response.ErrMsgCheckFriend)
	}

	if alreadyFriends {
		return errors.New(response.ErrMsgAlreadyFriends)
	}

	// Add friend connection using user IDs
	err = serv.repo.AddFriend(ctx, userID1, userID2)
	if err != nil {
		log.Printf("Failed to create friend connection: %v", err)
		return errors.Wrap(err, response.ErrMsgCreateFriend)
	}

	return nil
}
