package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/go-playground/log/v7"
)

// file lock
type FileLock struct {
	dir string
	f   *os.File
}

func NewFileLock(dir string) *FileLock {
	return &FileLock{
		dir: dir,
	}
}

func (l *FileLock) Lock() error {
	f, err := os.Open(l.dir)
	if err != nil {
		return err
	}
	l.f = f
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		return fmt.Errorf("cannot flock directory %s - %s", l.dir, err)
	}
	return nil
}

func (l *FileLock) Unlock() error {
	defer l.f.Close()
	return syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
}

// ReadFile reads csv file
func ReadFile(filename string) ([][]string, error) {
	csvFile, err := os.Open(filename)
	if err != nil {
		log.Error("read file error: %v", err)
		return nil, err
	}
	defer csvFile.Close()
	ReadCsv := csv.NewReader(csvFile)

	stringValue, _ := ReadCsv.ReadAll()
	return stringValue, nil
}

// demo job
func writeCsvByLine(path string, dataStruct *WriteDataByLine) error {
	//todo: bugs might remain, need mutex
	err := flock.Lock()
	defer flock.Unlock()
	if err != nil {
		log.Error("file lock err: ", err.Error())
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Error(err)
	}
	defer file.Close()

	WriterCsv := csv.NewWriter(file)

	startTime := strconv.Itoa(int(dataStruct.StartTime))
	stopTime := strconv.Itoa(int(dataStruct.StopTime))
	taskId := fmt.Sprintf("%v", dataStruct.TaskId)
	duration := strconv.Itoa(int(dataStruct.Duration) / int(time.Second))
	res := int(dataStruct.StopTime) - int(dataStruct.StartTime)
	stringRes := strconv.Itoa(res)
	dataLine := []string{taskId, duration, startTime, stopTime, stringRes}

	if err := WriterCsv.Write(dataLine); err != nil {
		log.Error(err)
	}

	WriterCsv.Flush()
	return nil
}
