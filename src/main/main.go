package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-playground/log/v7"
	"github.com/go-playground/log/v7/handlers/console"
)

var a, p = StartPaxos()
var flock *FileLock

func main() {
	//initialize logs
	cLog := console.New(true)
	log.AddHandler(cLog, log.AllLevels...)
	defer log.WithTrace().Info("time to run")

	//delete old log file
	err := os.Remove(Filepath + logFilename)
	if err != nil {
		fmt.Println("file remove Error!")
		fmt.Printf("%s", err)
	} else {
		fmt.Print("file remove OK!")
	}

	//file lock
	flock = NewFileLock(Filepath + logFilename)

	timeWheel := CreateTimeWheel(1*time.Second, slotsNums)
	timeWheel.startTW()
	log.Info("initialize rpc")
	timeWheel.serverTW()
	log.Info("start Batch register")
	time.Sleep(time.Duration(5) * time.Second)
	BatchRegister(time.Now())
	//for _, val := range TW.slots {
	//	for item := val.Front(); item != nil; item = item.Next() {
	//		log.Info("item is: ", item.Value.(*Task))
	//	}
	//}

	defer func() {
		for {
		}
	}()
}

func TaskJob(key interface{}) {
	fmt.Println(fmt.Sprintf("%v This is a task job with key: %v", time.Now().Format(Format), key))
}

func test(timeWheel *TimeWheel) {
	fmt.Println(fmt.Sprintf("%v Add task task-5s", time.Now().Format(time.RFC3339)))
	args := &AddTaskArgs{time.Duration(1) * time.Second, 1, time.Now(), TaskJob}
	reply := AddTaskReply{}
	err := timeWheel.AddTask(args, &reply)
	fmt.Println("finish1")
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("%v Add task task-10s", time.Now().Format(time.RFC3339)))
	args = &AddTaskArgs{time.Duration(10) * time.Second, 2, time.Now(), TaskJob}
	reply = AddTaskReply{}
	err = timeWheel.AddTask(args, &reply)
	fmt.Println("finish2")
	if err != nil {
		panic(err)
	}
}
