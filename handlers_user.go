package main

import (
	"encoding/json"
	"net/http"
	"project_1/internal/database"

	"github.com/badoux/checkmail"
	"github.com/google/uuid"
)

// Create user handler
func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p UserInput
	err := decoder.Decode(&p)
	if err != nil {
		responseWithError(w, 400, "Invalid request payload")
		return
	}
	hashPassword := HashPassword(p.Password)

	err = checkmail.ValidateFormat(p.Email)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid email")
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:       uuid.New(),
		Name:     p.Name,
		Email:    p.Email,
		Password: hashPassword,
	})
	if err != nil {
		responseWithError(w, http.StatusConflict, "Account already exists")
		return
	}

	responseWithJSON(w, 201, databaseUsertoUser(user))
}

// Get a user handler
func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	responseWithJSON(w, 200, databaseUsertoUser(user))
}

// Delete a User
func (apiCfg *apiConfig) handlerDeleteUser(w http.ResponseWriter, r *http.Request, user database.User) {
	err := apiCfg.DB.DeleteUser(r.Context(), user.ID)
	if err != nil {
		responseWithError(w, 500, "Can't delete user")
		return
	}
	responseWithJSON(w, 204, map[string]string{"status": "No Content"})
}

// Update a user handler
func (apiCfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request, user database.User) {
	decoder := json.NewDecoder(r.Body)
	var p UserInput
	err := decoder.Decode(&p)
	if err != nil {
		responseWithError(w, 400, "Invalid request payload")
		return
	}

	err = checkmail.ValidateFormat(p.Email)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid email")
		return
	}

	user, err = apiCfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:    user.ID,
		Name:  p.Name,
		Email: p.Email,
	})
	if err != nil {
		responseWithError(w, http.StatusConflict, "Account already exists")
		return
	}

	responseWithJSON(w, 200, databaseUsertoUser(user))
}
