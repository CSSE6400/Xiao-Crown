package main

import (
	"time"
)

const Format string = "2006/1/2 15:04:05"

var readFilename = "1.csv"
var Filepath = "."

//var Filepath = "/root/ft_local/Distributed-Timer-Service/src/test/"
var logFilename = "log1.csv"

const (
	serverID0   = 0
	serverID1   = 1
	serverID2   = 2
	serverID3   = 3
	serverID4   = 4
	proposerID  = 1
	numMachines = 3
	slotsNums   = 290
)

var listenIds = []int{serverID0, serverID1, serverID2, serverID3, serverID4}
var AcceptorIds = []int{serverID2}
var LearnerIds = []int{serverID4}

var registerIds = []int{serverID1, serverID3}

var port1 = "8006"

const (
	Socketname0 string = "9.134.131.104" //mVUFd@2873tB
	//Socketname1 string = "9.134.81.172"  //BdKJg@3697Uz //windows
	Socketname2 string = "9.134.72.227"  //PMqpN@5628eJ
	Socketname3 string = "9.135.113.126" //wpgqs*9728Jn
	Socketname4 string = "9.134.167.39"  //dpksZ*5439Cp
)

var SocketNames = []string{Socketname0, Socketname2, Socketname2, Socketname3, Socketname4}

//local server name
var LocalName = Socketname0

// WriteDataByLine
type WriteDataByLine struct {
	TaskId    interface{}
	Duration  time.Duration
	StartTime int64
	StopTime  int64
}

type writeData struct {
	writeData []WriteDataByLine
}
