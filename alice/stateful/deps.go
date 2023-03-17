package stateful

import (
	"awesomeProject1/config"
	"go.uber.org/zap"
)

type Deps interface {
	GetLogger() *zap.Logger
	GetConfig() *config.Config
}
