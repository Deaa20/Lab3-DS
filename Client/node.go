package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"sync"
	"time"
)

//status 0 means that the node is leaving the  chord ring
//status 1 means that the node is joining the Chord ring
//status 2 means that the node is in the chord ring

const fingerTableSize int = 161

type NodeClient struct {
	Address                  string
	Port                     int
	FingerTable              [fingerTableSize]NodeInfo
	Predecessor              NodeInfo
	Successors               []NodeInfo
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
	next                     int
}

var Node NodeClient

type NodeInfo struct {
	Address string
	Port    int
	NodeID  *big.Int
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
		Status:                   1,
		next:                     0,
	}

	if clientID == "" {
		Node.ClientID = HashString(fmt.Sprintf("%s:%d", address, port))

	} else {
		c := new(big.Int)
		c.SetString(clientID, 10)
		Node.ClientID = c
	}

	fmt.Println("successor" + fmt.Sprint(Node.NumSuccessors))

	if Node.JoinAddress == "" {
		NewChord()
	} else {
		JoinChord()
	}
	return Node

}

func NewChord() {
	Node.Predecessor = NodeInfo{}
	Node.Successors = make([]NodeInfo, Node.NumSuccessors)
	for i := range Node.Successors {
		Node.Successors[i].Address = Node.Address
		Node.Successors[i].Port = Node.Port
		Node.Successors[i].NodeID = Node.ClientID
	}

	for i := range Node.FingerTable {
		Node.FingerTable[i].Address = Node.Address
		Node.FingerTable[i].Port = Node.Port
		Node.FingerTable[i].NodeID = Node.ClientID
	}
	go Node.Server(Node.Port)
}

func JoinChord() {

	Node.Predecessor = NodeInfo{}
	Node.Successors = make([]NodeInfo, Node.NumSuccessors)
	for i := range Node.Successors {
		Node.Successors[i].Address = Node.Address
		Node.Successors[i].Port = Node.Port
		Node.Successors[i].NodeID = Node.ClientID
	}

	for i := range Node.FingerTable {
		Node.FingerTable[i].Address = Node.Address
		Node.FingerTable[i].Port = Node.Port
		Node.FingerTable[i].NodeID = Node.ClientID
	}
	go Node.Server(Node.Port)
	args := FindSuccessorArgs{}
	args.ID = Node.ClientID
	reply := FindSuccessorReply{}

	fmt.Print(fmt.Sprint(Node.Port) + " joining a chord ring using the port of the curent node \n")
	done := false
	address := Node.JoinAddress
	port := fmt.Sprint(Node.JoinPort)

	fmt.Println("address" + address)
	fmt.Println("Port" + port)
	maxtries := 5

	for done != true && maxtries > 0 {
		ok := callNode("NodeClient.FindSuccessor", &args, &reply, address, port)

		if ok {

			address = reply.Successor.Address
			port = fmt.Sprint(reply.Successor.Port)
			done = reply.Final
			fmt.Println("THIS IS THE SUCCESSOR I WAS GIVEN ; GONNA TRY IT NOW ", reply.Successor.NodeID)

		} else {
			fmt.Println("error finding successor ")
		}

		maxtries = maxtries - 1

	}
	Node.Successors[0].Address = address
	Node.Successors[0].NodeID = reply.Successor.NodeID
	num, err := strconv.Atoi(port)

	if err != nil {

		fmt.Println("Error:", err)
		return
	}
	Node.Successors[0].Port = num

}

