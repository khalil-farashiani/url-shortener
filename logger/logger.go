package logger

import (
	"encoding/json"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	cfg    zap.Config
	Logger *zap.Logger
)

func init() {
	// get config from json file
	currentDir, _ := os.Getwd()
	file, err := os.Open(currentDir + "/logger/zap.config")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	cfg.EncoderConfig.EncodeTime = zapcore.EpochTimeEncoder
	if err := json.Unmarshal([]byte(reader), &cfg); err != nil {
		panic(err)
	}
	Logger, err = cfg.Build()
	if err != nil {
		panic(err)
	}
	defer Logger.Sync()
}
