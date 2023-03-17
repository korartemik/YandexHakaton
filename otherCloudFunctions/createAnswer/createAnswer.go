package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

type RequestBody struct {
	HttpMethod string `json:"httpMethod"`
	Body       []byte `json:"body"`
}

type Body struct {
	ID        string `json:"idTag"`
	Question  string `json:"questionName"`
	IdRequest string `json:"idRequest"`
}

type Config struct {
	IAMToken string

	KeyFromChatGPT   string
	DatabaseEndpoint string
	OAthToken        string
}

func GetConfig() *Config {
	return &Config{
		IAMToken:         "",
		KeyFromChatGPT:   os.Getenv("KEY_FROM_CHAT_GPT"),
		DatabaseEndpoint: os.Getenv("DATABASE_ENDPOINT"),
		OAthToken:        os.Getenv("OAUTH_TOKEN"),
	}
}

var config = GetConfig()

func Handler(ctx context.Context, request []byte) (*Response, error) {
	requestBody := &RequestBody{}
	// Массив байтов, содержащий тело запроса, преобразуется в соответствующий объект
	err := json.Unmarshal(request, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("an error has occurred when parsing request: %v", err)
	}

	// В журнале будет напечатано название HTTP-метода, с помощью которого осуществлен запрос, а также тело запроса
	fmt.Println(requestBody.HttpMethod, string(requestBody.Body))

	req := &Body{}
	// Поле body запроса преобразуется в объект типа Request для получения переданного имени
	err = json.Unmarshal(requestBody.Body, &req)
	if err != nil {
		return nil, fmt.Errorf("an error has occurred when parsing body: %v", err)
	}
	id := req.ID
	questionName := req.Question
	if req.IdRequest == "1" {
		questionName = "Расскажи про " + questionName + " подробно с примерами."
	}
	answer := getAnswer(questionName)
	log.Printf("Answer generated")

	createToken()

	log.Printf("Token generated succses")
	// строка подключения
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
		return &Response{
			StatusCode: 500,
			Body:       err,
		}, nil
	}
	// закрытие драйвера по окончании работы программы обязательно
	defer db.Close(ctx)

	log.Printf("Connection succses")

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
	if err != nil {
		return &Response{
			StatusCode: 500,
			Body:       err,
		}, nil
	}
	return &Response{
		StatusCode: 200,
		Body:       "Ok",
	}, nil
}

func getAnswer(questionName string) string {
	client := openai.NewClient(GetConfig().KeyFromChatGPT)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()
	resp, err := client.CreateChatCompletion(
		ctx,
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
		answer = "Вопрос слишком сложный для меня, либо слшиком много запросов и мне не хватило времени, попробуй задать его попозже и я постараюсь на него ответить"
	} else {
		answer = resp.Choices[0].Message.Content
	}
	for strings.Index(answer, "\n") == 0 {
		answer = answer[1:]
	}
	cancel()
	return answer
}

func createToken() {
	values := map[string]string{"yandexPassportOauthToken": config.OAthToken}

	jsonValue, _ := json.Marshal(values)
	fmt.Println(string(jsonValue))
	resp, _ := http.Post("https://iam.api.cloud.yandex.net/iam/v1/tokens", "application/json", bytes.NewBuffer(jsonValue))
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Error with generating token %s", resp.Status))
	}
	var data struct {
		IAMToken string `json:"iamToken"`
	}
	err := json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("Error with generating token")
		panic(err)
	}
	config.IAMToken = data.IAMToken
}
