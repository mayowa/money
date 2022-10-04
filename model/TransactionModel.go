package model

import (
	"context"
	"time"

	"github.com/rs/xid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type TransactionTypes int

const (
	ExpenseTransaction TransactionTypes = iota + 1
	IncomeTransaction
)

type Transaction struct {
	ID              xid.ID
	BookID          xid.ID
	Description     string
	Amount          decimal.Decimal
	Type            TransactionTypes
	DateTransaction time.Time
	DateCreated     time.Time
}

type TransactionModeler interface {
	DBSetter
	Get(id string) (*Transaction, error)
	Insert(record *Transaction) error
	Delete(where string, args ...any) error
	Update(data map[string]any, where string, args ...any) error
}

var _ TransactionModeler = (*TransactionModel)(nil)

type TransactionModel struct {
	db  *bun.DB
	ctx context.Context
}

func (m *TransactionModel) SetDB(db *bun.DB) {
	m.db = db
	m.ctx = context.Background()
}

func (m *TransactionModel) Get(id string) (*Transaction, error) {
	var record Transaction
	err := GetRecord(m.ctx, m.db, "", id, &record)
	return &record, err
}

func (m *TransactionModel) Insert(record *Transaction) error {
	return InsertRecord(m.ctx, m.db, "", record)
}

func (m *TransactionModel) Delete(where string, args ...any) error {
	return DeleteRecord(m.ctx, m.db, "transactions", where, args...)
}

func (m *TransactionModel) Update(data map[string]any, where string, args ...any) error {
	return UpdateRecord(m.ctx, m.db, "transactions", data, where, args...)
}
