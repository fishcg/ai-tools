package server

import (
	"context"

	"github.com/fish/ai-tools/controllers"
)

type Server interface {
	// Name() string
	Run(m controllers.Mount)
	Shutdown(ctx context.Context) error
}
