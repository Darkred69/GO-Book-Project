package main

import (
	"encoding/json"
	"net/http"
	"project_1/internal/database"

	"github.com/badoux/checkmail"
	"github.com/google/uuid"
)

// handlerCreateUser creates a new user account
// @Summary      Create user
// @Description  Register a new user with name, email, and password
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user  body      UserInput  true  "User registration input"
// @Success      201   {object}  User
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Router       /v1/user [post]
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

// handlerGetUser returns the authenticated user's data
// @Summary      Get user
// @Description  Retrieve the current authenticated user's profile
// @Tags         user
// @Produce      json
// @Success      200  {object}  User
// @Failure      401  {object}  map[string]string
// @Router       /v1/user [get]
// @Security     BearerAuth
func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	responseWithJSON(w, 200, databaseUsertoUser(user))
}

// handlerDeleteUser deletes the authenticated user's account
// @Summary      Delete user
// @Description  Delete the current authenticated user's account
// @Tags         user
// @Produce      json
// @Success      204  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /v1/user [delete]
// @Security     BearerAuth
func (apiCfg *apiConfig) handlerDeleteUser(w http.ResponseWriter, r *http.Request, user database.User) {
	err := apiCfg.DB.DeleteUser(r.Context(), user.ID)
	if err != nil {
		responseWithError(w, 500, "Can't delete user")
		return
	}
	responseWithJSON(w, 204, map[string]string{"status": "No Content"})
}

// handlerUpdateUser updates the authenticated user's profile
// @Summary      Update user
// @Description  Update name and email of the authenticated user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user  body      UserInput  true  "Updated user info"
// @Success      200   {object}  User
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Router       /v1/user [put]
// @Security     BearerAuth
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
