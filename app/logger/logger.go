package logger

import (
	"fmt"

	"github.com/kainobor/eth-client/app/args"
	"github.com/kainobor/eth-client/app/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Logger is implementation of ILogger
	Logger struct {
		*zap.SugaredLogger
		env    string
		config *config.LoggerConfig
	}
)

var (
	getSyncerForPaths = func(paths []string) (zapcore.WriteSyncer, error) {
		syncer, _, err := zap.Open(paths...)

		return syncer, err
	}
)

// New logger
func New() *Logger {
	return new(Logger)
}

// Init logger through zap
func (l *Logger) Init(env string, c *config.LoggerConfig) error {
	l.config = c
	l.env = env

	highPriorityFunc, lowPriorityFunc := l.getLevelCheckers()

	highSyncer, lowSyncer, err := l.getSyncers()
	if err != nil {
		return fmt.Errorf("can't get syncers: %v", err)
	}

	encoder := l.getEncoder(env)
	core := l.createLoggerCore(highPriorityFunc, lowPriorityFunc, highSyncer, lowSyncer, encoder)

	l.initLogger(core)

	return nil
}

// getLevelCheckers returns two callbacks, that decide which messages should be set as high priority, which as low
func (l *Logger) getLevelCheckers() (zap.LevelEnablerFunc, zap.LevelEnablerFunc) {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	return highPriority, lowPriority
}

// getSyncers returns file syncers for each priority level files
func (l *Logger) getSyncers() (zapcore.WriteSyncer, zapcore.WriteSyncer, error) {
	errSync, err := getSyncerForPaths(l.config.ErrPaths)
	if err != nil {
		return nil, nil, fmt.Errorf("error while trying to open error paths: %v", err)
	}

	logSync, err := getSyncerForPaths(l.config.InfoPaths)
	if err != nil {
		return nil, nil, fmt.Errorf("error while trying to open log paths: %v", err)
	}

	return errSync, logSync, nil
}

// getEncoder returns encoder for logger messages
func (l *Logger) getEncoder(env string) zapcore.Encoder {
	var encoder zapcore.Encoder
	if env == args.EnvDev {
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	} else {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	}

	return encoder
}

func (l *Logger) createLoggerCore(
	highPriority zap.LevelEnablerFunc,
	lowPriority zap.LevelEnablerFunc,
	highSyncer zapcore.WriteSyncer,
	lowSyncer zapcore.WriteSyncer,
	encoder zapcore.Encoder,
) zapcore.Core {
	highCore := l.createCore(encoder, highSyncer, highPriority)
	lowCore := l.createCore(encoder, lowSyncer, lowPriority)

	return l.mergeCores(highCore, lowCore)
}

func (l *Logger) createCore(encoder zapcore.Encoder, syncer zapcore.WriteSyncer, priorityFunc zap.LevelEnablerFunc) zapcore.Core {
	return zapcore.NewCore(encoder, syncer, priorityFunc)
}

func (l *Logger) mergeCores(core1 zapcore.Core, core2 zapcore.Core) zapcore.Core {
	return zapcore.NewTee(core1, core2)
}

func (l *Logger) initLogger(core zapcore.Core) {
	parentLogger := zap.New(core)
	l.SugaredLogger = parentLogger.Sugar()
}
