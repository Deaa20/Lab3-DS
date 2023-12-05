package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
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
	fingerTable := make([]string, 160)
	for i := range fingerTable {
		fingerTable[i] = Node.Address + strconv.Itoa(Node.Port)
	}

	//Successor of node is the node itself
	Node.Predecessor = ""
	Node.FingerTable = fingerTable
	Node.Successors = successors
	go Node.Server(Node.Port)
}

func JoinChord() {
	Node.Predecessor = ""
	args := ExampleArgs{}
	args.String = "I have joined the Chord"
	args.Number = 1
	reply := ExampleReply{}
	fmt.Println(args.String+ " im here at joinChord")

	ok := callNode("NodeClient.SendTest", &args, &reply, Node.JoinAddress, Node.JoinPort)



	//ok := callNode("NodeClient.SendTest", &args, &reply, Node.JoinAddress, Node.JoinPort)
	if ok {
		fmt.Println("Chord is working")
		fmt.Println(reply.String)
		fmt.Println(reply.Number)
		//Node.ReciveTest(&args, &reply)

	} else {
		fmt.Print("JoinChord is not OK")
	}


	// Fund successor of node
	argsFindSuccessor := FindSuccessorArgs{}
	replyFindSuccessor := FindSuccessorReply{}
	ok = callNode("NodeClient.FindSuccessor", &argsFindSuccessor, &replyFindSuccessor, Node.JoinAddress, Node.JoinPort)

	if ok {
		fmt.Println("Found successor")
		fmt.Println(replyFindSuccessor.Successor)
		Node.Successors[0] = replyFindSuccessor.Successor

	} else {
		fmt.Print("Error when looking for successor")
	}

	// Notify successor of node that the node is the predecessor
	ok = callNode("NodeClient.Notify", &args, &reply, Node.Successors[0], Node.JoinPort)
	if !ok {
		fmt.Print("Error when notifying successor")
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
	return false
}

func (n *NodeClient) Server(port int) {
	fmt.Print("Server is running \n " )

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

func (n *NodeClient) SendTest(args *ExampleArgs, reply *ExampleReply) error {
	
	reply.String = args.String
	reply.Number = args.Number + 10
	return nil
}


func (n *NodeClient) ReciveTest(args *ExampleArgs, reply *ExampleReply) error {
	fmt.Print(reply.String)
	return nil
}

func (n *NodeClient) FindSuccessor(args *FindSuccessorArgs, reply *FindSuccessorReply) error {
	reply.Successor = "I am the successor"
	return nil
}

func LeaveChord() {

}
func Lookup() {}

func StoreFile() {

}

func PrintState() {
}
