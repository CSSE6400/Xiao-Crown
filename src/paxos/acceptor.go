package paxos

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"

	"github.com/go-playground/log/v7"
)

type Acceptor struct {
	lis net.Listener
	// node id
	id int
	// promise number, if 0 indicate that node haven't recive any Prepare yet.
	promiseNumber int
	// accepted number, if 0 indicate that node haven't accept any Proposal yet.
	acceptedNumber int
	// nil if no accepted value
	acceptedValue interface{}
	// learners id
	learners []int
	mutex    sync.Mutex
}

func newAcceptor(id int) *Acceptor {
	acceptor := &Acceptor{
		id: id,
	}
	acceptor.server()
	return acceptor
}

func (a *Acceptor) LockAcceptor(args *PaxosMsgArgs, reply *PaxosMsgReply) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return nil
}

func (a *Acceptor) Prepare(args *PaxosMsgArgs, reply *PaxosMsgReply) error {
	a.mutex.Lock()
	fmt.Println("Prepare from ", args.From, " to ", args.To)
	fmt.Println("args.num ", args.Number, "a.promise ", a.promiseNumber)
	if args.Number > a.promiseNumber {
		a.promiseNumber = args.Number
		fmt.Println("prepare promiseNumber ", a.promiseNumber)
		reply.Number = a.acceptedNumber
		fmt.Println("prepare accepted number ", a.acceptedNumber)
		reply.Value = a.acceptedValue
		fmt.Println("prepare acceptedValue ", a.acceptedValue)
		reply.Ok = true
	} else {
		reply.Ok = false
	}
	return nil
}

func (a *Acceptor) Accept(args *PaxosMsgArgs, reply *PaxosMsgReply) error {
	defer a.mutex.Unlock()
	fmt.Println("Accept from ", args.From, " to ", args.To)
	fmt.Println("args.num ", args.Number, "a.promise ", a.promiseNumber)
	if args.Number >= a.promiseNumber {
		a.promiseNumber = args.Number
		a.acceptedNumber = args.Number
		a.acceptedValue = args.Value
		reply.Ok = true

		fmt.Println("accept value: ", args.Value)
		_, existed := TW.finishedTasks.Load(args.Value.TaskId)
		if !existed {
			go WriteToMap(args.Value.TaskId)
		}
	} else {
		reply.Ok = false
	}
	//clean the acceptor
	a.promiseNumber = 0
	a.acceptedValue = nil
	a.acceptedNumber = 0
	return nil
}

// run tcp in passive mode, accept connection from others.
func (a *Acceptor) server() {
	rpcs := rpc.NewServer()
	rpcs.Register(a)
	addr := fmt.Sprintf(":6%d", a.id)
	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("listen error 3:", e)
	}
	a.lis = l
	go func() {
		for {
			conn, err := a.lis.Accept()
			if err != nil {
				continue
			}
			go rpcs.ServeConn(conn)
		}
	}()
}

// close connection
func (a *Acceptor) close() {
	a.lis.Close()
}
