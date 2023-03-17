package stateful

import (
	"context"

	"awesomeProject1/alice/api"
	"awesomeProject1/errors"
)

type scenario = func(context.Context, *api.Request) (*api.Response, errors.Err)

func (h *Handler) setupScenarios() {
	h.stateScenarios = map[api.State]scenario{
		api.StateWaitAnswerFromChatGPT: h.getAnswer,
		api.StateCreateReqQuestion:     h.creatAnswerOnQuestion,
		api.StateCreateReqAskBot:       h.askBot,
		api.StateWaitAnswerFromClient:  h.sendAnswerToClient,
		api.StateNextQuestion:          h.nextQuestion,
		api.StateCreateReqTheme:        h.askQuestion,
		/*api.StateAddItemReqItem: h.addItemReqItem,
		api.StateAddItemReqList: h.addItemReqList,
		api.StateCreateReqName:  h.createRequireName,
		api.StateDelItemReqList: h.deleteItemReqList,
		api.StateDelItemReqItem: h.deleteItemReqItem,
		api.StateDelReqName:     h.deleteListReqList,
		api.StateDelReqConfirm:  h.deleteListReqConfirm,
		api.StateViewReqName:    h.viewListReqName,*/
	}
	h.scratchScenarios = []scenario{
		h.creatAnswerOnQuestion,
		h.helper,
		h.whatCanIDo,
		h.stop,
		h.agree,
		h.bye,
		h.thx,
		h.askBot,
		h.askQuestion,
	}
}
