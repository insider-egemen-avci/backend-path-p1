package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type transactionPayload struct {
	Amount float64 `json:"amount"`
}

func (h *APIHandler) Credit(w http.ResponseWriter, r *http.Request) {
	authedUser, ok := UserFromContext(r.Context())
	if !ok {
		WriteJSONError(w, http.StatusInternalServerError, "User not found in context")
		return
	}

	var payload transactionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	tx, err := h.transactionService.Deposit(r.Context(), authedUser.ID, payload.Amount)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, tx)
}

func (h *APIHandler) Debit(w http.ResponseWriter, r *http.Request) {
	authedUser, ok := UserFromContext(r.Context())
	if !ok {
		WriteJSONError(w, http.StatusInternalServerError, "User not found in context")
		return
	}

	var payload transactionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	tx, err := h.transactionService.Withdraw(r.Context(), authedUser.ID, payload.Amount)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, tx)
}

type transferPayload struct {
	ToUserID int64   `json:"to_user_id"`
	Amount   float64 `json:"amount"`
}

func (h *APIHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	fromUser, ok := UserFromContext(r.Context())
	if !ok {
		WriteJSONError(w, http.StatusInternalServerError, "User not found in context")
		return
	}

	var payload transferPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	tx, err := h.transactionService.Transfer(r.Context(), fromUser.ID, payload.ToUserID, payload.Amount)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, tx)
}

func (h *APIHandler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	authedUser, ok := UserFromContext(r.Context())
	if !ok {
		WriteJSONError(w, http.StatusInternalServerError, "User not found in context")
		return
	}

	txID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	tx, err := h.transactionService.GetTransactionByID(r.Context(), txID)
	if err != nil {
		WriteJSONError(w, http.StatusNotFound, "Transaction not found")
		return
	}

	isParticipant := authedUser.ID == tx.FromUserID || authedUser.ID == tx.ToUserID
	isAdmin := authedUser.Role == "admin"

	if !isParticipant && !isAdmin {
		WriteJSONError(w, http.StatusForbidden, "You are not authorized to view this transaction")
		return
	}

	WriteJSON(w, http.StatusOK, tx)
}

func (h *APIHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	authedUser, ok := UserFromContext(r.Context())
	if !ok {
		WriteJSONError(w, http.StatusInternalServerError, "User not found in context")
		return
	}

	history, err := h.transactionService.GetUserHistory(r.Context(), authedUser.ID)
	if err != nil {
		slog.Error("Failed to get user history", "error", err)
		WriteJSONError(w, http.StatusInternalServerError, "Could not retrieve transaction history")
		return
	}

	WriteJSON(w, http.StatusOK, history)
}
