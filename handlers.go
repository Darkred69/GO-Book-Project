package main

import (
	"net/http"
	"project_1/internal/database"
)

// handlerGetPosts retrieves posts for the authenticated user
// @Summary      Get user posts
// @Description  Retrieve a list of posts belonging to the authenticated user
// @Tags         posts
// @Produce      json
// @Success      200  {array}   map[string]interface{} "List of posts"
// @Failure      500  {object}  map[string]interface{} "Internal Server Error"
// @Router       /v1/posts [get]
// @Security     BearerAuth
func (apiCfg *apiConfig) handlerGetPosts(w http.ResponseWriter, r *http.Request, user database.User) {

	posts, err := apiCfg.DB.GetPosts(r.Context(), database.GetPostsParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		responseWithError(w, 500, "Can't get posts")
		return
	}
	responseWithJSON(w, 200, databasePoststoPosts(posts))
}
