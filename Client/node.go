package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

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
	Lock                     sync.Mutex
	Status                   int
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
	go Node.Server(Node.Port)
}

func JoinChord() {
	Node.Predecessor = ""
	args := FindSuccessorArgs{String: "join"}
	reply := FindSuccessorReply{}

	ok := callNode("NodeClient.SendTest", &args, &reply, Node.JoinAddress, Node.JoinPort)

	//ok := callNode("NodeClient.SendTest", &args, &reply, Node.JoinAddress, Node.JoinPort)
	if ok {
		Node.ReciveTest(&args, &reply)

	} else {
		fmt.Print("here i am")
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

	fmt.Println(err)
	return false
}

func (n *NodeClient) Server(port int) {
	rpc.Register(n)
	rpc.HandleHTTP()
	portString := ":" + fmt.Sprint(port)

	fmt.Print(portString + "...\n")
	l, e := net.Listen("tcp", portString)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)

}

func (n *NodeClient) Done() bool {
	n.Lock.Lock()
	defer n.Lock.Unlock()
	ret := false
	if n.Status == 0 {
		ret = true
	}
	return ret
}

func (n *NodeClient) SendTest(args *FindSuccessorArgs, reply *FindSuccessorReply) error {
	fmt.Print("checking if im here or not")

	reply.String = "it works mf"
	fmt.Print("it works mf-------------------------")
	return nil
}
func (n *NodeClient) ReciveTest(args *FindSuccessorArgs, reply *FindSuccessorReply) error {
	fmt.Print(reply.String)
	fmt.Print("it works mf-------------------------")
	return nil
}

func LeaveChord() {

}
func Lookup() {}

func StoreFile() {

}

func PrintState() {
}
