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
		/*h.viewListFromScratch,
		h.listAllListsFromScratch,
		h.createFromScratch,
		h.addItemFromScratch,
		h.deleteListFromScratch,
		h.deleteItemFromScratch,*/
	}
}