func callNode(rpcname string, args interface{}, reply interface{}, address string, port string) bool {
	//Check if Sprint works

	addressPort := address + ":" + port
	fmt.Print(addressPort + "NodeChord" + "\n")
	c, err := rpc.DialHTTP("tcp", addressPort)

	fmt.Println("here im")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Print("the error is " + fmt.Sprint(err))
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

func stabilizeloop() {
	for Node.Done() == false {
		Stabilize()
		time.Sleep(time.Duration(Node.StabilizeInterval) * time.Millisecond)
	}
}

func fixfingersloop() {
	for Node.Done() == false {
		FixFingers()
		time.Sleep(time.Duration(Node.FixFingersInterval) * time.Millisecond)
	}
}

func checkpredecessorloop() {
	for Node.Done() == false {
		CheckPredecessor()
		time.Sleep(time.Duration(Node.CheckPredecessorInterval) * time.Millisecond)
	}
}

//---------------------------------------------------RPC functions--------------------------------------------------------

func between(start, elt, end *big.Int, inclusive bool) bool {
	if end.Cmp(start) > 0 {
		return (start.Cmp(elt) < 0 && elt.Cmp(end) < 0) || (inclusive && elt.Cmp(end) == 0)
	} else {
		return start.Cmp(elt) < 0 || elt.Cmp(end) < 0 || (inclusive && elt.Cmp(end) == 0)
	}
}

func (n *NodeClient) FindSuccessor(args *FindSuccessorArgs, reply *FindSuccessorReply) error {
	fmt.Println("Successor: ", len(n.Successors))
	End := Node.Successors[0].NodeID
	if between(n.ClientID, args.ID, End, true) {
		reply.Successor = Node.Successors[0]
		reply.Final = true
	} else {
		reply.Successor = Node.Successors[0]
		reply.Final = false
	}
	return nil
}

func HashToString(value *big.Int) string {
	hex := fmt.Sprintf("%040x", value)
	return hex
}

func (n *NodeClient) Ping(args *PingArgs, reply *PingArgs) error {

	return nil
}

func CheckPredecessor() {
	args := PingArgs{}
	reply := PingReply{}

	address := Node.Predecessor.Address
	port := Node.Predecessor.Port
	ok := callNode("NodeClient.Ping", &args, &reply, address, fmt.Sprint(port))
	if !ok {

		Node.Predecessor = NodeInfo{}
	}

}

func Stabilize() {
	args := GetPredecessorArgs{}
	reply := GetPredecessorReply{}

	address := Node.Successors[0].Address
	port := Node.Successors[0].Port

	if Node.Successors[0].NodeID == Node.ClientID {

		ok := Node.GetPredecessor(&args, &reply)
		if ok == nil {

			fmt.Println(Node.ClientID)
			fmt.Println(reply.Predecessor.NodeID)
			fmt.Println(Node.Successors[0].NodeID)
			fmt.Println("was it here")
			if reply.Predecessor.NodeID != nil {
				Node.Successors[0] = reply.Predecessor

			}
			fmt.Println("Successors", reply.Successors)

			address = Node.Successors[0].Address
			port = Node.Successors[0].Port
			NotifyArgs := NotifyNodesArgs{NodeInfo{Address: Node.Address, Port: Node.Port, NodeID: Node.ClientID}}
			NotifyReply := NotifyNodesReply{}
			fmt.Print("im at the stablization and it wokr")
			if Node.Successors[0].NodeID != Node.ClientID {
				bok := callNode("NodeClient.NotifyNodes", &NotifyArgs, &NotifyReply, address, fmt.Sprint(port))
				if bok {

					if NotifyReply.Successors != nil {
						Node.Successors = NotifyReply.Successors
					}

				} else {
					fmt.Println("error")
				}
			}

		}

	} else {
		fmt.Println("\n", Node.Successors[0].NodeID, "\n")
		ok := callNode("NodeClient.GetPredecessor", &args, &reply, address, fmt.Sprint(port))
		if ok {
			fmt.Println("TEST TEST TEST TEST PREDECESSOR WAS ")
			fmt.Println(Node.ClientID)
			fmt.Println(reply.Predecessor.NodeID)
			fmt.Println(Node.Successors[0].NodeID)
			if reply.Predecessor.NodeID == nil {

			} else if between(Node.ClientID, reply.Predecessor.NodeID, Node.Successors[0].NodeID, true) {

				fmt.Println("was it here 3")
				Node.Successors[0] = reply.Predecessor

				fmt.Println("THIS SHOULD NOT HAPPEN WHEN SUCC IS 2 ")
			}
			address = Node.Successors[0].Address
			port = Node.Successors[0].Port
			NotifyArgs := NotifyNodesArgs{NodeInfo{Address: Node.Address, Port: Node.Port, NodeID: Node.ClientID}}
			NotifyReply := NotifyNodesReply{}
			fmt.Print("address", address)
			ok = callNode("NodeClient.NotifyNodes", &NotifyArgs, &NotifyReply, address, fmt.Sprint(port))
			fmt.Println("reply", NotifyReply.Successors)
			if NotifyReply.Successors != nil {
				Node.Successors = NotifyReply.Successors
			}

		}

	}
}

func (n *NodeClient) NotifyNodes(args *NotifyNodesArgs, reply *NotifyNodesReply) error {
	if Node.Predecessor.Address == "" || between(Node.Predecessor.NodeID, args.NodeAdress.NodeID, Node.ClientID, true) {

		fmt.Println(Node.Predecessor.Address)

		Node.Predecessor = args.NodeAdress

	}
	reply.Successors = CopySuccessors()

	fmt.Println("THese are the successors i am giving to my next freind", reply.Successors)

	return nil
}
func GetPredecessorLocal() {

}

func (n *NodeClient) GetPredecessor(args *GetPredecessorArgs, reply *GetPredecessorReply) error {
	fmt.Println("Hey i am node ," + Node.ClientID.String() + "i am giving you " + Node.Predecessor.NodeID.String())
	reply.Predecessor = Node.Predecessor
	reply.Successors = CopySuccessors()

	return nil
}

func CopySuccessors() []NodeInfo {

	// COPY NODE SUCESSORS,  SHIFT ALL MEMEBERS OF NODE SUCESSORS TO THE RIGHT AND ADD CURRENT NODE TO FIRST POSITION
	Succesors := make([]NodeInfo, Node.NumSuccessors)
	for i := range Succesors {
		Succesors[i] = Node.Successors[i]
	}
	for i := len(Node.Successors) - 1; i > 0; i-- {
		Succesors[i] = Node.Successors[i-1]
	}
	Succesors[0] = NodeInfo{Address: Node.Address, Port: Node.Port, NodeID: Node.ClientID}

	return Succesors

}

func FixFingers() {
	Node.next = Node.next + 1
	if Node.next >= fingerTableSize {
		Node.next = 1
	}
	entry := AddEntry(Node.Address, Node.next)
	args := FindSuccessorArgs{ID: entry}
	reply := FindSuccessorReply{}
	Node.FindSuccessor(&args, &reply)

	address := reply.Successor.Address
	port := reply.Successor.Port
	done := reply.Final
	maxtries := 5
	for done != true && maxtries > 0 {
		ok := callNode("NodeClient.FindSuccessor", &args, &reply, address, fmt.Sprint(port))

		if ok {
			address = reply.Successor.Address
			port = reply.Successor.Port
			done = reply.Final
		} else {
			fmt.Println("error finding successor ")
		}
		maxtries = maxtries - 1

	}
	Node.FingerTable[Node.next].Address = address
	Node.FingerTable[Node.next].NodeID = reply.Successor.NodeID
	Node.FingerTable[Node.next].Port = port
}

func LeaveChord() {
	Node.Status = 0
}

func Lookup(filename string, node *NodeClient) NodeInfo {
	key := HashString(filename)
	fmt.Println("key" + key.String())

	args := FindSuccessorArgs{ID: key}
	reply := FindSuccessorReply{}
	maxTries := 5
	done := false
	address := Node.Successors[0].Address
	port := fmt.Sprint(Node.Successors[0].Port)
	for done != true && maxTries > 0 {
		if Node.Successors[0].NodeID == Node.ClientID {
			ok := Node.FindSuccessor(&args, &reply)
			if ok == nil {
				address = reply.Successor.Address
				port = fmt.Sprint(reply.Successor.Port)
				done = reply.Final
				maxTries = maxTries - 1
			}
		} else {
			ok := callNode("NodeClient.FindSuccessor", &args, &reply, address, port)

			if ok {

				address = reply.Successor.Address
				port = fmt.Sprint(reply.Successor.Port)
				done = reply.Final
				maxTries = maxTries - 1

			} else {
				fmt.Println("error finding successor ")
			}

		}

	}
	return reply.Successor

}

func StoreFile() {
}

func PrintState(node *NodeClient) {
	fmt.Printf("Printing the state of the ccurrent node\n ")
	fmt.Printf("Address:" + node.Address + "\n ")
	fmt.Printf("port:" + fmt.Sprint(node.Port) + "\n ")
	fmt.Printf("predecessor:" + fmt.Sprintf("%040x", node.Predecessor.NodeID) + "\n ")
	fmt.Printf("JoinAddress:" + node.JoinAddress + "\n ")
	fmt.Printf("JoinPort:" + fmt.Sprint(node.JoinPort) + "\n ")
	fmt.Printf("StabilizeInterval:" + fmt.Sprint(node.StabilizeInterval) + "\n ")
	fmt.Printf("FixFingersInterval:" + fmt.Sprint(node.FixFingersInterval) + "\n ")
	fmt.Printf("CheckPredecessorInterval:" + fmt.Sprint(node.CheckPredecessorInterval) + "\n ")
	fmt.Printf("NumSuccessors:" + fmt.Sprint(node.NumSuccessors) + "\n ")
	fmt.Printf("ClientID:" + fmt.Sprintf("%040x", (node.ClientID)) + "\n ")

	for i := 0; i < len(node.Successors); i++ {
		fmt.Printf("successor nr :[" + fmt.Sprint(i) + "]:" + node.Successors[i].NodeID.String() + "\n ")
	}
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

//---------------------------------------------------------------------
//                    Pseudo code from discord
//---------------------------------------------------------------------

func Closest_Preceding_Node(ID *big.Int) NodeInfo {
	for i := fingerTableSize - 1; i >= 1; i-- {
		if between(Node.ClientID, Node.FingerTable[i].NodeID, ID, true) {
			return Node.FingerTable[i]
		}
	}
	return NodeInfo{Address: Node.Address, Port: Node.Port, NodeID: Node.ClientID}
}
