package main

import (
	"encoding/json"
	"net/http"
	"project_1/internal/database"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// Create Feed
func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	decoder := json.NewDecoder(r.Body)
	var p FeedInput
	err := decoder.Decode(&p)
	if err != nil {
		responseWithError(w, 400, "Invalid request payload")
		return
	}

	// Function in Utils Check for valid url
	if !isValidURL(p.URL) {
		responseWithError(w, http.StatusBadRequest, "Invalid URL")
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   p.Name,
		Url:    p.URL,
		UserID: user.ID,
	})
	if err != nil {
		responseWithError(w, http.StatusConflict, "Feed exist")
		return
	}

	responseWithJSON(w, http.StatusCreated, databaseFeedtoFeed(feed))
}

// Get Feeds
func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := apiCfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		responseWithError(w, 500, "Can't get feeds")
		return
	}

	responseWithJSON(w, 200, databaseFeedstoFeeds(feeds))
}

// Update Feed
func (apiCfg *apiConfig) handlerUpdateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	decoder := json.NewDecoder(r.Body)
	var p FeedInput
	err := decoder.Decode(&p)
	if err != nil {
		responseWithError(w, 400, "Invalid request payload")
		return
	}
	feed_id := chi.URLParam(r, "feed_id")
	feedID, err := uuid.Parse(feed_id)
	if err != nil {
		responseWithError(w, 400, "Invalid feed id")
		return
	}
	feed, err := apiCfg.DB.GetFeed(r.Context(), feedID)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "Feed don't exsist")
		return
	}

	if feed.UserID != user.ID {
		responseWithError(w, 403, "Forbidden")
		return
	}

	// Function in Utils Check for valid url
	if !isValidURL(p.URL) {
		responseWithError(w, http.StatusBadRequest, "Invalid URL")
		return
	}

	feed, err = apiCfg.DB.UpdateFeed(r.Context(), database.UpdateFeedParams{
		Name:   p.Name,
		Url:    p.URL,
		UserID: user.ID,
		ID:     feedID,
	})

	if err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			responseWithError(w, http.StatusConflict, "Duplicate feed exist")
			return
		}
		responseWithError(w, 500, "Can't update feed")
		return
	}

	responseWithJSON(w, 200, databaseFeedtoFeed(feed))
}

// Delete Feed
func (apiCfg *apiConfig) handlerDeleteFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	id := chi.URLParam(r, "feed_id")
	feedID, err := uuid.Parse(id)
	if err != nil {
		responseWithError(w, 400, "Invalid feed id")
		return
	}

	feed, err := apiCfg.DB.GetFeed(r.Context(), feedID)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "Feed don't exsist")
		return
	}

	if feed.UserID != user.ID {
		responseWithError(w, 403, "Forbidden")
		return
	}
	err = apiCfg.DB.DeleteFeed(r.Context(), database.DeleteFeedParams{
		ID:     feedID,
		UserID: user.ID,
	})
	if err != nil {
		responseWithError(w, 500, "Can't delete feed")
		return
	}

	responseWithJSON(w, 204, map[string]string{"status": "No Content"})
}
