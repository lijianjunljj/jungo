package jun_log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

// levels
const (
	debugLevel   = 0
	releaseLevel = 1
	errorLevel   = 2
	fatalLevel   = 3
)

const (
	printDebugLevel   = "[debug  ] "
	printReleaseLevel = "[release] "
	printErrorLevel   = "[error  ] "
	printFatalLevel   = "[fatal  ] "
)

type Logger struct {
	mu         sync.Mutex
	level      int
	baseLogger *log.Logger
	baseFile   *os.File

	// 当使用 JUN_LOG_TIME_FORMAT 时，按周期检测并切换文件（与 gocommon/logging 行为对齐）
	logDir       string
	rotateFormat string
	fileFlag     int
	currentPath  string
	rotateStop   chan struct{}
	wg           sync.WaitGroup
}

func openLogFileAppend(filePath string) (*os.File, error) {
	return os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func (logger *Logger) rotateIfNeeded() {
	tf := strings.TrimSpace(logger.rotateFormat)
	if tf == "" || logger.logDir == "" {
		return
	}

	logger.mu.Lock()
	defer logger.mu.Unlock()

	newPath := path.Join(logger.logDir, time.Now().Format(tf)+".log")
	if newPath == logger.currentPath && logger.baseFile != nil {
		return
	}

	f, err := openLogFileAppend(newPath)
	if err != nil {
		return
	}

	old := logger.baseFile
	logger.baseFile = f
	logger.baseLogger = log.New(f, "", logger.fileFlag)
	logger.currentPath = newPath

	if old != nil {
		_ = old.Close()
	}
}

func (logger *Logger) startPeriodicRotate() {
	logger.rotateStop = make(chan struct{})
	logger.wg.Add(1)
	go func() {
		defer logger.wg.Done()
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				logger.rotateIfNeeded()
			case <-logger.rotateStop:
				return
			}
		}
	}()
}

func New(strLevel string, pathname string, flag int) (*Logger, error) {
	// level
	var level int
	switch strings.ToLower(strLevel) {
	case "debug":
		level = debugLevel
	case "release":
		level = releaseLevel
	case "error":
		level = errorLevel
	case "fatal":
		level = fatalLevel
	default:
		return nil, errors.New("unknown level: " + strLevel)
	}

	// logger
	var baseLogger *log.Logger
	var baseFile *os.File
	logger := &Logger{
		level:    level,
		logDir:   pathname,
		fileFlag: flag,
	}

	if pathname != "" {
		now := time.Now()
		tf := strings.TrimSpace(os.Getenv("JUN_LOG_TIME_FORMAT"))

		if tf != "" {
			logger.rotateFormat = tf
			fullPath := path.Join(pathname, now.Format(tf)+".log")
			file, err := openLogFileAppend(fullPath)
			if err != nil {
				return nil, err
			}
			baseLogger = log.New(file, "", flag)
			baseFile = file
			logger.currentPath = fullPath
			logger.startPeriodicRotate()
		} else {
			filename := fmt.Sprintf("%d%02d%02d_%02d_%02d_%02d.log",
				now.Year(),
				now.Month(),
				now.Day(),
				now.Hour(),
				now.Minute(),
				now.Second())
			file, err := os.Create(path.Join(pathname, filename))
			if err != nil {
				return nil, err
			}
			baseLogger = log.New(file, "", flag)
			baseFile = file
		}
	} else {
		baseLogger = log.New(os.Stdout, "", flag)
	}

	logger.baseLogger = baseLogger
	logger.baseFile = baseFile

	return logger, nil
}

// It's dangerous to call the method on logging
func (logger *Logger) Close() {
	if logger.rotateStop != nil {
		close(logger.rotateStop)
		logger.wg.Wait()
		logger.rotateStop = nil
	}

	logger.mu.Lock()
	defer logger.mu.Unlock()

	if logger.baseFile != nil {
		_ = logger.baseFile.Close()
	}

	logger.baseLogger = nil
	logger.baseFile = nil
}

func (logger *Logger) doPrintf(level int, printLevel string, format string, a ...interface{}) {
	if level < logger.level {
		return
	}

	logger.mu.Lock()
	defer logger.mu.Unlock()

	if logger.baseLogger == nil {
		panic("logger closed")
	}

	format = printLevel + format
	logger.baseLogger.Output(3, fmt.Sprintf(format, a...))

	if level == fatalLevel {
		os.Exit(1)
	}
}

func (logger *Logger) Debug(format string, a ...interface{}) {
	logger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func (logger *Logger) Release(format string, a ...interface{}) {
	logger.doPrintf(releaseLevel, printReleaseLevel, format, a...)
}

func (logger *Logger) Error(format string, a ...interface{}) {
	logger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func (logger *Logger) Fatal(format string, a ...interface{}) {
	logger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

var gLogger, _ = New("debug", "", log.LstdFlags)

// It's dangerous to call the method on logging
func Export(logger *Logger) {
	if logger != nil {
		gLogger = logger
	}
}

func Debug(format string, a ...interface{}) {
	gLogger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func Release(format string, a ...interface{}) {
	gLogger.doPrintf(releaseLevel, printReleaseLevel, format, a...)
}

func Error(format string, a ...interface{}) {
	gLogger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func Fatal(format string, a ...interface{}) {
	gLogger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

func Close() {
	gLogger.Close()
}
