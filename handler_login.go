package main

import (
	"net/http"
	"strings"

	"project_1/auth"

	"github.com/badoux/checkmail"
)

// @Summary Login to the system
// @Description Returns an authentication token
// @Tags authentication
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Auth Token Response"
// @Failure 400 {object} map[string]string "Bad Request Error"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /auth/login [post]
func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("username")
	password := r.FormValue("password")
	err := checkmail.ValidateFormat(email)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid email")
		return
	}

	user, err := apiCfg.DB.GetUser(r.Context(), email)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "User not found")
		return
	}
	if !checkHash(password, user.Password) {
		responseWithError(w, http.StatusUnauthorized, "Wrong password")
		return
	}

	token, err := auth.GenerateToken(user.ID.String())
	if err != nil {
		responseWithError(w, 500, "Can't generate token")
		return
	}
	token_type := strings.Split(token, " ")[0]
	access_token := strings.Split(token, " ")[1]
	responseWithJSON(w, 200, ResponseToken(token_type, access_token))
}
