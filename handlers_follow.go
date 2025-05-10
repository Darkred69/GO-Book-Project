package main

import (
	"encoding/json"
	"net/http"
	"project_1/internal/database"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

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

func (apiCfg *apiConfig) handlerGetFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	follows, err := apiCfg.DB.GetFollows(r.Context(), user.ID)
	if err != nil {
		responseWithError(w, 500, "Can't get follows")
		return
	}
	responseWithJSON(w, 200, databaseFollowstoFollows(follows))
}

// Using Path Operation for delete
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
