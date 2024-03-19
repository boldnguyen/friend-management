package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/boldnguyen/friend-management/internal/pkg/response"
	"github.com/boldnguyen/friend-management/internal/service"
)

// FriendHandler defines the HTTP handler interface for managing friend connections.
type FriendHandler interface {
	// AddFriend handles the HTTP request to create a friend connection between two users.
	AddFriend(w http.ResponseWriter, r *http.Request)
}

// friendHandler is the concrete implementation of FriendHandler.
type friendHandler struct {
	service service.FriendService
}

// NewFriendHandler creates a new instance of FriendHandler with the provided service.
func NewFriendHandler(service service.FriendService) FriendHandler {
	return &friendHandler{service: service}
}

// AddFriend handles the HTTP request to create a friend connection between two users.
func (h *friendHandler) AddFriend(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode the JSON request body
	var friendRequest struct {
		Friends []string `json:"friends"`
	}
	if err := json.NewDecoder(r.Body).Decode(&friendRequest); err != nil {
		response.RespondErr(ctx, w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Call the FriendService to add the friend connection
	err := h.service.AddFriend(ctx, friendRequest.Friends)
	if err != nil {
		// Check if the error is due to the users already being friends
		if strings.Contains(err.Error(), "already friends") {
			response.RespondErr(ctx, w, http.StatusBadRequest, "They are already friends")
			return
		}

		// Respond with an error message for other errors
		response.RespondErr(ctx, w, http.StatusInternalServerError, "Failed to add friend")
		return
	}

	// Respond with success if the friend connection is added successfully
	response.RespondSuccess(ctx, w, nil)
}
