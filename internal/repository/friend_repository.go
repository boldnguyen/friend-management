package repository

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

// FriendRepository defines the interface for interacting with friend-related data in the database.
type FriendRepository interface {
	// AddFriend creates a friend connection between two users in the database.
	AddFriend(ctx context.Context, userID1, userID2 int) error

	// GetUserIDByEmail retrieves the user ID associated with a given email address from the database.
	GetUserIDByEmail(ctx context.Context, email string) (int, error)

	// CreateUser inserts a new user with a default name and the given email address into the database.
	CreateUser(ctx context.Context, email string) (int, error)
}

// friendRepository is the concrete implementation of FriendRepository.
type friendRepository struct {
	db *sql.DB
}

// NewFriendRepository creates a new instance of FriendRepository with the provided database connection.
func NewFriendRepository(db *sql.DB) FriendRepository {
	return &friendRepository{db: db}
}

// AddFriend creates a friend connection between two users in the database.
func (r *friendRepository) AddFriend(ctx context.Context, userID1, userID2 int) error {
	// Begin a database transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer func() {
		if err := recover(); err != nil {
			_ = tx.Rollback()
		}
	}()

	// Check if the friendship already exists in the database
	var count int
	err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM friend_connections WHERE (user_id1 = $1 AND user_id2 = $2) OR (user_id1 = $2 AND user_id2 = $1)", userID1, userID2).Scan(&count)
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "failed to check existing friendship")
	}

	if count > 0 {
		_ = tx.Rollback()
		return errors.New("friendship already exists in the database")
	}

	// Insert the new friend connection
	_, err = tx.ExecContext(ctx, "INSERT INTO friend_connections (user_id1, user_id2) VALUES ($1, $2)", userID1, userID2)
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "failed to insert friend connection")
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// GetUserIDByEmail retrieves the user ID associated with a given email address from the database.
func (r *friendRepository) GetUserIDByEmail(ctx context.Context, email string) (int, error) {
	var userID int
	err := r.db.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", email).Scan(&userID)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get user ID by email")
	}
	return userID, nil
}

// CreateUser inserts a new user with a default name and the given email address into the database.
func (r *friendRepository) CreateUser(ctx context.Context, email string) (int, error) {
	var userID int
	err := r.db.QueryRowContext(ctx, "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", "DefaultName", email).Scan(&userID)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create user")
	}
	return userID, nil
}