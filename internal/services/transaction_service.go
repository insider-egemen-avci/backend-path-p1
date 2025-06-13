package services

import (
	"context"
	"fmt"
	"insider-egemen-avci/backend-path-p1/internal/models"
	"time"
)

type transactionService struct {
	transactionRepository models.TransactionRepository
	balanceRepository     models.BalanceRepository
	historyRepository     models.HistoryRepository
}

func NewTransactionService(transactionRepository models.TransactionRepository, balanceRepository models.BalanceRepository) models.TransactionService {
	return &transactionService{
		transactionRepository: transactionRepository,
		balanceRepository:     balanceRepository,
	}
}

func (service *transactionService) Deposit(ctx context.Context, userID int64, amount float64) (*models.Transaction, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	balance, err := service.balanceRepository.GetBalance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	balance.Credit(amount)

	if err := service.balanceRepository.UpdateBalance(ctx, balance); err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	transaction := &models.Transaction{
		ToUserID:  userID,
		Amount:    amount,
		Type:      models.TransactionTypeDeposit,
		Status:    models.TransactionStatusCompleted,
		CreatedAt: time.Now(),
	}

	history := &models.BalanceHistory{
		UserID:                userID,
		BalanceAfter:          balance.GetBalance(),
		TriggeringTransaction: transaction.ID,
	}

	if err := service.historyRepository.CreateHistory(ctx, history); err != nil {
		return nil, fmt.Errorf("failed to create history: %w", err)
	}

	if err := service.transactionRepository.CreateTransaction(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return transaction, nil
}

func (service *transactionService) Withdraw(ctx context.Context, userID int64, amount float64) (*models.Transaction, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	balance, err := service.balanceRepository.GetBalance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	if err := balance.Debit(amount); err != nil {
		return nil, fmt.Errorf("insufficient balance: %w", err)
	}

	if err := service.balanceRepository.UpdateBalance(ctx, balance); err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	transaction := &models.Transaction{
		FromUserID: userID,
		Amount:     amount,
		Type:       models.TransactionTypeWithdrawal,
		Status:     models.TransactionStatusCompleted,
		CreatedAt:  time.Now(),
	}

	if err := service.transactionRepository.CreateTransaction(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return transaction, nil
}

func (service *transactionService) Transfer(ctx context.Context, fromUserID, toUserID int64, amount float64) (*models.Transaction, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	if fromUserID == toUserID {
		return nil, fmt.Errorf("cannot transfer to yourself")
	}

	fromBalance, err := service.balanceRepository.GetBalance(ctx, fromUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	toBalance, err := service.balanceRepository.GetBalance(ctx, toUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	if err := fromBalance.Debit(amount); err != nil {
		return nil, fmt.Errorf("insufficient balance: %w", err)
	}

	toBalance.Credit(amount)

	if err := service.balanceRepository.UpdateBalance(ctx, fromBalance); err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	if err := service.balanceRepository.UpdateBalance(ctx, toBalance); err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	transaction := &models.Transaction{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Amount:     amount,
		Type:       models.TransactionTypeTransfer,
		Status:     models.TransactionStatusCompleted,
		CreatedAt:  time.Now(),
	}

	if err := service.transactionRepository.CreateTransaction(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return transaction, nil
}

func (service *transactionService) Rollback(ctx context.Context, transactionID int64) (*models.Transaction, error) {
	transaction, err := service.transactionRepository.GetTransactionByID(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	if transaction.Status != models.TransactionStatusPending {
		return nil, fmt.Errorf("transaction is not pending")
	}

	if models.TransactionType(transaction.Type) != models.TransactionTypeTransfer {
		return nil, fmt.Errorf("transaction is not a transfer")
	}

	reversedTransaction, err := service.Transfer(ctx, transaction.ToUserID, transaction.FromUserID, transaction.Amount)
	if err != nil {
		return nil, fmt.Errorf("failed to reverse transaction: %w", err)
	}

	if err := service.transactionRepository.UpdateTransactionStatus(ctx, transactionID, models.TransactionStatusRolledBack); err != nil {
		return nil, fmt.Errorf("failed to update transaction status: %w", err)
	}

	return reversedTransaction, nil
}
