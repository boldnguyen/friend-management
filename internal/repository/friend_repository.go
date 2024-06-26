package repository

import (
	"context"
	"strconv"

	"github.com/boldnguyen/friend-management/internal/models"
	"github.com/boldnguyen/friend-management/internal/pkg/response"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// GetUserByEmail retrieves the user by email from the database.
func (repo friendRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := models.Users(qm.Where("email = ?", email)).One(ctx, repo.DB)
	if err != nil {
		return nil, errors.Wrap(err, response.ErrMsgGetUserByEmail)
	}
	return user, nil
}

// AddFriend creates a friend connection between two users using their user IDs.
func (repo friendRepository) AddFriend(ctx context.Context, userID1, userID2 int) error {
	friendship := models.FriendConnection{
		UserID1: userID1,
		UserID2: userID2,
	}

	err := friendship.Insert(ctx, repo.DB, boil.Infer())
	if err != nil {
		return errors.Wrap(err, response.ErrMsgCreateFriend)
	}
	return nil
}

// CheckFriends checks if two users are already friends.
func (repo friendRepository) CheckFriends(ctx context.Context, userID1, userID2 int) (bool, error) {
	// Check if a friend connection exists between the two users
	exists, err := models.FriendConnections(
		qm.Where("user_id1 = ? AND user_id2 = ?", userID1, userID2),
	).Exists(ctx, repo.DB)
	if err != nil {
		return false, errors.Wrap(err, response.ErrMsgCheckFriend)
	}

	if !exists {
		// Check the reverse direction as well
		exists, err = models.FriendConnections(
			qm.Where("user_id1 = ? AND user_id2 = ?", userID2, userID1),
		).Exists(ctx, repo.DB)
		if err != nil {
			return false, errors.Wrap(err, response.ErrMsgCheckFriend)
		}
	}

	return exists, nil
}

// GetFriendsList retrieves the list of friends for a given user ID.
func (repo friendRepository) GetFriendsList(ctx context.Context, userID int) ([]string, error) {
	friendConnections, err := models.FriendConnections(
		qm.Select("user_id1", "user_id2"),
		qm.Where("user_id1 = ? OR user_id2 = ?", userID, userID),
	).All(ctx, repo.DB)
	if err != nil {
		return nil, errors.Wrap(err, response.ErrMsgGetFriendsList)
	}

	var friends []string
	for _, fc := range friendConnections {
		var friendID int
		if fc.UserID1 == userID {
			friendID = fc.UserID2
		} else {
			friendID = fc.UserID1
		}

		user, err := models.Users(qm.Where("id = ?", friendID)).One(ctx, repo.DB)
		if err != nil {
			return nil, errors.Wrap(err, response.ErrMsgGetFriendsList)
		}
		friends = append(friends, user.Email)
	}

	return friends, nil
}

// GetCommonFriends retrieves the list of common friends between two user IDs.
func (repo friendRepository) GetCommonFriends(ctx context.Context, userID1, userID2 int) ([]string, error) {
	// Get friends list for both users
	friends, err := models.FriendConnections(
		qm.Select("user_id2"),
		qm.Where("user_id1 IN (?, ?)", userID1, userID2), // Select friends of both users
	).All(ctx, repo.DB)
	if err != nil {
		return nil, errors.Wrap(err, response.ErrMsgGetCommonFriends)
	}

	// Find common friends
	commonFriends := make(map[int]int) // Using map for efficient lookup
	var common []string
	for _, friend := range friends {
		commonFriends[friend.UserID2]++ // Count occurrences of each friend
	}

	// Retrieve users corresponding to common friends
	for friendID, count := range commonFriends {
		if count == 2 { // If a friend is found for both users
			user, err := models.Users(qm.Where("id = ?", friendID)).One(ctx, repo.DB)
			if err != nil {
				return nil, errors.Wrap(err, response.ErrMsgGetCommonFriends)
			}
			common = append(common, user.Email)
		}
	}

	return common, nil
}

// CheckSubscription checks if a subscription already exists between requestor and target.
func (repo *friendRepository) CheckSubscription(ctx context.Context, requestor, target string) (bool, error) {
	exists, err := models.Subscriptions(
		qm.Where("requestor = ? AND target = ?", requestor, target),
	).Exists(ctx, repo.DB)
	if err != nil {
		return false, errors.Wrap(err, response.ErrMsgCheckSubscription)
	}
	return exists, nil
}

// SubscribeUpdates subscribes the requestor to updates from the target.
func (repo *friendRepository) SubscribeUpdates(ctx context.Context, requestor, target string) error {
	// Check if the subscription already exists
	exists, err := repo.CheckSubscription(ctx, requestor, target)
	if err != nil {
		return errors.Wrap(err, response.ErrMsgSubscribeUpdates)
	}
	if exists {
		return errors.New(response.ErrMsgSubscriptionAlreadyExists)
	}

	subscription := models.Subscription{
		Requestor: requestor,
		Target:    target,
	}
	err = subscription.Insert(ctx, repo.DB, boil.Infer())
	if err != nil {
		return errors.Wrap(err, response.ErrMsgSubscribeUpdates)
	}
	return nil
}

// DeleteSubscription deletes a subscription from the database.
func (repo *friendRepository) DeleteSubscription(ctx context.Context, requestorID, targetID int) error {
	_, err := models.Subscriptions(
		qm.Where("requestor = ? AND target = ?", requestorID, targetID),
	).DeleteAll(ctx, repo.DB)
	if err != nil {
		return errors.Wrap(err, response.ErrMsgBlockUpdates)
	}
	return nil
}

// BlockUser inserts a block entry into the database.
func (repo *friendRepository) BlockUser(ctx context.Context, requestorID, targetID int) error {
	block := models.Block{
		Requestor: requestorID,
		Target:    targetID,
	}
	err := block.Insert(ctx, repo.DB, boil.Infer())
	if err != nil {
		return errors.Wrap(err, response.ErrMsgBlockUpdates)
	}
	return nil
}

// GetSubscribers retrieves the list of subscribers for a given user ID.
func (repo *friendRepository) GetSubscribers(ctx context.Context, userID int) ([]string, error) {
	userIDStr := strconv.Itoa(userID) // Chuyển đổi userID từ int sang string

	query := `
        SELECT u.email
        FROM users u
        JOIN subscriptions s ON u.id = CAST(s.requestor AS INTEGER)
        WHERE s.target = $1;
    `

	rows, err := repo.DB.QueryContext(ctx, query, userIDStr)
	if err != nil {
		return nil, errors.Wrap(err, response.ErrMsgCheckSubscription)
	}
	defer rows.Close()

	var subscribers []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, errors.Wrap(err, response.ErrMsgCheckSubscription)
		}
		subscribers = append(subscribers, email)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, response.ErrMsgCheckSubscription)
	}

	return subscribers, nil
}

// HasBlockedUpdates checks if the target user has blocked updates from the sender.
func (repo *friendRepository) HasBlockedUpdates(ctx context.Context, targetEmail, senderEmail string) (bool, error) {
	targetUser, err := repo.GetUserByEmail(ctx, targetEmail)
	if err != nil {
		return false, errors.Wrap(err, response.ErrMsgGetUserByEmail)
	}

	senderUser, err := repo.GetUserByEmail(ctx, senderEmail)
	if err != nil {
		return false, errors.Wrap(err, response.ErrMsgGetUserByEmail)
	}

	exists, err := models.Blocks(
		qm.Where("requestor = ? AND target = ?", targetUser.ID, senderUser.ID),
	).Exists(ctx, repo.DB)
	if err != nil {
		return false, errors.Wrap(err, response.ErrMsgBlockUpdates)
	}

	return exists, nil
}
