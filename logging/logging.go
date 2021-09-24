package logging

import (
	"io"
	"log"
	"os"
	"time"
)

const _INFO_PREFIX = "[*] "
const _ERROR_PREFIX = "[!] "

func InitLogs() {
	t := time.Now()
	logFile, err := os.OpenFile(t.Format("20060102150405")+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func Print(v ...interface{}) {
	log.Print(v...)
}

func Println(v ...interface{}) {
	log.Println(v...)
}

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Info(v ...interface{}) {
	log.Print(_INFO_PREFIX)
	log.Print(v...)
}

func Infoln(v ...interface{}) {
	log.Print(_INFO_PREFIX)
	log.Println(v...)
}

func Infof(format string, v ...interface{}) {
	log.Print(_INFO_PREFIX)
}

func Error(v ...interface{}) {
	log.Print(_ERROR_PREFIX)
	log.Print(v...)
}

func Errorln(v ...interface{}) {
	log.Print(_ERROR_PREFIX)
	log.Println(v...)
}

func Errorf(format string, v ...interface{}) {
	log.Print(_ERROR_PREFIX)
}
