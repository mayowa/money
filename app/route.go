package app

func (a *App) Route() {
	if a.mux == nil {
		a.log.Error().Msg("mux == nil")
		return
	}

}
