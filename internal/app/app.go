package app

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cart-overflow/cart/internal/cache"
	"github.com/cart-overflow/cart/internal/cart"
)

type Deps struct {
	Getenv  func(string) string
	Logwr   io.Writer
	Now     func() time.Time
	Started chan<- struct{}
}

func Run(ctx context.Context, deps Deps) {
	getenv := deps.Getenv
	log := log.New(deps.Logwr, "", log.LstdFlags)

	cacheCl, err := cache.NewCacheClient(getenv(ValkeyAddrKey))
	if err != nil {
		log.Fatalf("failed to init cache client")
	}
	defer cacheCl.Close()

	cartHandler := cart.Compose(
		cart.Deps{
			Log:     log,
			Now:     deps.Now,
			CacheCl: cacheCl,
		},
		cart.Config{},
	)

	server, err := NewServer(
		getenv(AddrKey),
		Handlers{
			Cart: cartHandler,
		},
	)
	if err != nil {
		log.Fatalf("init server error: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Go(func() {
		log.Printf("server started")
		deps.Started <- struct{}{}
		err := server.Run()
		if err != nil {
			log.Printf("server error: %v", err)
		}
		log.Printf("server stopped")
	})

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Printf("shutting down: context done")
	case sig := <-sc:
		log.Printf("shutting down: %v", sig)
	}

	server.Stop()
	wg.Wait()
	log.Printf("service stopped")
}
