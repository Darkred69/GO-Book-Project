package main

import (
	"net/http"
	"project_1/internal/database"
)

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
