package main

import (
	"encoding/json"
	"net/http"
	"project_1/internal/database"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// handlerCreateFeed creates a new feed
// @Summary      Create feed
// @Description  Add a new RSS feed for the authenticated user
// @Tags         feeds
// @Accept       json
// @Produce      json
// @Param        feed  body      map[string]interface{}  true  "Feed data"
// @Success      201   {object}  map[string]interface{}    "Created feed data"
// @Failure      400   {object}  map[string]interface{}    "Bad request error"
// @Failure      409   {object}  map[string]interface{}    "Conflict error"
// @Router       /v2/feeds [post]
// @Security     BearerAuth
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

// handlerGetFeeds returns all available feeds
// @Summary      Get all feeds
// @Description  Retrieve a list of all feeds (publicly available or owned)
// @Tags         feeds
// @Produce      json
// @Success      200  {array}   map[string]interface{}    "List of feeds"
// @Failure      500  {object}  map[string]interface{}    "Internal server error"
// @Router       /v2/feeds [get]
// @Security     BearerAuth
func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := apiCfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		responseWithError(w, 500, "Can't get feeds")
		return
	}

	responseWithJSON(w, 200, databaseFeedstoFeeds(feeds))
}

// handlerUpdateFeed updates an existing feed
// @Summary      Update feed
// @Description  Modify the name or URL of a feed
// @Tags         feeds
// @Accept       json
// @Produce      json
// @Param        feed_id  path      string     true  "Feed ID"
// @Param        feed     body      map[string]interface{}  true  "Updated feed data"
// @Success      200      {object}  map[string]interface{}  "Updated feed data"
// @Failure      400      {object}  map[string]interface{}  "Bad request error"
// @Failure      403      {object}  map[string]interface{}  "Forbidden error"
// @Failure      404      {object}  map[string]interface{}  "Feed not found error"
// @Failure      409      {object}  map[string]interface{}  "Conflict error"
// @Failure      500      {object}  map[string]interface{}  "Internal server error"
// @Router       /v2/feeds/{feed_id} [put]
// @Security     BearerAuth
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

// handlerDeleteFeed deletes an existing feed
// @Summary      Delete feed
// @Description  Remove a feed owned by the authenticated user
// @Tags         feeds
// @Produce      json
// @Param        feed_id  path      string  true  "Feed ID"
// @Success      204      {object}  map[string]string  "Status: No Content"
// @Failure      400      {object}  map[string]string  "Bad request error"
// @Failure      403      {object}  map[string]string  "Forbidden error"
// @Failure      404      {object}  map[string]string  "Feed not found error"
// @Failure      500      {object}  map[string]string  "Internal server error"
// @Router       /v2/feeds/{feed_id} [delete]
// @Security     BearerAuth
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
