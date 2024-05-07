package repository

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/boldnguyen/friend-management/internal/models"
	"github.com/boldnguyen/friend-management/internal/pkg/db"
	"github.com/boldnguyen/friend-management/internal/pkg/response"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// LoadSqlTestFile reads the SQL test file and executes it on the provided test database connection.
func LoadSqlTestFile(t *testing.T, tx *sql.DB, sqlFile string) {
	b, err := os.ReadFile(sqlFile)
	require.NoError(t, err)

	_, err = tx.Exec(string(b))
	require.NoError(t, err)
}

// LoadSqlTestFileWithTx reads the SQL test file and executes it on the provided transaction.
func LoadSqlTestFileWithTx(t *testing.T, tx *sql.Tx, sqlFile string) {
	b, err := os.ReadFile(sqlFile)
	require.NoError(t, err)

	_, err = tx.Exec(string(b))
	require.NoError(t, err)
}

// TestFriendRepository_GetUserByEmail tests the GetUserByEmail method of the friendRepository.
func TestFriendRepository_GetUserByEmail(t *testing.T) {
	// Your test cases
	tcs := map[string]struct {
		email       string
		expUser     *models.User
		expectedErr string
	}{
		"existing_user": {
			email: "john@example.com",
			expUser: &models.User{
				ID:    1,
				Name:  "John Doe",
				Email: "john@example.com",
			},
			expectedErr: "",
		},
		"non_existing_user": {
			email:       "nonexistent@example.com",
			expUser:     nil,
			expectedErr: response.ErrMsgGetUserByEmail,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			ctx := context.Background()

			// Open test database connection
			dbTest, err := db.ConnectDB("postgres://friend-management:1234@localhost:5432/friend-management?sslmode=disable")
			require.NoError(t, err)
			defer dbTest.Close()

			// Load test data
			LoadSqlTestFile(t, dbTest, "../testdata/friends.sql")

			// Initialize repository with the mock
			mockRepo := &MockRepo{}

			// When GetUserByEmail is called
			mockRepo.On("GetUserByEmail", ctx, tc.email).Return(tc.expUser, nil)

			// Initialize repository with the mocked repository
			repo := friendRepository{DB: dbTest}

			// When
			user, err := repo.GetUserByEmail(ctx, tc.email)

			// Then
			if tc.expUser == nil {
				assert.Nil(t, user)
			} else {
				require.NotNil(t, user) // Ensure user is not nil
				assert.Equal(t, tc.expUser.ID, user.ID)
				assert.Equal(t, tc.expUser.Name, user.Name)
				assert.Equal(t, tc.expUser.Email, user.Email)
			}

			if tc.expectedErr != "" {
				assert.Contains(t, err.Error(), tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestFriendRepository_AddFriend tests the AddFriend method of the friendRepository.
func TestFriendRepository_AddFriend(t *testing.T) {
	tcs := map[string]struct {
		userID1     int
		userID2     int
		expectedErr string
	}{
		"add_friend_success": {
			userID1:     1,
			userID2:     2,
			expectedErr: "",
		},
		"add_friend_failure": {
			userID1:     1,
			userID2:     3, // Assuming this user does not exist
			expectedErr: response.ErrMsgCreateFriend,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			ctx := context.Background()

			// Initialize mock repository
			mockRepo := &MockRepo{}

			// When AddFriend is called
			mockRepo.On("AddFriend", ctx, tc.userID1, tc.userID2).Return(nil)

			// When
			err := mockRepo.AddFriend(ctx, tc.userID1, tc.userID2)

			// Then
			if tc.expectedErr != "" {
				assert.Contains(t, err.Error(), tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestFriendRepository_CheckFriends tests the CheckFriends method of the friendRepository.
func TestFriendRepository_CheckFriends(t *testing.T) {
	tcs := map[string]struct {
		userID1     int
		userID2     int
		areFriends  bool // Expected result
		expectedErr string
	}{
		"existing_friends": {
			userID1:    1,
			userID2:    2,
			areFriends: true, // John and Jane are friends
		},
		"non_existing_friends": {
			userID1:    1,
			userID2:    3, // Assuming user 3 is not friends with user 1
			areFriends: false,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			ctx := context.Background()

			// Initialize repository with the mock
			mockRepo := &MockRepo{}

			// When CheckFriends is called, we expect it to be called with the given context and user IDs
			mockRepo.On("CheckFriends", ctx, tc.userID1, tc.userID2).Return(tc.areFriends, nil)

			// When
			areFriends, err := mockRepo.CheckFriends(ctx, tc.userID1, tc.userID2)

			// Then
			if tc.expectedErr != "" {
				assert.Contains(t, err.Error(), tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.areFriends, areFriends)
			}
		})
	}
}

// TestFriendRepository_GetFriendsList tests the GetFriendsList method of the friendRepository.
func TestFriendRepository_GetFriendsList(t *testing.T) {
	// Your test cases
	tcs := map[string]struct {
		userID      int
		expFriends  []string
		expectedErr string
	}{
		"existing_user_with_friends": {
			userID:      1,
			expFriends:  []string{"jane@example.com", "bob@example.com"}, // Expect John to be friends with Jane
			expectedErr: "",
		},
		"existing_user_without_friends": {
			userID:      3,
			expFriends:  []string{}, // Expect Alice to have no friends
			expectedErr: "",
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			ctx := context.Background()

			// Open test database connection
			dbTest, err := db.ConnectDB("postgres://friend-management:1234@localhost:5432/friend-management?sslmode=disable")
			require.NoError(t, err)
			defer dbTest.Close()

			// Load test data
			LoadSqlTestFile(t, dbTest, "../testdata/friends.sql")

			// Initialize mock repository
			mockRepo := &MockRepo{}

			// Mock the GetFriendsList method
			mockRepo.On("GetFriendsList", ctx, tc.userID).Return(tc.expFriends, nil)

			// Initialize repository with the mocked repository
			repo := friendRepository{DB: dbTest}

			// When GetFriendsList is called
			friends, err := repo.GetFriendsList(ctx, tc.userID)

			// Then
			if tc.expectedErr != "" {
				assert.Contains(t, err.Error(), tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tc.expFriends, friends)
			}
		})
	}
}

// TestFriendRepository_GetCommonFriends tests the GetCommonFriends method of the friendRepository.
func TestFriendRepository_GetCommonFriends(t *testing.T) {
	// Your test cases
	tcs := map[string]struct {
		userID1     int
		userID2     int
		expCommon   []string
		expectedErr string
	}{
		"existing_common_friends": {
			userID1:     1,
			userID2:     4,
			expCommon:   []string{"jane@example.com"}, // Expect John and Bob to have a common friend
			expectedErr: "",
		},
		"no_common_friends": {
			userID1:     1,
			userID2:     3,
			expCommon:   []string{}, // Expect John and Alice to have no common friends
			expectedErr: "",
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			ctx := context.Background()

			// Open test database connection
			dbTest, err := db.ConnectDB("postgres://friend-management:1234@localhost:5432/friend-management?sslmode=disable")
			require.NoError(t, err)
			defer dbTest.Close()

			// Load test data
			LoadSqlTestFile(t, dbTest, "../testdata/friends.sql")

			// Initialize mock repository
			mockRepo := &MockRepo{}

			// Mock the GetCommonFriends method
			mockRepo.On("GetCommonFriends", ctx, tc.userID1, tc.userID2).Return(tc.expCommon, nil)

			// Initialize repository with the mocked repository
			repo := friendRepository{DB: dbTest}

			// When GetCommonFriends is called
			common, err := repo.GetCommonFriends(ctx, tc.userID1, tc.userID2)

			// Then
			if tc.expectedErr != "" {
				assert.Contains(t, err.Error(), tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tc.expCommon, common)
			}
		})
	}
}

// TestFriendRepository_CheckSubscription tests the CheckSubscription method of the friendRepository.
func TestFriendRepository_CheckSubscription(t *testing.T) {
	// Your test cases
	tcs := map[string]struct {
		requestor    string
		target       string
		subscription bool  // Expected result
		expectedErr  error // Expected error
	}{
		"existing_subscription": {
			requestor:    "user1@example.com",
			target:       "user2@example.com",
			subscription: true, // Expect a subscription to exist between user1 and user2
			expectedErr:  nil,
		},
		"non_existing_subscription": {
			requestor:    "user1@example.com",
			target:       "user3@example.com", // Assuming user3 is not subscribed by user1
			subscription: false,
			expectedErr:  nil,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			ctx := context.Background()

			// Open test database connection
			dbTest, err := db.ConnectDB("postgres://friend-management:1234@localhost:5432/friend-management?sslmode=disable")
			require.NoError(t, err)
			defer dbTest.Close()

			// Start a transaction
			tx, err := dbTest.Begin()
			require.NoError(t, err)

			// Load test data using the transaction
			LoadSqlTestFileWithTx(t, tx, "../testdata/subscriptions.sql")

			// Initialize mock repository
			mockRepo := &MockRepo{}

			// When CheckSubscription is called, we expect it to be called with the given context, requestor, and target
			mockRepo.On("CheckSubscription", ctx, tc.requestor, tc.target).Return(tc.subscription, tc.expectedErr)

			// Initialize repository with the mock and transaction
			repo := friendRepository{DB: tx}

			// When
			subscription, err := repo.CheckSubscription(ctx, tc.requestor, tc.target)

			// Then
			assert.Equal(t, tc.expectedErr, err) // Check for expected error
			assert.Equal(t, tc.subscription, subscription)

			// Rollback the transaction
			err = tx.Rollback()
			require.NoError(t, err)
		})
	}
}

// TestFriendRepository_SubscribeUpdates tests the SubscribeUpdates method of the friendRepository.
func TestFriendRepository_SubscribeUpdates(t *testing.T) {
	// Your test cases
	tcs := map[string]struct {
		requestor   string
		target      string
		expectedErr error // Expected error
	}{
		"subscribe_success": {
			requestor:   "user1@example.com",
			target:      "user3@example.com",
			expectedErr: nil,
		},
		"subscribe_failure": {
			requestor:   "user1@example.com",
			target:      "user2@example.com", // Assuming that this subscription already exists
			expectedErr: errors.New(response.ErrMsgSubscriptionAlreadyExists),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			ctx := context.Background()

			// Open test database connection
			dbTest, err := db.ConnectDB("postgres://friend-management:1234@localhost:5432/friend-management?sslmode=disable")
			require.NoError(t, err)
			defer dbTest.Close()

			// Start transaction
			tx, err := dbTest.Begin()
			require.NoError(t, err)

			// Load test data
			LoadSqlTestFileWithTx(t, tx, "../testdata/subscriptions.sql")

			// Initialize mock repository
			mockRepo := &MockRepo{}

			// Mock SubscribeUpdates method behavior
			mockRepo.On("SubscribeUpdates", ctx, tc.requestor, tc.target).Return(tc.expectedErr)

			// Initialize repository with the mock and transaction
			repo := friendRepository{DB: tx}

			// When
			err = repo.SubscribeUpdates(ctx, tc.requestor, tc.target)

			// Rollback transaction
			require.NoError(t, tx.Rollback())

			// Then
			if tc.expectedErr != nil {
				require.EqualError(t, err, tc.expectedErr.Error(), "Error message mismatch")
			} else {
				require.NoError(t, err, "Unexpected error occurred")
			}

			// Log expected and actual error messages
			t.Logf("Expected error: %q", tc.expectedErr)
			t.Logf("Actual error: %q", err)
		})
	}
}

// TestFriendRepository_BlockUpdates tests the BlockUpdates method of the friendRepository.
func TestFriendRepository_BlockUpdates(t *testing.T) {
	tcs := map[string]struct {
		requestor   string
		target      string
		expectedErr error
	}{
		"block_success": {
			requestor:   "user1@example.com",
			target:      "user2@example.com",
			expectedErr: nil,
		},
		"block_failure_subscription_does_not_exist": {
			requestor:   "user1@example.com",
			target:      "user4@example.com",
			expectedErr: errors.New(response.ErrMsgSubscriptionDoesNotExist),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given
			ctx := context.Background()

			// Open test database connection
			dbTest, err := db.ConnectDB("postgres://friend-management:1234@localhost:5432/friend-management?sslmode=disable")
			require.NoError(t, err)
			defer dbTest.Close()

			// Start transaction
			tx, err := dbTest.Begin()
			require.NoError(t, err)

			// Load test data
			LoadSqlTestFileWithTx(t, tx, "../testdata/subscriptions.sql")

			// Initialize mock repository
			mockRepo := &MockRepo{}

			// Mock BlockUpdates method behavior
			mockRepo.On("BlockUpdates", ctx, tc.requestor, tc.target).Return(tc.expectedErr)

			// Initialize repository with the mock and transaction
			repo := friendRepository{DB: tx}

			// When
			err = repo.BlockUpdates(ctx, tc.requestor, tc.target)

			// Rollback transaction
			require.NoError(t, tx.Rollback())

			// Then
			if tc.expectedErr != nil {
				require.EqualError(t, err, tc.expectedErr.Error(), "Error message mismatch")
			} else {
				require.NoError(t, err, "Unexpected error occured")
			}

			// Log expected and actual error messages
			t.Logf("Expected error: %q", tc.expectedErr)
			t.Logf("Acutal error: %q", err)

		})
	}
}
