package handler

import (
	"encoding/json"
	"net/http"

	"github.com/boldnguyen/friend-management/internal/pkg/response"
	"github.com/boldnguyen/friend-management/internal/service"
)

// CreateFriendConnectionRequest defines the structure of the request for creating a friend connection.
type CreateFriendConnectionRequest struct {
	Friends []string `json:"friends" validate:"required,min=2,max=2,dive,email"`
}

// NewHandler creates a new HTTP handler for creating a friend connection.
func NewHandler(friendService service.FriendService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateFriendConnectionRequest
		ctx := r.Context()

		// Decode the JSON data from the request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.RespondErr(ctx, w, http.StatusBadRequest, response.ErrMsgDecodeRequest+": "+err.Error())
			return
		}

		// Validate input
		if len(req.Friends) != 2 {
			response.RespondErr(ctx, w, http.StatusBadRequest, response.ErrMsgInvalidRequest)
			return
		}

		// Call the friend service to create the friend connection
		err := friendService.CreateFriend(ctx, req.Friends[0], req.Friends[1])
		if err != nil {
			response.RespondErr(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		// Respond with success
		response.RespondSuccess(ctx, w, map[string]bool{"success": true})
	}
}
