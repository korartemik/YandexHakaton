package bd

import (
	"context"

	"awesomeProject1/config"
)

type Deps interface {
	GetConfig() *config.Config
	GetContext() context.Context
}
