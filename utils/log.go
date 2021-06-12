package utils

import (
	"log"
	"os"
	"time"
)

type Log struct {
	LogDir     string
	TimeFormat string
}

func (obj *Log) Write(f *os.File, text string) {

	t := time.Now()
	var pre = "[" + t.Format(obj.TimeFormat) + "]==>"
	_, err := f.WriteString(pre + text + "\n")
	log.Println(text)
	obj.CheckError(err)
}

func (obj *Log) CreateFile(fileName string) *os.File {

	err := os.MkdirAll(obj.LogDir, os.ModePerm)
	obj.CheckError(err)

	file, err := os.Create(obj.LogDir + fileName)
	obj.CheckError(err)

	return file
}

func (obj *Log) CheckError(err error) {
	if err != nil {
		log.Panic(err)
		panic(err)
	}
}
