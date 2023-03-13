package errors

import (
	"context"
	"fmt"

	"awesomeProject1/log"
	"go.uber.org/zap"
)

func Log(ctx context.Context, err error) {
	if e, ok := err.(Err); ok {
		if e.GetCode().IsUser() {
			log.Info(ctx, "user error", zap.Error(err))
			return
		}
		log.Error(ctx, fmt.Sprintf("code %s", e.GetCode()), zap.Error(e.Unwrap()))
		return
	}
	log.Error(ctx, "unexpected error", zap.Error(err))
}
