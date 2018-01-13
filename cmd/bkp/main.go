package main

import (
	"os"
	"runtime"

	"github.com/jojomi/bkp"
	script "github.com/jojomi/go-script"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

func main() {
	setupLogger()
	defer logger.Sync()

	if forceRoot() {
		os.Exit(0)
	}

	err := bkp.CheckEnvironment()
	if err != nil {
		sugar.Fatal(err)
	}

	// warn about nice (Linux, MacOS X) and ionice (Linux)
	sc := script.NewContext()
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		if !sc.CommandExists("nice") {
			sugar.Warn("\"nice\" command not found. Please make sure it is in your PATH to keep your system responsive while doing backups.")
		}
	}
	if runtime.GOOS == "linux" {
		if !sc.CommandExists("ionice") {
			sugar.Warn("\"ionice\" command not found. Please make sure it is in your PATH to keep your system responsive while doing backups.")
		}
	}

	// start handling commandline
	rootCmd := makeRootCmd()
	rootCmd.AddCommand(getJobsCmd())
	rootCmd.Execute()
}

func setupLogger() {
	encoderCfg := zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:       false,
		DisableStacktrace: true,
		DisableCaller:     true,
		Encoding:          "console",
		EncoderConfig:     encoderCfg,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
	}
	logger, _ = config.Build()
	sugar = logger.Sugar()
}
