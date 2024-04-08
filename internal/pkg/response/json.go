package response

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/httplog"
)

// Error messages
const (
	ErrMsgDecodeRequest  = "failed to decode JSON data from request body"
	ErrMsgInvalidRequest = "Exactly two email addresses are required to create a friend connection"
	ErrMsgCreateFriend   = "failed to create friend connection"
	ErrMsgGetUserByEmail = "failed to get user by email"
	ErrMsgUserNotFound   = "user not found"
	ErrMsgCheckFriend    = "failed to check if users are already friends"
	ErrMsgAlreadyFriends = "They are already friends"
	ErrMsgGetFriendsList = "failed to get friends list"
)

// RespondSuccess responds basic success response
func RespondSuccess(ctx context.Context, w http.ResponseWriter, data interface{}) error {
	log := httplog.LogEntry(ctx)

	w.Header().Set("Content-Type", "application/json")

	resp := map[string]interface{}{
		"success": true,
	}
	if data != nil {
		resp["data"] = data
	}
	respByte, err := json.Marshal(resp)
	if err != nil {
		log.Error().Msgf("Failed to marshal success response, err: %s", err)
		if err := RespondErr(ctx, w, http.StatusInternalServerError, err.Error()); err != nil {
			return err
		}
		return nil
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respByte)
	return nil
}

// RespondErr responds basic error response
func RespondErr(
	ctx context.Context,
	w http.ResponseWriter,
	code int,
	msg string,
) error {
	log := httplog.LogEntry(ctx)

	w.Header().Set("Content-Type", "application/json")

	respByte, err := json.Marshal(map[string]interface{}{
		"success":       false,
		"error_message": msg,
	})
	if err != nil {
		log.Error().Msgf("Failed to marshal error response, err: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(code)
	w.Write(respByte)
	return nil
}
