package logging

import (
	"io"
	"log"
	"os"
	"time"
)

const _ERROR_PREFIX string = "[ERROR] "
const _DEBUG_PREFIX string = "[DEBUG] "
const _INFO_PREFIX string = "[INFO] "

func InitLogs() {
	t := time.Now()
	logFile, err := os.OpenFile(t.Format("20060102150405")+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func Print(v ...interface{}) {
	Info(v...)
}

func Println(v ...interface{}) {
	Infoln(v...)
}

func Printf(format string, v ...interface{}) {
	Infof(format, v...)
}

func Info(v ...interface{}) {
	log.SetPrefix(_INFO_PREFIX)
	log.Print(v...)
}

func Infoln(v ...interface{}) {
	log.SetPrefix(_INFO_PREFIX)
	log.Println(v...)
}

func Infof(format string, v ...interface{}) {
	log.SetPrefix(_INFO_PREFIX)
	log.Printf(format, v...)
}

func Debug(v ...interface{}) {
	log.SetPrefix(_DEBUG_PREFIX)
	log.Print(v...)
}

func Debugln(v ...interface{}) {
	log.SetPrefix(_DEBUG_PREFIX)
	log.Println(v...)
}

func Debugf(format string, v ...interface{}) {
	log.SetPrefix(_DEBUG_PREFIX)
	log.Printf(format, v...)
}

func Error(v ...interface{}) {
	log.SetPrefix(_ERROR_PREFIX)
	log.Print(v...)
}

func Errorln(v ...interface{}) {
	log.SetPrefix(_ERROR_PREFIX)
	log.Println(v...)
}

func Errorf(format string, v ...interface{}) {
	log.SetPrefix(_ERROR_PREFIX)
	log.Printf(format, v...)
}
