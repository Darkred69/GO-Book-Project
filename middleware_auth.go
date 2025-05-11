package main

import (
	"net/http"
	"project_1/auth"
	"project_1/internal/database"

	"github.com/google/uuid"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

// middlewareAuth is a middleware that checks if the user is authenticated
// and retrieves the user information from the database.
func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user_id, err := auth.GetUserID(r.Header)
		if err != nil {
			responseWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		parsedUserID, err := uuid.Parse(user_id)
		if err != nil {
			responseWithError(w, 400, "Invalid user ID format")
			return
		}
		user, err := apiCfg.DB.GetUserByID(r.Context(), parsedUserID)
		if err != nil {
			responseWithError(w, http.StatusNotFound, "User don't exsist")
			return
		}
		handler(w, r, user)
	}
}
