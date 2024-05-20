package handler

import (
	"encoding/json"
	"net/http"

	"github.com/boldnguyen/friend-management/internal/pkg/response"
	"github.com/boldnguyen/friend-management/internal/service"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// CreateFriendConnectionRequest defines the structure of the request for creating a friend connection.
type CreateFriendConnectionRequest struct {
	Friends []string `json:"friends" validate:"required,min=2,max=2,dive,email"`
}

// CommonFriendsRequest defines the structure of the request for retrieving common friends.
type CommonFriendsRequest struct {
	Friends []string `json:"friends" validate:"required,min=2,max=2,dive,email"`
}

// SubscribeRequest defines the structure of the request for subscribing to updates.
type SubscribeRequest struct {
	Requestor string `json:"requestor" validate:"required,email"`
	Target    string `json:"target" validate:"required,email"`
}
type BlockUpdatesRequest struct {
	Requestor string `json:"requestor" validate:"required,email"`
	Target    string `json:"target" validate:"required,email"`
}

type getRecipientsRequest struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
}

type getRecipientsResponse struct {
	Success    bool     `json:"success"`
	Recipients []string `json:"recipients"`
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
		response.RespondSuccess(ctx, w, nil)
	}
}

// FriendListHandler creates a new HTTP handler for retrieving the friends list.
func FriendListHandler(friendService service.FriendService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email string `json:"email" validate:"required,email"`
		}
		ctx := r.Context()

		// Decode the JSON data from the request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.RespondErr(ctx, w, http.StatusBadRequest, response.ErrMsgDecodeRequest+": "+err.Error())
			return
		}

		// Call the friend service to retrieve the friends list
		friends, err := friendService.GetFriendsList(ctx, req.Email)
		if err != nil {
			response.RespondErr(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		// Respond with success and friends list
		response.RespondSuccess(ctx, w, map[string]interface{}{
			"success": true,
			"friends": friends,
			"count":   len(friends),
		})
	}
}

// CommonFriendsHandler creates a new HTTP handler for retrieving common friends list.
func CommonFriendsHandler(friendService service.FriendService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CommonFriendsRequest
		ctx := r.Context()

		// Decode the JSON data from request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.RespondErr(ctx, w, http.StatusBadRequest, response.ErrMsgDecodeRequest)
			return
		}
		// Validate input
		if len(req.Friends) != 2 {
			response.RespondErr(ctx, w, http.StatusBadRequest, response.ErrMsgInvalidRequest)
			return
		}

		// Call the friend service to retrieve common friends list
		commonFriends, err := friendService.GetCommonFriends(ctx, req.Friends[0], req.Friends[1])
		if err != nil {
			response.RespondErr(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		// Respond with success and common friends list
		response.RespondSuccess(ctx, w, map[string]interface{}{
			"success": true,
			"friends": commonFriends,
			"count":   len(commonFriends),
		})
	}
}

// SubscribeHandler creates a new HTTP handler for subscribing to updates.
func SubscribeHandler(friendService service.FriendService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SubscribeRequest
		ctx := r.Context()

		// Decode the JSON data from the request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.RespondErr(ctx, w, http.StatusBadRequest, response.ErrMsgDecodeRequest)
			return
		}

		// Validate the request
		if err := validate.Struct(req); err != nil {
			response.RespondErr(ctx, w, http.StatusBadRequest, response.ErrMsgInvalidRequest)
		}

		// Call the friend service to subscribe to updates
		err := friendService.SubscribeUpdates(ctx, req.Requestor, req.Target)
		if err != nil {
			response.RespondErr(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		// Respond with success
		response.RespondSuccess(ctx, w, nil)
	}
}

// BlockUpdatesHandler creates a new HTTP handler for blocking updates from an email address.
func BlockUpdatesHandler(friendService service.FriendService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BlockUpdatesRequest
		ctx := r.Context()

		// Decode the JSON data from the request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.RespondErr(ctx, w, http.StatusBadRequest, response.ErrMsgDecodeRequest)
			return
		}

		// Call the friend service to block updates
		err := friendService.BlockUpdates(ctx, req.Requestor, req.Target)
		if err != nil {
			response.RespondErr(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}

		// Respond with success
		response.RespondSuccess(ctx, w, map[string]bool{"success": true})
	}
}

func GetRecipientsHandler(friendService service.FriendService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req getRecipientsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		recipients, err := friendService.GetEligibleRecipients(r.Context(), req.Sender, req.Text)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res := getRecipientsResponse{
			Success:    true,
			Recipients: recipients,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
