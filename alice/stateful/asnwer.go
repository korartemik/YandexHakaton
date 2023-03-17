package stateful

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	"net/http"

	aliceapi "awesomeProject1/alice/api"
	"awesomeProject1/config"
	"awesomeProject1/errors"
	"awesomeProject1/util"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result/named"
	"log"
)

func (h *Handler) creatAnswerOnQuestion(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}

	intent := req.Request.NLU.Intents.CreateAnswer
	if intent == nil {
		if req.State.Session.State == aliceapi.StateCreateReqQuestion {
			intent = req.Request.NLU.Intents.CompleteQuestion
			if intent == nil {
				return nil, nil
			}
		} else {
			return nil, nil
		}
	}
	questionName, ok := intent.Slots.ListName.AsString()
	if !ok || len(questionName) == 0 {
		return &aliceapi.Response{Response: &aliceapi.Resp{Text: "Рассказать про что?"},
			State: &aliceapi.StateData{State: aliceapi.StateCreateReqQuestion}}, nil
	}
	if len(questionName) == 0 {
		return &aliceapi.Response{Response: &aliceapi.Resp{Text: "Рассказать про что?"},
			State: &aliceapi.StateData{State: aliceapi.StateCreateReqQuestion}}, nil
	}
	id := util.GenerateID()
	createAnswer(ctx, id, questionName, h.config)
	var buttons []*aliceapi.Button
	buttons = append(buttons, &aliceapi.Button{Title: "Дальше"})
	text, tts := config.GetPhraseAfterAnswerRequest()
	return &aliceapi.Response{Response: &aliceapi.Resp{Text: text, TTS: tts,
		Buttons: buttons},
		State: &aliceapi.StateData{AnswerID: id, State: aliceapi.StateWaitAnswerFromChatGPT}}, nil
}

