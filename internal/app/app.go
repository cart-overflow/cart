package app

import (
	"context"
	"io"
)

type Deps struct {
	Getenv  func(string) string
	Logwr   io.Writer
	Started chan<- struct{}
}

func Run(ctx context.Context, deps Deps) {

}
