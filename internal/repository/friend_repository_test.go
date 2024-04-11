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

// LoadSqlTestFile reads the SQL test file and executes it on the provided test database connection.
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
