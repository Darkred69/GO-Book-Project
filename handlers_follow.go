package main

import (
	"encoding/json"
	"net/http"
	"project_1/internal/database"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// handlerFollowFeed allows the user to follow a feed
// @Summary      Follow a feed
// @Description  Create a follow relationship for a specific feed
// @Tags         follow
// @Accept       json
// @Produce      json
// @Param        follow  body      map[string]string  true  "Feed ID to follow"
// @Success      201     {object}  Follow
// @Failure      400     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Failure      409     {object}  map[string]string
// @Router       /v3/follow [post]
// @Security     BearerAuth
func (apiCfg *apiConfig) handlerFollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	var p params
	err := decoder.Decode(&p)
	if err != nil {
		responseWithError(w, 400, "Invalid request payload")
		return
	}
	follow, err := apiCfg.DB.CreateFollow(r.Context(), database.CreateFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: p.FeedID,
	})
	if err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			responseWithError(w, http.StatusConflict, "Feed already followed")
			return
		}
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			responseWithError(w, http.StatusNotFound, "Feed not found")
			return
		}
		responseWithError(w, http.StatusBadRequest, "Can't follow feed")
		return
	}

	responseWithJSON(w, http.StatusCreated, databaseFollowtoFollow(follow))
}

// handlerGetFollows returns all feeds followed by the user
// @Summary      Get followed feeds
// @Description  Retrieve a list of feeds the authenticated user is following
// @Tags         follow
// @Produce      json
// @Success      200  {array}   Follow
// @Failure      500  {object}  map[string]string
// @Router       /v3/follow [get]
// @Security     BearerAuth
func (apiCfg *apiConfig) handlerGetFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	follows, err := apiCfg.DB.GetFollows(r.Context(), user.ID)
	if err != nil {
		responseWithError(w, 500, "Can't get follows")
		return
	}
	responseWithJSON(w, 200, databaseFollowstoFollows(follows))
}

// handlerUnfollow removes a feed from the user's followed list
// @Summary      Unfollow a feed
// @Description  Unfollow a feed by ID
// @Tags         follow
// @Produce      json
// @Param        feed_id  path      string  true  "Feed ID to unfollow"
// @Success      204      {object}  map[string]string "status": "No Content"
// @Failure      400      {object}  map[string]string "error": "Invalid feed ID"
// @Failure      404      {object}  map[string]string "error": "Feed not found"
// @Failure      500      {object}  map[string]string "error": "Internal server error"
// @Router       /v3/follow/{feed_id} [delete]
// @Security     BearerAuth
func (apiCfg *apiConfig) handlerUnfollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feed := chi.URLParam(r, "feed_id")
	feedID, err := uuid.Parse(feed)
	if err != nil {
		responseWithError(w, 400, "Invalid feed id")
		return
	}

	_, err = apiCfg.DB.GetFeed(r.Context(), feedID)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "Feed not found")
		return
	}

	_, err = apiCfg.DB.GetFollowsByFeedID(r.Context(), database.GetFollowsByFeedIDParams{
		FeedID: feedID,
		UserID: user.ID,
	})
	if err != nil {
		responseWithError(w, http.StatusNotFound, "Feed not followed")
		return
	}

	err = apiCfg.DB.Unfollow(r.Context(), database.UnfollowParams{
		UserID: user.ID,
		FeedID: feedID,
	})
	if err != nil {
		responseWithError(w, 500, "Can't unfollow feed")
		return
	}
	responseWithJSON(w, 204, map[string]string{"status": "No Content"})
}
