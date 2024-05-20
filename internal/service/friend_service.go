package service

import (
	"context"
	"log"
	"strings"

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

// GetFriendsList retrieves the list of friends for a given email address.
func (serv friendService) GetFriendsList(ctx context.Context, email string) ([]string, error) {
	// Get user ID from email
	user, err := serv.repo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("Failed to get user for email %s: %v", email, err)
		return nil, errors.Wrap(err, response.ErrMsgGetUserByEmail)
	}
	if user == nil {
		return nil, errors.New(response.ErrMsgUserNotFound)
	}
	userID := user.ID

	// Get friends list
	friends, err := serv.repo.GetFriendsList(ctx, userID)
	if err != nil {
		log.Printf("Failed to get friends list for user ID %d: %v", userID, err)
		return nil, errors.Wrap(err, response.ErrMsgGetFriendsList)
	}

	return friends, nil
}

// GetCommonFriends retrieves the list of common friends between two email addresses.
func (serv friendService) GetCommonFriends(ctx context.Context, email1, email2 string) ([]string, error) {
	// Get user IDs from emails
	user1, err := serv.repo.GetUserByEmail(ctx, email1)
	if err != nil {
		log.Printf("Failed to get user for email %s: %v", email1, err)
		return nil, errors.Wrap(err, response.ErrMsgGetUserByEmail)
	}
	if user1 == nil {
		return nil, errors.New(response.ErrMsgUserNotFound)
	}
	userID1 := user1.ID

	user2, err := serv.repo.GetUserByEmail(ctx, email2)
	if err != nil {
		log.Printf("Failed to get user for email %s: %v", email2, err)
		return nil, errors.Wrap(err, response.ErrMsgGetUserByEmail)
	}
	if user2 == nil {
		return nil, errors.New(response.ErrMsgUserNotFound)
	}
	userID2 := user2.ID

	// Get common friends list
	commonFriends, err := serv.repo.GetCommonFriends(ctx, userID1, userID2)
	if err != nil {
		log.Printf("Failed to get common friends list: %v", err)
		return nil, errors.Wrap(err, response.ErrMsgGetCommonFriends)
	}

	return commonFriends, nil
}

// SubscribeUpdates subscribes requestor to updates from target.
func (serv *friendService) SubscribeUpdates(ctx context.Context, requestor, target string) error {
	// Check if the subscription already exists
	exists, err := serv.repo.CheckSubscription(ctx, requestor, target)
	if err != nil {
		return errors.Wrap(err, response.ErrMsgCheckSubscription)
	}
	if exists {
		return errors.New(response.ErrMsgAlreadySubscribed)
	}

	// Subscribe to updates
	err = serv.repo.SubscribeUpdates(ctx, requestor, target)
	if err != nil {
		log.Printf("Failed to subscribe updates: %v", err)
		return errors.Wrap(err, response.ErrMsgSubscribeUpdates)
	}
	return nil
}

// BlockUpdates handles the business logic for blocking updates.
func (serv *friendService) BlockUpdates(ctx context.Context, requestorEmail, targetEmail string) error {
	// Retrieve users by email
	requestorUser, err := serv.repo.GetUserByEmail(ctx, requestorEmail)
	if err != nil {
		return errors.Wrap(err, response.ErrMsgGetUserByEmail)
	}

	targetUser, err := serv.repo.GetUserByEmail(ctx, targetEmail)
	if err != nil {
		return errors.Wrap(err, response.ErrMsgGetUserByEmail)
	}

	requestorID := requestorUser.ID
	targetID := targetUser.ID

	// Check if users are friends
	areFriends, err := serv.repo.CheckFriends(ctx, requestorID, targetID)
	if err != nil {
		return errors.Wrap(err, response.ErrMsgCheckFriend)
	}

	if areFriends {
		// Delete the subscription
		err = serv.repo.DeleteSubscription(ctx, requestorID, targetID)
		if err != nil {
			return errors.Wrap(err, response.ErrMsgRemoveSubscription)
		}
	}

	// Block the user
	err = serv.repo.BlockUser(ctx, requestorID, targetID)
	if err != nil {
		return errors.Wrap(err, response.ErrMsgBlockUser)
	}

	return nil
}

// GetEligibleRecipients retrieves all email addresses that can receive updates from an email address.
func (serv *friendService) GetEligibleRecipients(ctx context.Context, senderEmail, text string) ([]string, error) {
	senderUser, err := serv.repo.GetUserByEmail(ctx, senderEmail)
	if err != nil {
		return nil, errors.Wrap(err, response.ErrMsgGetUserByEmail)
	}

	senderID := senderUser.ID

	friends, err := serv.repo.GetFriendsList(ctx, senderID)
	if err != nil {
		return nil, errors.Wrap(err, response.ErrMsgGetFriendsList)
	}

	subscribers, err := serv.repo.GetSubscribers(ctx, senderID)
	if err != nil {
		return nil, errors.Wrap(err, response.ErrMsgCheckSubscription)
	}

	mentionedEmails := extractMentionedEmails(text)

	recipientsSet := make(map[string]struct{})
	for _, friendEmail := range friends {
		blocked, err := serv.repo.HasBlockedUpdates(ctx, friendEmail, senderEmail)
		if err != nil {
			return nil, errors.Wrap(err, response.ErrMsgBlockUpdates)
		}
		if !blocked {
			recipientsSet[friendEmail] = struct{}{}
		}
	}

	for _, subscriberEmail := range subscribers {
		blocked, err := serv.repo.HasBlockedUpdates(ctx, subscriberEmail, senderEmail)
		if err != nil {
			return nil, errors.Wrap(err, response.ErrMsgBlockUpdates)
		}
		if !blocked {
			recipientsSet[subscriberEmail] = struct{}{}
		}
	}

	for _, email := range mentionedEmails {
		mentionedUser, err := serv.repo.GetUserByEmail(ctx, email)
		if err == nil {
			blocked, err := serv.repo.HasBlockedUpdates(ctx, mentionedUser.Email, senderEmail)
			if err != nil {
				return nil, errors.Wrap(err, response.ErrMsgBlockUpdates)
			}
			if !blocked {
				recipientsSet[mentionedUser.Email] = struct{}{}
			}
		}
	}

	recipients := make([]string, 0, len(recipientsSet))
	for email := range recipientsSet {
		recipients = append(recipients, email)
	}

	return recipients, nil
}

func extractMentionedEmails(text string) []string {
	words := strings.Fields(text)
	var emails []string
	for _, word := range words {
		if strings.Contains(word, "@") {
			emails = append(emails, word)
		}
	}
	return emails
}
