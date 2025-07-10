package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *APIHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	authedUser, ok := UserFromContext(r.Context())
	if !ok {
		WriteJSONError(w, http.StatusInternalServerError, "Failed to retrieve user from context")
		return
	}

	idStr := chi.URLParam(r, "id")
	targetUserID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid user ID specified")
		return
	}

	if authedUser.Role != "admin" && authedUser.ID != targetUserID {
		WriteJSONError(w, http.StatusForbidden, "You are not authorized to view this user's profile")
		return
	}

	user, err := h.userService.GetByID(r.Context(), targetUserID)
	if err != nil {
		slog.Error("Failed to get user by ID", "error", err)
		WriteJSONError(w, http.StatusNotFound, "User not found")
		return
	}

	WriteJSON(w, http.StatusOK, user)
}

func (h *APIHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers(r.Context())
	if err != nil {
		slog.Error("Failed to get all users", "error", err)
		WriteJSONError(w, http.StatusInternalServerError, "Could not retrieve users")
		return
	}

	WriteJSON(w, http.StatusOK, users)
}

type updateUserPayload struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
}

func (h *APIHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	authedUser, ok := UserFromContext(r.Context())
	if !ok {
		WriteJSONError(w, http.StatusInternalServerError, "User not found in context")
		return
	}
	targetUserID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	if authedUser.Role != "admin" && authedUser.ID != targetUserID {
		WriteJSONError(w, http.StatusForbidden, "Not authorized to update this user")
		return
	}

	userToUpdate, err := h.userService.GetByID(r.Context(), targetUserID)
	if err != nil {
		WriteJSONError(w, http.StatusNotFound, "User to update not found")
		return
	}

	var payload updateUserPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if payload.Username != nil {
		userToUpdate.Username = *payload.Username
	}
	if payload.Email != nil {
		userToUpdate.Email = *payload.Email
	}

	if err := h.userService.UpdateUser(r.Context(), userToUpdate); err != nil {
		slog.Error("Failed to update user", "error", err)
		WriteJSONError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	WriteJSON(w, http.StatusOK, userToUpdate)
}

func (h *APIHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	authedUser, ok := UserFromContext(r.Context())
	if !ok {
		WriteJSONError(w, http.StatusInternalServerError, "User not found in context")
		return
	}
	targetUserID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	if authedUser.Role != "admin" && authedUser.ID != targetUserID {
		WriteJSONError(w, http.StatusForbidden, "Not authorized to delete this user")
		return
	}

	if err := h.userService.DeleteUser(r.Context(), targetUserID); err != nil {
		slog.Error("Failed to delete user", "error", err)
		WriteJSONError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	w.WriteHeader(http.StatusOK)
}
