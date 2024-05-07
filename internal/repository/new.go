package repository

import (
	"context"
	"database/sql"

	"github.com/boldnguyen/friend-management/internal/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// FriendRepository provides methods for interacting with friend data in the database.
type FriendRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	AddFriend(ctx context.Context, userID1, userID2 int) error
	CheckFriends(ctx context.Context, userID1, userID2 int) (bool, error)
	GetFriendsList(ctx context.Context, userID int) ([]string, error)
	GetCommonFriends(ctx context.Context, userID1, userID2 int) ([]string, error)
	CheckSubscription(ctx context.Context, requestor, target string) (bool, error) // New method
	SubscribeUpdates(ctx context.Context, requestor, target string) error
	BlockUpdates(ctx context.Context, requestor, target string) error
}

// friendRepository implements the FriendRepository interface.
type friendRepository struct {
	DB boil.ContextExecutor
}

// NewFriendRepository creates a new instance of FriendRepository.
func NewFriendRepository(db *sql.DB) FriendRepository {
	return &friendRepository{DB: db}
}
