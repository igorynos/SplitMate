package model

import "time"

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
type Ledger struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	OwnerID   int64     `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
}
type Member struct {
	LedgerID int64     `json:"ledger_id"`
	UserID   int64     `json:"user_id"`
	WalletID int64     `json:"wallet_id"`
	JoinedAt time.Time `json:"joined_at"`
}
type Wallet struct {
	ID       int64  `json:"id"`
	LedgerID int64  `json:"ledger_id"`
	Name     string `json:"name"`
}
type Expense struct {
	ID           int64     `json:"id"`
	LedgerID     int64     `json:"ledger_id"`
	PaidBy       int64     `json:"paid_by"`
	Description  string    `json:"description"`
	Amount       int64     `json:"amount"`
	Participants []int64   `json:"participants"`
	CreatedAt    time.Time `json:"created_at"`
}
type Payment struct {
	ID        int64     `json:"id"`
	LedgerID  int64     `json:"ledger_id"`
	FromUser  int64     `json:"from_user"`
	ToUser    int64     `json:"to_user"`
	Amount    int64     `json:"amount"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
type Transfer struct {
	FromUser int64 `json:"from_user"`
	ToUser   int64 `json:"to_user"`
	Amount   int64 `json:"amount"`
}
