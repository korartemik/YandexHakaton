package alice

import (
	aliceapi "awesomeProject1/alice/api"
	"context"
)

type Handler interface {
	Handle(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, error)
}
