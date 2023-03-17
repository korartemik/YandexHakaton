package stateful

import (
	"awesomeProject1/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	aliceapi "awesomeProject1/alice/api"
	"awesomeProject1/errors"
	"awesomeProject1/util"
)

func (h *Handler) askBot(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}

	intent := req.Request.NLU.Intents.AskBot
	if intent == nil {
		if req.State.Session.State == aliceapi.StateCreateReqAskBot {
			intent = req.Request.NLU.Intents.CompleteAsk
			if intent == nil {
				return nil, nil
			}
		} else {
			return nil, nil
		}
	}
	questionName, ok := intent.Slots.ListName.AsString()
	if !ok || len(questionName) == 0 {
		return &aliceapi.Response{Response: &aliceapi.Resp{Text: "Узнать о чём?", TTS: "Узнать о чём?"},
			State: &aliceapi.StateData{State: aliceapi.StateCreateReqAskBot}}, nil
	}
	if len(questionName) == 0 {
		return &aliceapi.Response{Response: &aliceapi.Resp{Text: "Узнать о чём?", TTS: "Узнать о чём?"},
			State: &aliceapi.StateData{State: aliceapi.StateCreateReqAskBot}}, nil
	}
	id := util.GenerateID()
	createAnswerOnAsk(ctx, id, questionName, h.config)
	var buttons []*aliceapi.Button
	buttons = append(buttons, &aliceapi.Button{Title: "Дальше"})
	text, tts := config.GetPhraseAfterAnswerRequest()
	return &aliceapi.Response{Response: &aliceapi.Resp{Text: text, TTS: tts,
		Buttons: buttons},
		State: &aliceapi.StateData{AnswerID: id, State: aliceapi.StateWaitAnswerFromChatGPT}}, nil
}

func createAnswerOnAsk(ctx context.Context, id string, questionName string, config *config.Config) {
	values := map[string]string{"idTag": id, "questionName": questionName, "idRequest": "2"}

	jsonValue, _ := json.Marshal(values)
	fmt.Println(string(jsonValue))

	go http.Post(config.FunctionEndpoint, "application/x-www-form-urlencoded", bytes.NewBuffer(jsonValue))
}
