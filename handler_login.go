package main

import (
	"net/http"
	"strings"

	"project_1/auth"

	"github.com/badoux/checkmail"
)

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
