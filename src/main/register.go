package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/log/v7"
)

// receive tasks from Client, parameter from client might be string, need parsing interval and uuid from it.
func Register(interval time.Duration, uuid int) {
	register(interval, uuid)
	//rpc to other servers
	// args := BackupArgs{
	// 	Interval: interval,
	// 	Uuid:     uuid,
	// }
	// backup(args)
}

// scan the csv
func BatchRegister(timeStart time.Time) {
	result, err := ReadFile(Filepath + readFilename)
	if err != nil {
		log.Error("err in read file %v", err)
	}
	//for each line in csv data structure:
	for i, items := range result {
		log.Info("the ", i, " task is registered since ", time.Since(timeStart))
		d, err := strconv.Atoi(items[1])
		if err != nil {
			log.Error(err)
		}
		duration := time.Duration(d) * time.Second
		uuid, err := strconv.Atoi(items[0])
		if err != nil {
			log.Error(err)
		}
		log.Info("duration: ", duration, " uuid: ", uuid)
		Register(duration, uuid)
	}
}

// rpc calling, public method
// finish tasks from other server
func (tw *TimeWheel) Backup(args *BackupArgs, reply *BackupReply) error {
	uuid := args.Uuid
	interval := args.Interval
	register(interval, uuid)
	reply.Msg = "Backup Succeed.\n"
	reply.Code = 0
	return nil
}

// register tasks
func register(interval time.Duration, uuid int) {
	log.Info(fmt.Sprintf("%s Add Task ID: %d", time.Now().Format(Format), uuid))
	args := &AddTaskArgs{interval, uuid, time.Now(), TaskJob}
	reply := AddTaskReply{}
	err := TW.AddTask(args, &reply)
	if err != nil {
		log.Error(err)
	}
}

func backup(args BackupArgs) {
	reply := BackupReply{}
	//fmt.Println(SocketNames[registerIds[0]])
	if ok := call(Socketname2, "TimeWheel.Backup", args, &reply, port1); !ok {
		//if ok := call(SocketNames[registerIds[0]], "Backup", args, &reply); !ok {
		log.Error("call not ok ", ok)
	}
	//if ok := call(SocketNames[registerIds[1]], "Backup", args, &reply); !ok {
	//	fmt.Printf("Register: backup register error\n")
	//}
}
