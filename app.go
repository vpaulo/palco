package main

import (
	"context"
	"fmt"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) domready(ctx context.Context) {
	// runtime.EventsEmit(a.ctx, "loaded")
}

func (a *App) shutdown(ctx context.Context) {
	// TODO:
}

func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	// TODO:
	return false
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
