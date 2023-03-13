package stateful

import (
	"context"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"

	aliceapi "awesomeProject1/alice/api"
	"awesomeProject1/config"
	"awesomeProject1/errors"
	"awesomeProject1/util"
	openai "github.com/sashabaranov/go-openai"
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
		return nil, nil
	}
	questionName, ok := intent.Slots.ListName.AsString()
	if !ok {
		return &aliceapi.Response{Response: &aliceapi.Resp{Text: "Так, ваш вопрос оченб и очень странный", TTS: "Так, ваш вопрос оченб и очень странный"}}, nil
	}
	id := util.GenerateID()
	go createAnswer(ctx, id, questionName)
	var buttons []*aliceapi.Button
	buttons = append(buttons, &aliceapi.Button{Title: "Дальше"})

	/*if err != nil {
		return &aliceapi.Response{Response: &aliceapi.Resp{Text: "Хм, это тема слишком сложная для меня, попробуй узнать у своих одногруппников или преподавателя"}}, nil
	}*/
	return &aliceapi.Response{Response: &aliceapi.Resp{Text: "Связываюсь с сервером... Передаю запрос... Жду ответ... Готово! Чтобы продолжить, скажите «Дальше».",
		TTS:     "\"Связываюсь с сервером... sil <[2500]> Передаю запрос... sil <[2500]> Жду ответ... sil <[2500]> Готово! Чтобы продолжить, скажите «Дальше».\"",
		Buttons: buttons},
		State: &aliceapi.StateData{ItemText: id, State: aliceapi.StateWaitAnswerFromChatGPT}}, nil
	/*questionName, _ := intent.Slots.ListName.AsString()
	questionName = "Расскажи про " + questionName + " подробно с примерами."
	client := openai.NewClient("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")
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: questionName,
				},
			},
		},
	)

	if err != nil {
		return &aliceapi.Response{Response: &aliceapi.Resp{Text: "Хм, это тема слишком сложная для меня, попробуй узнать у своих одногруппников или преподавателя"}}, nil
	}

	return &aliceapi.Response{Response: &aliceapi.Resp{Text: resp.Choices[0].Message.Content}}, nil*/
}

func (h *Handler) getAnswer(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	answerID := req.State.Session.ItemText
	if req.Request.Type == aliceapi.RequestTypeButton || (req.Request.Type == aliceapi.RequestTypeSimple && (req.Request.NLU.Intents.Next != nil || req.Request.NLU.Intents.ButtonNext != nil)) {
		// строка подключения
		dsn := config.Get().DatabaseEndpoint
		// IAM-токен
		token := config.Get().IAMToken
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
			return &aliceapi.Response{Response: &aliceapi.Resp{Text: "Хм, это тема слишком сложная для меня, попробуй узнать у своих одногруппников или преподавателя, либо попробуй узнать об этом попозже, когда я прокачаю свои скилы"}}, nil
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
			var buttons []*aliceapi.Button
			buttons = append(buttons, &aliceapi.Button{Title: "Дальше"})
			return &aliceapi.Response{Response: &aliceapi.Resp{Text: "Связываюсь с сервером... Передаю запрос... Жду ответ... Готово! Думаешь  почему опять, а потому что слишком быстро нажал кнопку готов в первый раз... Чтобы продолжить, скажите «Дальше».",
				TTS:     "\"Связываюсь с сервером... sil <[2500]> Передаю запрос... sil <[2500]> Жду ответ... sil <[2500]> Готово! Думаешь  почему опять, а потому что слишком быстро нажал кнопку готов в первый раз... Чтобы продолжить, скажите «Дальше».\"",
				Buttons: buttons},
				State: &aliceapi.StateData{ItemText: answerID, State: aliceapi.StateWaitAnswerFromChatGPT}}, nil
		}
		if len(answerMain) == 0 {
			var buttons []*aliceapi.Button
			buttons = append(buttons, &aliceapi.Button{Title: "Дальше"})
			return &aliceapi.Response{Response: &aliceapi.Resp{Text: "Связываюсь с сервером... Передаю запрос... Жду ответ... Готово! Думаешь  почему опять, а потому что слишком быстро нажал кнопку готов в первый раз... Чтобы продолжить, скажите «Дальше».",
				TTS:     "\"Связываюсь с сервером... sil <[2500]> Передаю запрос... sil <[2500]> Жду ответ... sil <[2500]> Готово! Думаешь  почему опять, а потому что слишком быстро нажал кнопку готов в первый раз... Чтобы продолжить, скажите «Дальше».\"",
				Buttons: buttons},
				State: &aliceapi.StateData{ItemText: answerID, State: aliceapi.StateWaitAnswerFromChatGPT}}, nil
		}
		return &aliceapi.Response{Response: &aliceapi.Resp{Text: answerMain}}, nil

	} else {
		return nil, nil
	}
}
func createAnswer(ctx context.Context, id string, questionName string) {
	questionName = "Расскажи про " + questionName + " подробно с примерами."
	client := openai.NewClient(config.Get().KeyFromChatGPT)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: questionName,
				},
			},
		},
	)
	answer := ""
	if err != nil {
		answer = "Вопрос слишком сложный для меня, соре"
	} else {
		answer = resp.Choices[0].Message.Content
	}
	// строка подключения
	dsn := config.Get().DatabaseEndpoint
	// IAM-токен
	token := config.Get().IAMToken
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
		return
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
		UPSERT INTO answers(id, answer) VALUES ($id, $answer);
      `,
				table.NewQueryParameters(
					table.ValueParam("$id", types.BytesValueFromString(id)),
					table.ValueParam("$answer", types.UTF8Value(answer)), // подстановка в условие запроса
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
}
