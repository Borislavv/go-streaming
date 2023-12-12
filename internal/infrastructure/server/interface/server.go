package server_interface

import (
	"context"
	"sync"
)

type Server interface {
	Listen(ctx context.Context, wg *sync.WaitGroup)
}
