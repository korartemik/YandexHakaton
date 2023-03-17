package config

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Config struct {
	OAuthToken            string
	KeyFromChatGPT        string
	DatabaseEndpoint      string
	FunctionEndpoint      string
	FunctionQandAEndpoint string
	IAMToken              string
}

func Get() *Config {
	return &Config{
		OAuthToken:            requireString("OAUTH_TOKEN"),
		KeyFromChatGPT:        requireString("KEY_CHAT_GPT"),
		DatabaseEndpoint:      requireString("DATABASE_ENDPOINT"),
		FunctionEndpoint:      requireString("FUNCTION_ENDPOINT"),
		FunctionQandAEndpoint: requireString("FUNCTION_QA_ENDPOINT"),
		IAMToken:              "",
	}
}

func NewConfig() *Config {
	c := Get()
	c.createIAMToken()
	return c
}

func requireBytes(name string) []byte {
	str := requireString(name)
	res, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		panic(fmt.Errorf("failed to decode var %v with base64: %w", err))
	}
	return res
}
func requireString(name string) string {
	res, ok := os.LookupEnv(name)
	if !ok {
		panic(fmt.Sprintf("required env var %s not found", name))
	}
	return res
}

func (c *Config) createIAMToken() {
	values := map[string]string{"yandexPassportOauthToken": c.OAuthToken}

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
	c.IAMToken = data.IAMToken
	log.Printf("IAMToken create succses")
}
