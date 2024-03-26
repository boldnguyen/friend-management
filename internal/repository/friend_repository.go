package repository

import (
	"context"

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
