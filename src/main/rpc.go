package main

import (
	"fmt"
	"net"
	"net/rpc"
	"time"

	"github.com/go-playground/log/v7"
)

type AddTaskArgs struct {
	interval time.Duration
	uuid     int
	execTime time.Time
	taskJob  interface{}
}

type AddTaskReply struct {
	stupid int
}

type BackupArgs struct {
	Interval time.Duration
	Uuid     int
}

type BackupReply struct {
	Msg  string
	Code int
}

type PaxosMsgArgs struct {
	Number int              // proposal turn
	Value  *WriteDataByLine // proposal value
	From   int              // sender id
	To     int              // reciver id
}

type PaxosMsgReply struct {
	Ok     bool
	Number int
	Value  *WriteDataByLine
}

// This method starts a RPC server for tw
func (tw *TimeWheel) serverTW() {
	rpcs := rpc.NewServer()
	rpcs.Register(tw)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port1))
	if err != nil {
		log.Error("listen error 1:", err)
	}
	tw.lis = lis
	go func() {
		for {
			conn, err := tw.lis.Accept()
			if err != nil {
				continue
			}
			go rpcs.ServeConn(conn)
		}
	}()
}

//
// send an RPC request to other servers, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(sockname string, rpcname string, args interface{}, reply interface{}, port string) bool {
	fmt.Println(fmt.Sprintf("%s:%s", sockname, port))
	c, err := rpc.Dial("tcp", fmt.Sprintf("%s:%s", sockname, port))
	fmt.Println("client is: ", c)
	if err != nil {
		log.Error("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Printf("calling:%s::%s() error: %v\n", sockname, rpcname, err)
	log.Error("rpc call err: ", err)

	return false
}
