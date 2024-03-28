package repository

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/boldnguyen/friend-management/internal/models"
	"github.com/boldnguyen/friend-management/internal/pkg/db"
	"github.com/boldnguyen/friend-management/internal/pkg/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// LoadSqlTestFile loads SQL data from a file into the provided test database connection.
func LoadSqlTestFile(t *testing.T, tx *sql.DB, sqlFile string) {
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
			LoadSqlTestFile(t, dbTest, "C:/Users/nguyen.nguyen/Desktop/friend-management/testdata/friends.sql")

			// Initialize your repository with the test database
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

			// Open test database connection
			dbTest, err := db.ConnectDB("postgres://friend-management:1234@localhost:5432/friend-management?sslmode=disable")
			require.NoError(t, err)
			defer dbTest.Close()

			// Clear the friend_connections table
			_, err = dbTest.Exec("DELETE FROM friend_connections")
			require.NoError(t, err)

			// Load test data using the original database connection
			LoadSqlTestFile(t, dbTest, "C:/Users/nguyen.nguyen/Desktop/friend-management/testdata/friends.sql")

			// Initialize your repository with the transaction
			repo := friendRepository{DB: dbTest}
			// Check if the friend connection already exists
			exists, err := repo.CheckFriends(ctx, tc.userID1, tc.userID2)
			require.NoError(t, err)

			if exists {
				// Skip the insert operation or handle as needed
				t.Skip("Friend connection already exists")
			}

			// When
			err = repo.AddFriend(ctx, tc.userID1, tc.userID2)

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
	// Your test cases
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
			// Open test database connection
			dbTest, err := db.ConnectDB("postgres://friend-management:1234@localhost:5432/friend-management?sslmode=disable")
			require.NoError(t, err)
			defer dbTest.Close()

			// Clear the friend_connections table
			_, err = dbTest.Exec("DELETE FROM friend_connections")
			require.NoError(t, err)

			// Load test data using the original database connection
			LoadSqlTestFile(t, dbTest, "C:/Users/nguyen.nguyen/Desktop/friend-management/testdata/friends.sql")

			// Initialize your repository with the transaction
			repo := friendRepository{DB: dbTest}

			// When
			areFriends, err := repo.CheckFriends(ctx, tc.userID1, tc.userID2)

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
