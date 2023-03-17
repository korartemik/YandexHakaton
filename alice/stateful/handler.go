package stateful

import (
	"awesomeProject1/config"
	"context"

	"awesomeProject1/alice/api"
	aliceapi "awesomeProject1/alice/api"
	"awesomeProject1/cache"
	"awesomeProject1/errors"
	"awesomeProject1/log"
	"go.uber.org/zap"
)

type Handler struct {
	logger           *zap.Logger
	config           *config.Config
	stateScenarios   map[api.State]scenario
	scratchScenarios []scenario
}

func NewHandler(deps Deps) (*Handler, error) {
	h := &Handler{
		logger: deps.GetLogger(),
		config: deps.GetConfig(),
	}
	h.setupScenarios()
	return h, nil
}

func (h *Handler) Handle(ctx context.Context, req *api.Request) (*api.Response, error) {
	sessionID := req.Session.SessionID
	ctx = log.CtxWithLogger(ctx, h.logger.With(zap.String("sessionID", string(sessionID))))
	ctx = cache.ContextWithCache(ctx)
	resp, err := h.handle(ctx, req)
	if err != nil {
		return h.reportError(ctx, err)
	}
	resp.Version = req.Version
	return resp, nil
}

func (h *Handler) handle(ctx context.Context, req *api.Request) (*api.Response, errors.Err) {
	var buttons []*aliceapi.Button
	buttons = append(buttons, &aliceapi.Button{Title: "Что ты умеешь?"})
	buttons = append(buttons, &aliceapi.Button{Title: "Расскажи про"})
	buttons = append(buttons, &aliceapi.Button{Title: "Задай мне вопрос по теме"})
	buttons = append(buttons, &aliceapi.Button{Title: "Узнай о"})
	if req.Session.New || req.AccountLinkingComplete != nil {
		text, tts := config.GetGreeting()
		return &api.Response{Response: &api.Resp{
			Text:    text,
			TTS:     tts,
			Buttons: buttons},
		}, nil
	}
	if state := req.State.Session; state.State != api.StateInit {
		intents := req.Request.NLU.Intents
		if intents.Reject != nil {
			return &api.Response{
				Response: &api.Resp{Text: "Чем я могу помочь?"},
			}, nil
		}
		scenario, ok := h.stateScenarios[state.State]
		if ok {
			resp, err := scenario(ctx, req)
			if err != nil {
				return nil, err
			}
			if resp != nil {
				return resp, nil
			}
		}
	}
	for _, s := range h.scratchScenarios {
		resp, err := s(ctx, req)
		if err != nil {
			return nil, err
		}
		if resp != nil {
			return resp, err
		}
	}
	text, tts := config.GetPhraseDoNotUnderstand()
	return &api.Response{Response: &api.Resp{
		Text: text,
		TTS:  tts,
	}}, nil
}

func (h *Handler) reportError(ctx context.Context, err errors.Err) (*api.Response, error) {
	errors.Log(ctx, err)
	return nil, err
}
