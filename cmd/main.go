package main

import (
	"context"

	"github.com/imirjar/rb-diver/internal/app"
)

func main() {
	if err := app.Run(context.Background()); err != nil {
		panic(err.Error())
	}
}
