package handlers

import (
	"insider-egemen-avci/backend-path-p1/internal/models"
)

type APIHandler struct {
	userService        models.UserService
	transactionService models.TransactionService
	balanceService     models.BalanceService
	jwtSecret          string
}

func NewAPIHandler(
	userService models.UserService,
	txService models.TransactionService,
	balanceService models.BalanceService,
	jwtSecret string,
) *APIHandler {
	return &APIHandler{
		userService:        userService,
		transactionService: txService,
		balanceService:     balanceService,
		jwtSecret:          jwtSecret,
	}
}
