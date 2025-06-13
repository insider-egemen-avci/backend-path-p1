package models

import (
	"fmt"
	"time"
)

type TransactionStatus string
type TransactionType string

const (
	TransactionStatusPending    TransactionStatus = "pending"
	TransactionStatusCompleted  TransactionStatus = "completed"
	TransactionStatusFailed     TransactionStatus = "failed"
	TransactionStatusRolledBack TransactionStatus = "rolled_back"
)

const (
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
	TransactionTypeTransfer   TransactionType = "transfer"
	TransactionTypePayment    TransactionType = "payment"
)

type Transaction struct {
	ID         int64             `json:"id"`
	FromUserID int64             `json:"from_user_id,omitempty"`
	ToUserID   int64             `json:"to_user_id,omitempty"`
	Amount     float64           `json:"amount"`
	Type       TransactionType   `json:"type"`
	Status     TransactionStatus `json:"status"`
	CreatedAt  time.Time         `json:"created_at"`
}

func (transaction *Transaction) Complete() error {
	if transaction.Status != TransactionStatusPending {
		return fmt.Errorf("transaction is not pending")
	}

	transaction.Status = TransactionStatusCompleted

	return nil
}

func (transaction *Transaction) Fail(reason string) error {
	if transaction.Status != TransactionStatusPending {
		return fmt.Errorf("transaction is not pending")
	}

	transaction.Status = TransactionStatusFailed

	return nil
}
