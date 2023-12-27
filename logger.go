package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var lock = &sync.Mutex{}
var logger *LoggerType

type AllLoggers struct {
	Info  *log.Logger
	Warn  *log.Logger
	Err   *log.Logger
	Fatal *log.Logger
}
type LoggerType struct {
	Info  *log.Logger
	Warn  *log.Logger
	Err   *log.Logger
	Fatal *log.Logger
}

func GetLoggerType() AllLoggers {
	var allloggers AllLoggers
	homedir, _ := os.UserHomeDir()
	//@todo check if they exist first
	Infolog, err := os.OpenFile(homedir+"/.flowcat/logs/info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		fmt.Println("info log file not created", err.Error())
	}

	Errorlog, err := os.OpenFile(homedir+"/.flowcat/logs/error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		fmt.Println("error log file not created", err.Error())
	}

	lock.Lock()
	defer lock.Unlock()
	if logger == nil {
		logger = &LoggerType{}
	}
	logger.Info = log.New(Infolog, "[info]", log.LstdFlags|log.Lshortfile|log.Ldate|log.Ltime|log.LUTC)
	logger.Warn = log.New(Infolog, "[warn]", log.LstdFlags|log.Lshortfile|log.Ldate|log.Ltime|log.LUTC)
	logger.Err = log.New(Errorlog, "[ERROR]", log.LstdFlags|log.Lshortfile|log.Ldate|log.Ltime|log.LUTC)
	logger.Fatal = log.New(Errorlog, "[FATAL]", log.LstdFlags|log.Lshortfile|log.Ldate|log.Ltime|log.LUTC)

	allloggers.Info = logger.Info
	allloggers.Warn = logger.Warn
	allloggers.Err = logger.Err
	allloggers.Fatal = logger.Fatal

	return allloggers
}
