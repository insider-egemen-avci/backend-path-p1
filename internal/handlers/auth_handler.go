package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"insider-egemen-avci/backend-path-p1/internal/models"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *APIHandler) Register(w http.ResponseWriter, r *http.Request) {
	var payload registerPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	newUser, err := h.userService.Register(
		r.Context(),
		payload.Username,
		payload.Email,
		payload.Password,
	)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			WriteJSONError(w, http.StatusConflict, err.Error())
		} else {
			WriteJSONError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	WriteJSON(w, http.StatusCreated, newUser)
}

func (h *APIHandler) generateTokens(user *models.User) (accessToken string, refreshToken string, err error) {
	accessExpirationTime := time.Now().Add(15 * time.Minute)
	accessClaims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", user.ID),
		ExpiresAt: jwt.NewNumericDate(accessExpirationTime),
		Issuer:    "go-fintech-app",
	}
	accessTokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenJwt.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return "", "", err
	}

	refreshToken = uuid.New().String()
	err = h.userService.SaveRefreshToken(context.Background(), user.ID, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (h *APIHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload loginPayload

	user, err := h.userService.Login(r.Context(), payload.Email, payload.Password)
	if err != nil {
		WriteJSONError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	accessToken, refreshToken, err := h.generateTokens(user)
	if err != nil {
		slog.Error("Failed to generate tokens", "error", err)
		WriteJSONError(w, http.StatusInternalServerError, "Failed to process login")
		return
	}

	response := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	WriteJSON(w, http.StatusOK, response)
}

func (h *APIHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.userService.GetUserByRefreshToken(r.Context(), payload.RefreshToken)
	if err != nil {
		WriteJSONError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	accessToken, refreshToken, err := h.generateTokens(user)
	if err != nil {
		slog.Error("Failed to generate tokens on refresh", "error", err)
		WriteJSONError(w, http.StatusInternalServerError, "Failed to process refresh")
		return
	}

	response := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	WriteJSON(w, http.StatusOK, response)
}
