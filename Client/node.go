package main

import (
    "fmt"
    "log"
	"net/rpc")

type NodeClient struct {
	Address                  string
	Port                     int
	FingerTable              []string
	Predecessor              string
	Successors               []string
	Bucket                   map[Key]string
	JoinAddress              string
	JoinPort                 int
	StabilizeInterval        int
	FixFingersInterval       int
	CheckPredecessorInterval int
	NumSuccessors            int
	ClientID                 string
}

var Node NodeClient

type NodeInfo struct {
	Address string
	Port    int
}


func NewChord() {
	successors := make([]string, Node.NumSuccessors)
	for i := range successors {
		successors[i] = Node.Address
	}
	fingerTable := make([]NodeInfo, 160)
	for i := range fingerTable {
		fingerTable[i] = NodeInfo{Address: Node.Address, Port: Node.Port}
	}
	Node.Predecessor = ""
}

func JoinChord() {
	Node.Predecessor = ""
	args := FindSuccessorArgs{}
	reply := FindSuccessorReply{}
	ok := callNode("main.FindSuccessor", args, reply, Node.JoinAddress, Node.JoinPort)
	if ok {

	}

}

func callNode(rpcname string, args interface{}, reply interface{}, address string, port int) bool {
	//Check if Sprint works
	addressPort := address + ":" + fmt.Sprint(port)
	c, err := rpc.DialHTTP("tcp", addressPort)

	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	//fmt.Println(err)
	return false
}









func LeaveChord() {

}
func Lookup() {}

func StoreFile() {

}

func PrintState() {
}
