package main

import (
	"context"

	"github.com/ttrtcixy/workout/internal/app"
)

func main() {
	ctx := context.Background()

	a := app.NewApp(ctx)

	a.Run(context.Background())
}
