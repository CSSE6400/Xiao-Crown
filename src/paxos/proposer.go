package paxos

import (
	"strconv"
	"sync"

	"github.com/go-playground/log/v7"
)

type Proposer struct {
	// node id
	id int
	// Max round that currently known(learned) by node
	round int
	// proposal number
	number int
	// list of acceptor
	acceptors []int

	mutex sync.Mutex
}

func (p *Proposer) Propose(v interface{}) interface{} {
	p.round = p.round + numMachines
	p.number = p.round<<16 | p.id
	log.Info(" p.round: ", p.round, " p.number: ", p.number)

	// phase 1
	prepareCount := 0
	maxNumber := 0
	for _, aid := range p.acceptors {
		args := PaxosMsgArgs{
			Number: p.number,
			From:   p.id,
			To:     aid,
			Value:  v,
		}
		reply := PaxosMsgReply{}

		log.Info("before prepare, and the value is: ", args.Value)
		err := call(SocketNames[aid], "Acceptor.Prepare", args, &reply, strconv.Itoa(aid+60))
		log.Info("after prepare, and the value is: ", args.Value)
		if !err {
			continue
		} else {
			log.Info("Acceptor.Prepare:", err)
		}

		if reply.Ok {
			prepareCount++
			if reply.Number > maxNumber {
				maxNumber = reply.Number
				v = reply.Value
			}
		}

		if prepareCount == p.majority() {
			break
		}
	}

	// phase 2
	acceptCount := 0
	if prepareCount >= p.majority() {
		log.Info("phase two before for loop")
		for _, aid := range p.acceptors {
			log.Info("aid is ", aid)
			args := PaxosMsgArgs{
				Number: p.number,
				Value:  v,
				From:   p.id,
				To:     aid,
			}
			reply := PaxosMsgReply{}
			//todo change the address
			log.Info("before accept, and the value is: ", args.Value)
			err := call(SocketNames[aid], "Acceptor.Accept", args, &reply, strconv.Itoa(aid+60))
			log.Info("after accept, and the value is: ", args.Value)
			if !err {
				continue
			} else {
				log.Info("Acceptor.Accept:", err)
			}

			if reply.Ok {
				acceptCount++
			}
		}
	}

	if acceptCount >= p.majority() {
		// todo save locally
		return v
	}
	log.Info("after accept")
	return nil
}

func (p *Proposer) majority() int {
	return len(p.acceptors)/2 + 1
}
