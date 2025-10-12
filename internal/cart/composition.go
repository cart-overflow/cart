package cart

import (
	"log"
	"time"

	"github.com/valkey-io/valkey-go"
)

type Deps struct {
	Log     *log.Logger
	CacheCl valkey.Client
	Now     func() time.Time
}

type Config struct {
}

func Compose(deps Deps, cfg Config) *Handler {
	cache := newCache(deps.CacheCl, deps.Now, deps.Log)
	service := newService(cache, deps.Log)
	return newHandler(service)
}
