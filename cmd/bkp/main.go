package main

import (
	"os"

	"github.com/blang/semver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger

	minResticVersion semver.Version
)

func main() {
	setupLogger()
	defer logger.Sync()

	if forceRoot() {
		os.Exit(0)
	}

	// start handling commandline
	rootCmd := makeRootCmd()
	rootCmd.AddCommand(getEnvCmd(), getJobsCmd(), getSnapshotsCmd(), getMountCmd())
	rootCmd.Execute()
}

func init() {
	minResticVersion, _ = semver.Make("0.9.5")
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
