/*
 * Codes to bootstrap paxos
 */
package main

func start(acceptorIds []int) []*Acceptor {
	acceptors := make([]*Acceptor, 0)
	for _, aid := range listenIds {
		a := newAcceptor(aid)
		acceptors = append(acceptors, a)
	}

	return acceptors
}

func cleanup(acceptors []*Acceptor) {
	for _, a := range acceptors {
		a.close()
	}
}

func StartPaxos() ([]*Acceptor, *Proposer) {
	acceptors := start(AcceptorIds)
	// defer cleanup(acceptors) do this somewhere eles
	p := &Proposer{
		id:        proposerID,
		acceptors: AcceptorIds,
	}

	return acceptors, p
}
