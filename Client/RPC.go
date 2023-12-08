package main

type AddFingerEntryArgs struct {
	Port    int
	Address string
	Status  int
}
type AddFingerEntryReply struct {
	Status int
}

type FindSuccessorArgs struct {
	NodeAddress string
	NodePort    string
}

type FindSuccessorReply struct {
	Successor string
	Final     bool
}
