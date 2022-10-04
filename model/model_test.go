package model

import (
	"context"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mayowa/money/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

var cfg *internal.Config

func TestMain(m *testing.M) {
	var err error
	cfg, err = internal.ReadConfig(nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to read config")
		os.Exit(1)
	}

	order := []string{
		zerolog.TimestampFieldName,
		zerolog.LevelFieldName,
		zerolog.MessageFieldName,
		zerolog.CallerFieldName,
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, PartsOrder: order})

	if err := os.Remove("../data/test.db"); err != nil {
		log.Error().Err(err).Msg("failed to remove test.db")
	}

	mg, err := migrate.New(
		"file://../migrations",
		"sqlite3://../data/test.db?_journal=WAL&cache=shared&_foreign_keys=true")
	if err != nil {
		wd, _ := os.Getwd()
		log.Error().Err(err).Str("cwd", wd).Msg("migration setup failed")
		os.Exit(1)
	}
	if err = mg.Up(); err != nil {
		log.Error().Err(err).Msg("migration failed")
		os.Exit(1)
	}

	if ec := m.Run(); ec != 0 {
		log.Error().Msgf("tests failed: %d", ec)
		os.Exit(ec)
	}
}

func TestGetRecord(t *testing.T) {
	db, err := GetDb("../data/test.db", cfg)
	require.NoError(t, err, "error when connecting to test.db")

	err = insertSQL(t, db, [2][]string{
		{`insert into books (id, name) values('9bsv0s7lib40035evmug', 'foo')`},
		{`delete from books`},
	})
	require.NoError(t, err, "error inserting test data")

	record := Book{}
	err = GetRecord(context.Background(), db, "", "9bsv0s7lib40035evmug", &record)
	require.NoError(t, err, "GetRecord")
	assert.Equal(t, "foo", record.Name)
}

func insertSQL(t *testing.T, db *bun.DB, data [2][]string) error {
	exec := func(which int) error {
		return RunInTx(context.Background(), db, nil, func(ctx context.Context, tx bun.Tx) error {
			for _, s := range data[which] {
				if _, err := db.Exec(s); err != nil {
					return err
				}
			}

			return nil
		})
	}

	if err := exec(0); err != nil {
		return err
	}

	t.Cleanup(func() {
		if err := exec(1); err != nil {
			t.Log(err)
		}
	})

	return nil
}
