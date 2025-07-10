package handlers

import "net/http"

func (h *APIHandler) GetCurrentBalance(w http.ResponseWriter, r *http.Request) {
	authedUser, ok := UserFromContext(r.Context())
	if !ok {
		WriteJSONError(w, http.StatusInternalServerError, "User not found in context")
		return
	}

	balance, err := h.balanceService.GetBalance(r.Context(), authedUser.ID)
	if err != nil {
		WriteJSONError(w, http.StatusNotFound, "Balance not found for user")
		return
	}

	WriteJSON(w, http.StatusOK, balance)
}