func (h *Handler) getAnswer(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	answerID := req.State.Session.AnswerID
	if req.Request.Type == aliceapi.RequestTypeButton || (req.Request.Type == aliceapi.RequestTypeSimple && (req.Request.NLU.Intents.Next != nil || req.Request.NLU.Intents.ButtonNext != nil)) {
		var buttons []*aliceapi.Button
		buttons = append(buttons, &aliceapi.Button{Title: "Дальше"})
		// строка подключения
		dsn := h.config.DatabaseEndpoint
		// IAM-токен
		token := h.config.IAMToken
		// создаем объект подключения db, является входной точкой для сервисов YDB
		db, err := ydb.Open(ctx, dsn,
			//  yc.WithInternalCA(), // используем сертификаты Яндекс Облака
			ydb.WithAccessTokenCredentials(token), // аутентификация с помощью токена
			//  ydb.WithAnonimousCredentials(), // анонимная аутентификация (например, в docker ydb)
			//  yc.WithMetadataCredentials(token), // аутентификация изнутри виртуальной машины в Яндекс Облаке или из Яндекс Функции
			//  yc.WithServiceAccountKeyFileCredentials("~/.ydb/sa.json"), // аутентификация в Яндекс Облаке с помощью файла сервисного аккаунта
			//  environ.WithEnvironCredentials(ctx), // аутентификация с использованием переменных окружения
		)
		if err != nil {
			text, tts := config.GetPhraseProblemWithBD()
			return &aliceapi.Response{Response: &aliceapi.Resp{Text: text, TTS: tts}}, nil
		}
		// закрытие драйвера по окончании работы программы обязательно
		defer db.Close(ctx)

		var (
			readTx = table.TxControl(
				table.BeginTx(
					table.WithOnlineReadOnly(),
				),
				table.CommitTx(),
			)
		)
		answerMain := ""
		err = db.Table().Do(ctx,
			func(ctx context.Context, s table.Session) (err error) {
				var (
					res    result.Result
					answer *string // указатель - для опциональных результатов
				)
				_, res, err = s.Execute(
					ctx,
					readTx,
					`
        DECLARE $seriesID AS String;
        SELECT
          id,
          answer
        FROM
          answers
        WHERE
          id = $seriesID;
      `,
					table.NewQueryParameters(
						table.ValueParam("$seriesID", types.BytesValueFromString(answerID)), // подстановка в условие запроса
					),
				)
				if err != nil {
					return err
				}
				defer res.Close() // закрытие result'а обязательно
				log.Printf("> select_simple_transaction:\n")
				for res.NextResultSet(ctx) {
					for res.NextRow() {
						// в ScanNamed передаем имена колонок из строки сканирования,
						// адреса (и типы данных), куда следует присвоить результаты запроса
						err = res.ScanNamed(
							named.Optional("answer", &answer),
						)
						if err != nil {
							return err
						}

						answerMain = *answer

					}
				}
				return res.Err()
			},
		)
		if err != nil {
			text, tts := config.GetPhraseProblemWithBD()
			return &aliceapi.Response{Response: &aliceapi.Resp{Text: text, TTS: tts,
				Buttons: buttons},
				State: &aliceapi.StateData{AnswerID: answerID, State: aliceapi.StateWaitAnswerFromChatGPT}}, nil
		}
		if len(answerMain) == 0 {
			var buttons []*aliceapi.Button
			buttons = append(buttons, &aliceapi.Button{Title: "Дальше"})
			text, tts := config.GetWaitingText()
			return &aliceapi.Response{Response: &aliceapi.Resp{Text: text, TTS: tts,
				Buttons: buttons},
				State: &aliceapi.StateData{AnswerID: answerID, State: aliceapi.StateWaitAnswerFromChatGPT}}, nil
		}
		if len(answerMain) <= 1024 {
			return &aliceapi.Response{Response: &aliceapi.Resp{Text: answerMain}}, nil
		}
		startIndex := changeAnswer(ctx, answerID, answerMain, db, 1024)
		if startIndex == -1 {
			text, tts := config.GetPhraseProblemWithBD()
			return &aliceapi.Response{Response: &aliceapi.Resp{Text: text, TTS: tts,
				Buttons: buttons},
				State: &aliceapi.StateData{AnswerID: answerID, State: aliceapi.StateWaitAnswerFromChatGPT}}, nil
		}
		return &aliceapi.Response{Response: &aliceapi.Resp{Text: answerMain[:startIndex], Buttons: buttons},
			State: &aliceapi.StateData{AnswerID: answerID, State: aliceapi.StateWaitAnswerFromChatGPT}}, nil
	} else {
		return nil, nil
	}
}
func createAnswer(ctx context.Context, id string, questionName string, config *config.Config) {
	values := map[string]string{"idTag": id, "questionName": questionName, "idRequest": "1"}

	jsonValue, _ := json.Marshal(values)
	fmt.Println(string(jsonValue))

	go http.Post(config.FunctionEndpoint, "application/x-www-form-urlencoded", bytes.NewBuffer(jsonValue))
}

func changeAnswer(ctx context.Context, answerID string, answerMain string, db *ydb.Driver, index int) int {
	if index == 1022 {
		return -1
	}
	var (
		readTx2 = table.TxControl(
			table.BeginTx(
				table.WithSerializableReadWrite(),
			),
			table.CommitTx(),
		)
	)
	err := db.Table().Do(ctx,
		func(ctx context.Context, s table.Session) (err error) {
			var (
				res result.Result
			)
			_, res, err = s.Execute(
				ctx,
				readTx2,
				`
        DECLARE $id AS string;
		DECLARE $answer AS utf8;
		UPSERT INTO answers(id, answer) VALUES ($id, $answer);
      `,
				table.NewQueryParameters(
					table.ValueParam("$id", types.BytesValueFromString(answerID)),
					table.ValueParam("$answer", types.UTF8Value(answerMain[index:])), // подстановка в условие запроса
				),
			)
			if err != nil {
				return err
			}
			defer res.Close() // закрытие result'а обязательно
			log.Printf("< insert_simple_transaction:\n")

			return res.Err()
		},
	)
	if err != nil {
		return changeAnswer(ctx, answerID, answerMain, db, index-1)
	} else {
		return index
	}

}
