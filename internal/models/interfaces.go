package models

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id int64) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByRefreshToken(ctx context.Context, token string) (*User, error)
	SaveRefreshToken(ctx context.Context, userID int64, token string) error
	GetAll(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *Transaction) error
	GetTransactionByID(ctx context.Context, id int64) (*Transaction, error)
	UpdateTransactionStatus(ctx context.Context, id int64, status TransactionStatus) error
	GetHistoryByUserID(ctx context.Context, userID int64) ([]*Transaction, error)
}

type BalanceRepository interface {
	GetBalance(ctx context.Context, userID int64) (*Balance, error)
	UpdateBalance(ctx context.Context, balance *Balance) error
	GetBalanceHistory(ctx context.Context, userID int64) ([]*Balance, error)
}

type UserService interface {
	Register(ctx context.Context, username, email, password string) (*User, error)
	Login(ctx context.Context, email, password string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	GetUserByRefreshToken(ctx context.Context, token string) (*User, error)
	SaveRefreshToken(ctx context.Context, userID int64, token string) error
	GetAllUsers(ctx context.Context) ([]*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int64) error
}

type TransactionService interface {
	Deposit(ctx context.Context, userID int64, amount float64) (*Transaction, error)
	Withdraw(ctx context.Context, userID int64, amount float64) (*Transaction, error)
	Transfer(ctx context.Context, fromUserID, toUserID int64, amount float64) (*Transaction, error)
	Rollback(ctx context.Context, transactionID int64) (*Transaction, error)
	GetTransactionByID(ctx context.Context, txID int64) (*Transaction, error)
	GetUserHistory(ctx context.Context, userID int64) ([]*Transaction, error)
}

type BalanceService interface {
	GetBalance(ctx context.Context, userID int64) (*Balance, error)
	GetBalanceHistory(ctx context.Context, userID int64) ([]*Balance, error)
}

type HistoryRepository interface {
	CreateHistory(ctx context.Context, history *BalanceHistory) error
}
