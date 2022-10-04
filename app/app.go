package app

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mayowa/money/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type App struct {
	log zerolog.Logger
	mux *chi.Mux
	cfg *internal.Config
}

func (a *App) Init() error {
	var err error
	order := []string{
		zerolog.TimestampFieldName,
		zerolog.LevelFieldName,
		zerolog.MessageFieldName,
		zerolog.CallerFieldName,
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, PartsOrder: order})
	a.log = log.Logger

	// read config file
	a.cfg, err = internal.ReadConfig(nil)
	if err != nil {
		a.log.Error().Err(err).Msg("readConfig")
		return err
	}

	// configure mux
	a.mux = chi.NewMux()
	if a.cfg.Options.LogRequests {
		a.mux.Use(middleware.Logger)
	}
	a.mux.Use(middleware.CleanPath)
	a.mux.Use(middleware.Recoverer)

	if a.cfg.Options.Timeout > 0 {
		a.mux.Use(middleware.Timeout(time.Second * time.Duration(a.cfg.Options.Timeout)))
	}

	return nil
}

func (a *App) Run() error {
	if err := a.Init(); err != nil {
		return err
	}

	a.Route()

	// serve
	a.log.Debug().Msgf("listening on %s", a.cfg.Addr)
	err := http.ListenAndServe(a.cfg.Addr, a.mux)
	if err != nil {
		a.log.Error().Err(err).Msg("")
		return err
	}

	return nil
}
