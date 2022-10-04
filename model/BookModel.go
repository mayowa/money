package model

import (
	"context"

	"github.com/rs/xid"
	"github.com/uptrace/bun"
)

type Book struct {
	ID   xid.ID
	Name string
}

type BookModeler interface {
	DBSetter
	Get(id string) (*Book, error)
	Insert(record *Book) error
	Delete(where string, args ...any) error
	Update(data map[string]any, where string, args ...any) error
}

var _ BookModeler = (*BookModel)(nil)

type BookModel struct {
	db  *bun.DB
	ctx context.Context
}

func (m *BookModel) SetDB(db *bun.DB) {
	m.db = db
	m.ctx = context.Background()
}

func (m *BookModel) Get(id string) (*Book, error) {
	var record Book
	err := GetRecord(m.ctx, m.db, "", id, &record)
	return &record, err
}

func (m *BookModel) Insert(record *Book) error {
	return InsertRecord(m.ctx, m.db, "", record)
}

func (m *BookModel) Delete(where string, args ...any) error {
	return DeleteRecord(m.ctx, m.db, "books", where, args...)
}

func (m *BookModel) Update(data map[string]any, where string, args ...any) error {
	return UpdateRecord(m.ctx, m.db, "books", data, where, args...)
}
