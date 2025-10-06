package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/ttrtcixy/workout/internal/app/provider"
)

type App struct {
	wg sync.WaitGroup
	*provider.Provider
}

func NewApp(ctx context.Context) *App {
	const op = "App.NewApp"

	p, err := provider.New(ctx)
	if err != nil {
		log.Fatalf("%s: error initializing provider: %s", op, err.Error())
	}
	return &App{
		Provider: p,
	}
}

func (a *App) Run(ctx context.Context) {
	defer a.Closer().Close()

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		if err := a.HTTPServer().Start(ctx); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				a.Logger().Error(err.Error())
			}
		}
	}()

	a.wg.Wait()
}
