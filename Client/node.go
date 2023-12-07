package main


import (
	"fmt"
	"log"
	"math/big"
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
	ClientID                 *big.Int
	Lock                     sync.Mutex
	Status                   int
}


var Node NodeClient

type NodeInfo struct {
	Address string
	Port    int
}

func CreateNode(address string, port int, joinAddress string, joinPort int, stabilizeInterval int,
	fixFingersInterval int, checkPredecessorInterval int, numSuccessors int, clientID string) NodeClient {

	Node = NodeClient{
		Address:                  address,
		Port:                     port,
		JoinAddress:              joinAddress,
		JoinPort:                 joinPort,
		StabilizeInterval:        stabilizeInterval,
		FixFingersInterval:       fixFingersInterval,
		CheckPredecessorInterval: checkPredecessorInterval,
		NumSuccessors:            numSuccessors,
		ClientID:                 HashString(clientID),
		Status:                   1,
	}

	if Node.JoinAddress == "" {
		NewChord()
	} else {
		JoinChord()
	}
	return Node

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
	args := FindSuccessorArgs{}
	args.String = "join"
	reply := FindSuccessorReply{}
	fmt.Print(args.String + " im here at joinChord")

	ok := callNode("NodeClient.SendTest", &args, &reply, Node.JoinAddress, Node.JoinPort)

	//ok := callNode("NodeClient.SendTest", &args, &reply, Node.JoinAddress, Node.JoinPort)
	if ok {
		fmt.Print("Chord is working \n")
		fmt.Print(reply.String)
		//Node.ReciveTest(&args, &reply)

	} else {
		fmt.Print("JoinChord is not OK")
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
	fmt.Print("Server is running \n ")

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

	reply.String = args.String
	return nil
}
func (n *NodeClient) ReciveTest(args *FindSuccessorArgs, reply *FindSuccessorReply) error {
	fmt.Print(reply.String)
	return nil
}

func LeaveChord() {

}
func Lookup() {}

func StoreFile() {

}

func PrintState(node *NodeClient) {
	fmt.Printf("Printing the state of the ccurrent node\n ")
	fmt.Printf("Address:" + node.Address + "\n ")
	fmt.Printf("port:" + fmt.Sprint(node.Port) + "\n ")
	fmt.Printf("predecessor:" + node.Predecessor + "\n ")
	fmt.Printf("JoinAddress:" + node.JoinAddress + "\n ")
	fmt.Printf("JoinPort:" + fmt.Sprint(node.JoinPort) + "\n ")
	fmt.Printf("StabilizeInterval:" + fmt.Sprint(node.StabilizeInterval) + "\n ")
	fmt.Printf("FixFingersInterval:" + fmt.Sprint(node.FixFingersInterval) + "\n ")
	fmt.Printf("CheckPredecessorInterval:" + fmt.Sprint(node.CheckPredecessorInterval) + "\n ")
	fmt.Printf("NumSuccessors:" + fmt.Sprint(node.NumSuccessors) + "\n ")
	fmt.Printf("ClientID:" + fmt.Sprintf("%040x", (node.ClientID)) + "\n ")
}
