package main

import (
	"context"
	"os"
	"time"

	"github.com/cart-overflow/cart/internal/app"
)

func main() {
	app.Run(context.Background(), app.Deps{
		Getenv:  os.Getenv,
		Logwr:   os.Stdout,
		Now:     time.Now,
		Started: make(chan struct{}, 1),
	})
}
