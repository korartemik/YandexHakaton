package config

import (
	"encoding/base64"
	"fmt"
	"os"
)

type Config struct {
	IAMToken string

	KeyFromChatGPT   string
	DatabaseEndpoint string
}

func Get() *Config {
	return &Config{
		IAMToken:         "y0_AgAAAAAEqlnXAATuwQAAAADeYHM8MglXaQBDSvO684CBJ_RLs3H3ZyI",
		KeyFromChatGPT:   "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",
		DatabaseEndpoint: "grpcs://ydb.serverless.yandexcloud.net:2135/ru-central1/b1gsm4ottbrmg1pmmc79/etn0j1l74mga5f3vfeol",
	}
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
