package main

import "math/big"

type AddFingerEntryArgs struct {
	Port    int
	Address string
	Status  int
}
type AddFingerEntryReply struct {
	Status int
}

type FindSuccessorArgs struct {
	ID *big.Int
}

type FindSuccessorReply struct {
	Successor NodeInfo
	Final     bool
}

type GetPredecessorArgs struct {
}
type GetPredecessorReply struct {
	Predecessor NodeInfo
	Successors  []NodeInfo
}

type NotifyNodesArgs struct {
	NodeAdress NodeInfo
}

type NotifyNodesReply struct {
	Successors []NodeInfo
	Check      bool
}

type PingArgs struct {
}
type PingReply struct {
}
