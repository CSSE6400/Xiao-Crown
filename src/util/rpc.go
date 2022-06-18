package util

import (
	"fmt"
	"net/rpc"

	"github.com/go-playground/log/v7"
)

type PaxosMsgArgs struct {
	Number int         // proposal turn
	Value  interface{} // proposal value
	From   int         // sender id
	To     int         // reciver id
}

type PaxosMsgReply struct {
	Ok     bool
	Number int
	Value  interface{}
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
