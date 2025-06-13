package services

import (
	"context"
	"insider-egemen-avci/backend-path-p1/internal/models"
)

type balanceService struct {
	balanceRepo models.BalanceRepository
}

func NewBalanceService(balanceRepo models.BalanceRepository) models.BalanceService {
	return &balanceService{
		balanceRepo: balanceRepo,
	}
}

func (s *balanceService) GetBalance(ctx context.Context, userID int64) (*models.Balance, error) {
	return s.balanceRepo.GetBalance(ctx, userID)
}
