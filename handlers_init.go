package main

import "net/http"

// @Summary Get a hello message
// @Description Returns a simple hello message
// @Tags example
// @Accept  json
// @Produce json
// @Success 200 {string} string "Hello, World!"
// @Router /  [get]
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	responseWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Handler to check readiness of the server with an error
func handlerErr(w http.ResponseWriter, r *http.Request) {
	responseWithError(w, 400, "Internal Server Error")
}
