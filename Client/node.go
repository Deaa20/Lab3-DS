package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

//status 0 means that the node is leaving the  chord ring
//status 1 means that the node is joining the Chord ring
//status 2 means that the node is in the chord ring

const fingerTableSize int = 6

type NodeClient struct {
	Address                  string
	Port                     int
	FingerTable              [6]string
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
	args := AddFingerEntryArgs{}
	args.Port = Node.Port
	args.Address = Node.Address
	args.Status = Node.Status

	reply := AddFingerEntryReply{}
	fmt.Print(fmt.Sprint(Node.Port) + " joining a chord ring using the port of the curent node \n")

	ok := callNode("NodeClient.AddNode", &args, &reply)

	//ok := callNode("NodeClient.SendTest", &args, &reply, Node.JoinAddress, Node.JoinPort)
	if ok  {
		fmt.Print("JoinChord is working  \n")
		//Node.ReciveTest(&args, &reply)

	} else {
		fmt.Print("JoinChord has failed \n")
	}

}

func callNode(rpcname string, args interface{}, reply interface{}) bool {
	//Check if Sprint works

	addressPort := Node.JoinAddress + ":" + fmt.Sprint(Node.JoinPort)
	fmt.Print(addressPort + "\n")
	c, err := rpc.DialHTTP("tcp", addressPort)

	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Print("the error is " + fmt.Sprint(err) )
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

//---------------------------------------------------RPC functions--------------------------------------------------------

// A function to add a new node the Chord ring by a node which exist already in the ring

func (n *NodeClient) AddNode(args *AddFingerEntryArgs, reply *AddFingerEntryReply) error {
	fmt.Print("adding a node to the chord ring \n")
	//n.Lock.Lock()
	//defer n.Lock.Unlock()
	clientPort := args.Port
	//clientAdress := args.Address
	clientStatus := args.Status
	if clientStatus == 1 {
	

		for i := 0; i < len(n.FingerTable); i++ {
			fmt.Println("in the if statement")

			if n.FingerTable[i] == "" {
				var entry = AddEntry(fmt.Sprint(clientPort), i)
				n.FingerTable[i] = fmt.Sprintf("%040x", (entry))
				reply.Status = 2
				fmt.Println("Done adding")
				break
			}
		}
	} else {
		fmt.Print("Failed to add node to the finger table")
	}

	return nil
}

// test
func (n *NodeClient) SendTest(args *FindSuccessorArgs, reply *FindSuccessorReply) error {

	reply.String = args.String
	return nil
}

// test
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

// Computes n + 2^(i-1) mod
func AddEntry(address string, fingerentry int) *big.Int {
	const keySize = sha1.Size * 8
	var two = big.NewInt(2)
	var hashMod = new(big.Int).Exp(big.NewInt(2), big.NewInt(keySize), nil)

	n := HashString(address)
	fingerentryminus1 := big.NewInt(int64(fingerentry) - 1)
	jump := new(big.Int).Exp(two, fingerentryminus1, nil)
	sum := new(big.Int).Add(n, jump)

	return new(big.Int).Mod(sum, hashMod)
}

// Hash a string using sha1
func HashString(elt string) *big.Int {
	hasher := sha1.New()
	hasher.Write([]byte(elt))
	return new(big.Int).SetBytes(hasher.Sum(nil))
}



func(n *NodeClient) Stabilize(){
	
	// First get all successors
	// Successors := n.getSuccessors(n.Address)
	// Check for error
	// if err != nil {

	// 	fmt.Print("Could not get successors")

	// 	if n.Successors[0] == "" {
	// 		// Then use the node itself as successor
	// 		n.Successors[0] = n.Address
			
	// 	}else{
	// 		// Successor at place 0 might not be working so we need to find a new one
	// 		for i := 0; i < len(n.Successors); i++ {
	// 			if i < len(n.Successors)-1 {
	// 				n.Successors[i] = n.Successors[i+1]
	// 			}else{
	// 				n.Successors[i] = ""
	// 			}
	// 		}
	// 	}

	// }else{
	// 	for i := 0; i < len(Successors); i++ {
	// 		// This depends if we have the node in polace 0 or not?
	// 		n.Successors[i] = Successors[i]
	// 	}
	// }


	// // get predecessor of successor 0

	// if err == nil {
	// 	// Get successor 
	
	// 	// get predecessor identifier
	// 	// if predecessor identifier is between n and successor 0
	// 	// then successor 0 is predecessor
	// 	if between(n.Address, Predecessor.address, successorIdentifier, false) {
	// 		//n.Successors[0] = Predecessor
	// 	}
	// }

	// // Notify the successor that it is the predecessor

	// // Now delete all backups successors


	// // If only one node in the ring then DO SOMETHING ELSE
	// if n.Successors[0] == n.Address {
	// 	// Do something
	// }

	// // Iterate through nodes bucket and copy file to successor 0 backup
	// for key, value := range n.Bucket {
	// 	if value == "" {
	// 		break
	// 	}else{
	// 		copy files
	// 	}
	// }

	// return nil



}
