package stateful

import (
	"go.uber.org/zap"
)

type Deps interface {
	GetLogger() *zap.Logger
}
