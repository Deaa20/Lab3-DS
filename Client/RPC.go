package main


type AddFingerEntryArgs struct {
	Port int
	Address string
	Status int
	
}
type AddFingerEntryReply struct {
	Status int
	
}
type FindSuccessorArgs struct {
	String string
}

type FindSuccessorReply struct {
	String string
}

