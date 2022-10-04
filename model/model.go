package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mayowa/money/internal"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
)

type DBSetter interface {
	SetDB(db *bun.DB)
}

type Models struct {
	Book        BookModeler
	Transaction TransactionModeler
}

func (m *Models) SetDB(db *bun.DB) {
	m.Book.SetDB(db)
	m.Transaction.SetDB(db)
}

func NewModel(db *bun.DB) (*Models, error) {
	mdl := &Models{
		Book:        &BookModel{},
		Transaction: &TransactionModel{},
	}

	mdl.SetDB(db)
	return mdl, nil
}

func OpenDb(dsn string) (db *bun.DB, err error) {

	sqlDB, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	db = bun.NewDB(sqlDB, sqlitedialect.New(), bun.WithDiscardUnknownColumns())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))

	return db, nil
}

func GetDb(dbName string, cfg *internal.Config) (*bun.DB, error) {

	if dbName[0] != '.' && dbName[0] != '/' {
		dbName = filepath.Join(cfg.Folders.Data, dbName)
	}

	dsn := fmt.Sprint("file:", dbName, "?_journal=WAL&cache=shared&_foreign_keys=true")
	db, err := OpenDb(dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetRecord(ctx context.Context, db bun.IDB, table, id string, model any) (err error) {
	if db == nil {
		return errors.New("db == nil")
	}

	stm := db.NewSelect().Model(model).Where("id = ?", id)
	if len(table) > 0 {
		stm.TableExpr(table)
	}
	err = stm.Scan(ctx)
	return err
}

func InsertRecord(ctx context.Context, db bun.IDB, table string, model any) (err error) {
	if db == nil {
		return errors.New("db == nil")
	}

	stm := db.NewInsert().Model(model)
	if len(table) > 0 {
		stm.TableExpr(table)
	}
	_, err = stm.Exec(ctx)

	return nil
}

func UpdateRecord(ctx context.Context, db bun.IDB, table string, data map[string]any, where string, args ...any) (err error) {
	if db == nil {
		return errors.New("db == nil")
	}

	stm := db.NewUpdate().Model(data).TableExpr(table).Where(where, args...)
	_, err = stm.Exec(ctx)

	return err
}

func DeleteRecord(ctx context.Context, db bun.IDB, table string, where string, args ...any) (err error) {
	if db == nil {
		return errors.New("db == nil")
	}

	stm := db.NewDelete().TableExpr(table).Where(where, args...)
	_, err = stm.Exec(ctx)

	return err
}

type TxHandler func(ctx context.Context, tx bun.Tx) error

func RunInTx(ctx context.Context, dbc bun.IDB, opts *sql.TxOptions, fn TxHandler) error {
	var (
		db   *bun.DB
		tx   bun.Tx
		err  error
		isDb bool
		isTx bool
	)

	// create a new transaction if dbc is *DB or use the provided transaction if dbc = Tx
	db, isDb = dbc.(*bun.DB)
	if isDb {
		tx, err = db.BeginTx(ctx, opts)
		if err != nil {
			return err
		}
	} else {
		tx, isTx = dbc.(bun.Tx)
	}

	if !isTx && !isDb {
		return errors.New("invalid type in dbc: must be either *bun.DB or bun.Tx")
	}

	// run closure
	if err := fn(ctx, tx); err != nil {
		// rollback if the transaction was created here
		if isDb {
			rErr := tx.Rollback()
			if rErr != nil {
				err = fmt.Errorf("%s - %w", rErr, err)
			}
		}

		return err
	}

	// commit if the transaction was created here
	if isDb {
		return tx.Commit()
	}

	return nil
}
