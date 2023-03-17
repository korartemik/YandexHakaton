package stateful

import (
	aliceapi "awesomeProject1/alice/api"
	"awesomeProject1/config"
	"awesomeProject1/errors"
	"awesomeProject1/util"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result/named"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	"log"
	"net/http"
)

func (h *Handler) askQuestion(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}

	intent := req.Request.NLU.Intents.AskMeQuestion
	if intent == nil {
		if req.State.Session.State == aliceapi.StateCreateReqTheme {
			intent = req.Request.NLU.Intents.CompleteTheme
			if intent == nil {
				return nil, nil
			}
		} else {
			return nil, nil
		}
	}
	themeName, ok := intent.Slots.ListName.AsString()
	if !ok || len(themeName) == 0 {
		return &aliceapi.Response{Response: &aliceapi.Resp{Text: "Вопрос по какой теме?"},
			State: &aliceapi.StateData{State: aliceapi.StateCreateReqTheme}}, nil
	}
	if len(themeName) == 0 {
		return &aliceapi.Response{Response: &aliceapi.Resp{Text: "Вопрос по какой теме?"},
			State: &aliceapi.StateData{State: aliceapi.StateCreateReqTheme}}, nil
	}
	id := util.GenerateID()
	idTheme := util.GenerateID()
	err := insertTheme(ctx, themeName, idTheme, h.config)
	if err != nil {
		log.Printf("Problem with db")
		log.Println(err)
	}
	log.Printf("Insert transaction succeeded")
	createQuestion(ctx, id, idTheme, h.config)
	var buttons []*aliceapi.Button
	buttons = append(buttons, &aliceapi.Button{Title: "Дальше"})
	text, tts := config.GetPhraseAfterQuestionRequest()
	return &aliceapi.Response{Response: &aliceapi.Resp{Text: text,
		TTS:     tts,
		Buttons: buttons},
		State: &aliceapi.StateData{AnswerID: id, NextQuestionID: id, ThemeID: idTheme, State: aliceapi.StateNextQuestion}}, nil
}

func insertTheme(ctx context.Context, theme string, themId string, config *config.Config) error {
	dsn := config.DatabaseEndpoint
	// IAM-токен
	token := config.IAMToken
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
		return err
	}
	// закрытие драйвера по окончании работы программы обязательно
	defer db.Close(ctx)
	var (
		readTx = table.TxControl(
			table.BeginTx(
				table.WithSerializableReadWrite(),
			),
			table.CommitTx(),
		)
	)
	err = db.Table().Do(ctx,
		func(ctx context.Context, s table.Session) (err error) {
			var (
				res result.Result
			)
			_, res, err = s.Execute(
				ctx,
				readTx,
				`
        DECLARE $id AS string;
		DECLARE $answer AS utf8;
		UPSERT INTO themes(id, theme) VALUES ($id, $answer);
      `,
				table.NewQueryParameters(
					table.ValueParam("$id", types.BytesValueFromString(themId)),
					table.ValueParam("$answer", types.UTF8Value(theme)), // подстановка в условие запроса
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
	return err
}

func (h *Handler) nextQuestion(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type == aliceapi.RequestTypeButton || (req.Request.Type == aliceapi.RequestTypeSimple && req.Request.NLU.Intents.NextQuestion != nil) {
		questionID := req.State.Session.NextQuestionID
		themeID := req.State.Session.ThemeID
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
		questionMain := ""
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
          question
        FROM
          questions
        WHERE
          id = $seriesID;
      `,
					table.NewQueryParameters(
						table.ValueParam("$seriesID", types.BytesValueFromString(questionID)), // подстановка в условие запроса
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
							named.Optional("question", &answer),
						)
						if err != nil {
							return err
						}

						questionMain = *answer

					}
				}
				return res.Err()
			},
		)
		if err != nil {
			text, tts := config.GetPhraseProblemWithBD()
			return &aliceapi.Response{Response: &aliceapi.Resp{Text: text, TTS: tts,
				Buttons: buttons},
				State: &aliceapi.StateData{AnswerID: questionID, NextQuestionID: questionID, ThemeID: themeID, State: aliceapi.StateNextQuestion}}, nil
		}
		if len(questionMain) == 0 {
			var buttons []*aliceapi.Button
			buttons = append(buttons, &aliceapi.Button{Title: "Дальше"})
			text, tts := config.GetWaitingText()
			return &aliceapi.Response{Response: &aliceapi.Resp{Text: text, TTS: tts,
				Buttons: buttons},
				State: &aliceapi.StateData{AnswerID: questionID, NextQuestionID: questionID, ThemeID: themeID, State: aliceapi.StateNextQuestion}}, nil
		}

		nextQuestionID := util.GenerateID()
		createQuestion(ctx, nextQuestionID, themeID, h.config)
		questionMain = questionMain
		if len(questionMain) <= 1024 {
			return &aliceapi.Response{Response: &aliceapi.Resp{Text: questionMain},
				State: &aliceapi.StateData{AnswerID: questionID, NextQuestionID: nextQuestionID, ThemeID: themeID, State: aliceapi.StateWaitAnswerFromClient}}, nil
		}
		return &aliceapi.Response{Response: &aliceapi.Resp{Text: questionMain[:1024]},
			State: &aliceapi.StateData{AnswerID: questionID, NextQuestionID: nextQuestionID, ThemeID: themeID, State: aliceapi.StateWaitAnswerFromClient}}, nil
	} else {
		return nil, nil
	}
}

func createQuestion(ctx context.Context, nextQuestID string, idTheme string, config *config.Config) {
	values := map[string]string{"nextQuestID": nextQuestID, "idTheme": idTheme}

	jsonValue, _ := json.Marshal(values)
	fmt.Println(string(jsonValue))

	go http.Post(config.FunctionQandAEndpoint, "application/x-www-form-urlencoded", bytes.NewBuffer(jsonValue))
}

func (h *Handler) sendAnswerToClient(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type == aliceapi.RequestTypeSimple && req.Request.NLU.Intents.WaitAnswer != nil {
		nextQuestionID := req.State.Session.NextQuestionID
		answerID := req.State.Session.AnswerID
		themeID := req.State.Session.ThemeID

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
          questions
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
			return &aliceapi.Response{Response: &aliceapi.Resp{Text: text, TTS: tts}}, nil
		}
		var buttons []*aliceapi.Button
		buttons = append(buttons, &aliceapi.Button{Title: "Следующий вопрос"})
		if len(answerMain) <= 1024 {
			return &aliceapi.Response{Response: &aliceapi.Resp{Text: answerMain, Buttons: buttons},
				State: &aliceapi.StateData{AnswerID: nextQuestionID, NextQuestionID: nextQuestionID, ThemeID: themeID, State: aliceapi.StateNextQuestion}}, nil
		}
		return &aliceapi.Response{Response: &aliceapi.Resp{Text: answerMain[:1024], Buttons: buttons},
			State: &aliceapi.StateData{AnswerID: nextQuestionID, NextQuestionID: nextQuestionID, ThemeID: themeID, State: aliceapi.StateNextQuestion}}, nil
	} else {
		return nil, nil
	}
}
