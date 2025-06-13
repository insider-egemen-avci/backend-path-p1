package models

import (
	"fmt"
	"sync"
	"time"
)

type Balance struct {
	UserID        int64     `json:"user_id"`
	Amount        float64   `json:"amount"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
	mu            sync.RWMutex
}

func NewBalance(userID int64, initialAmount float64) *Balance {
	return &Balance{
		UserID:        userID,
		Amount:        initialAmount,
		LastUpdatedAt: time.Now(),
	}
}

func (balance *Balance) Credit(amount float64) {
	balance.mu.Lock()
	defer balance.mu.Unlock()

	balance.Amount += amount
	balance.LastUpdatedAt = time.Now()
}

func (balance *Balance) Debit(amount float64) error {
	balance.mu.Lock()
	defer balance.mu.Unlock()

	if balance.Amount < amount {
		return fmt.Errorf("insufficient balance")
	}

	balance.Amount -= amount
	balance.LastUpdatedAt = time.Now()

	return nil
}

func (balance *Balance) GetBalance() float64 {
	balance.mu.RLock()
	defer balance.mu.RUnlock()

	return balance.Amount
}

type BalanceHistory struct {
	UserID                int64     `json:"user_id"`
	BalanceAfter          float64   `json:"balance_after"`
	TriggeringTransaction int64     `json:"triggering_transaction"`
	CreatedAt             time.Time `json:"created_at"`
}

func NewBalanceHistory(userID int64, balanceAfter float64, transactionID int64) *BalanceHistory {
	return &BalanceHistory{
		UserID:                userID,
		BalanceAfter:          balanceAfter,
		TriggeringTransaction: transactionID,
		CreatedAt:             time.Now(),
	}
}
