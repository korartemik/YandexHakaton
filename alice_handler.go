package main

import (
	"context"

	"awesomeProject1/alice"
	aliceapi "awesomeProject1/alice/api"
	"awesomeProject1/alice/stateful"
	"awesomeProject1/config"
	"awesomeProject1/log"
	"go.uber.org/zap"
)

type aliceApp struct {
	ctx     context.Context
	logger  *zap.Logger
	handler alice.Handler
	config  *config.Config
}

func (a *aliceApp) GetLogger() *zap.Logger {
	assertInitialized(a.logger, "logger")
	return a.logger
}

func (a *aliceApp) GetContext() context.Context {
	assertInitialized(a.ctx, "ctx")
	return a.ctx
}

func (a *aliceApp) GetConfig() *config.Config {
	assertInitialized(a.config, "config")
	return a.config
}

var aliceAppInstance *aliceApp

func initAliceApp() (*aliceApp, error) {
	ctx, err := initLogging()
	if err != nil {
		return nil, err
	}
	log.Info(ctx, "initializing alice app")

	aliceAppInstance = &aliceApp{ctx: ctx, logger: log.FromCtx(ctx)}
	aliceAppInstance.config = config.NewConfig()
	aliceAppInstance.handler, err = stateful.NewHandler(aliceAppInstance)
	//aliceAppInstance.handler, err = stateless.NewHandler(aliceAppInstance)
	if err != nil {
		return nil, err
	}
	return aliceAppInstance, nil
}

func getAliceApp() (*aliceApp, error) {
	if aliceAppInstance == nil {
		return initAliceApp()
	}
	return aliceAppInstance, nil
}

func AliceHandler(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, error) {
	aliceApp, err := getAliceApp()
	if err != nil {
		return nil, err
	}
	return aliceApp.handler.Handle(ctx, req)
}
