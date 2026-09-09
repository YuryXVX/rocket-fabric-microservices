package app

import "context"

type App struct {
	diContainer *diContainer
}

func New(cxt context.Context) *App {
	a := &App{}

	return a
}

func (a *App) Run() error {
	return nil
}
