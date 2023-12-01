package main


//make predecesoor and successor store address and port

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type Key string

type NodeAddress string

type Node struct {
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

var Client Node

func main() {
	// Check if there are command-line arguments
	if len(os.Args) < 1 {
		fmt.Println("Please provide an input argument.")
		return
	}
	inputValidator()

	// Get the first command-line argument (excluding the program name)

}

func inputValidator() {
	var address string
	var port int
	var joinAddress string
	var joinPort int
	var stabilizeInterval int
	var fixFingersInterval int
	var checkPredecessorInterval int
	var numSuccessors int
	var clientID string

	flag.StringVar(&address, "a", "", "The IP address that the Chord client will bind to.")
	flag.IntVar(&port, "p", 0, "The port that the Chord client will bind to and listen on.")
	flag.StringVar(&joinAddress, "ja", "", "The IP address of the machine running a Chord node.")
	flag.IntVar(&joinPort, "jp", 0, "The port that an existing Chord node is bound to and listening on.")
	flag.IntVar(&stabilizeInterval, "ts", 0, "The time in milliseconds between invocations of ‘stabilize’.")
	flag.IntVar(&fixFingersInterval, "tff", 0, "The time in milliseconds between invocations of ‘fix fingers’.")
	flag.IntVar(&checkPredecessorInterval, "tcp", 0, "The time in milliseconds between invocations of ‘check predecessor’.")
	flag.IntVar(&numSuccessors, "r", 0, "The number of successors maintained by the Chord client.")
	flag.StringVar(&clientID, "i", "", "The identifier (ID) assigned to the Chord client.")

	flag.Parse()

	if address == "" {
		fmt.Printf("Error: Please input an adress ")
		os.Exit(1)
	}
	if port < 0 && port > 65535 {
		fmt.Printf("Error:Please provide a valid port")
		os.Exit(1)

	}
	if (joinAddress != "" && joinPort == 0) || (joinPort != 0 && joinAddress == "") {
		fmt.Println("Error: Both --ja and --jp must be specified if one is provided.")
		os.Exit(1)
	}
	if stabilizeInterval < 1 || stabilizeInterval > 60000 ||
		fixFingersInterval < 1 || fixFingersInterval > 60000 ||
		checkPredecessorInterval < 1 || checkPredecessorInterval > 60000 {
		fmt.Println("Error: --ts, --tff, and --tcp must be specified with values in the range [1, 60000].")
		flag.Usage()
		os.Exit(1)
	}

	if numSuccessors < 1 || numSuccessors > 32 {
		fmt.Println("Error: -r must be specified with a value in the range [1, 32].")
		flag.Usage()
		os.Exit(1)
	}
	if clientID == "" {
		//hash the adress
	}

	CreateNode(address, port, joinAddress, joinPort, stabilizeInterval,
		fixFingersInterval, checkPredecessorInterval, numSuccessors, clientID)

	/* 	if joinAddress == "" {
	   		create(&Client)
	   	} else {
	   		join(&Client)
	   	}
	*/
}

type NodeInfo struct {
	Address string
	Port    int
}

func CreateNode(address string, port int, joinAddress string, joinPort int, stabilizeInterval int,
	fixFingersInterval int, checkPredecessorInterval int, numSuccessors int, clientID string) {

	Client = Node{Address: address,
		Port:                     port,
		JoinAddress:              joinAddress,
		StabilizeInterval:        stabilizeInterval,
		FixFingersInterval:       fixFingersInterval,
		CheckPredecessorInterval: checkPredecessorInterval,
		NumSuccessors:            numSuccessors,
		ClientID:                 clientID,
	}

	if Client.JoinAddress == "" {
		newChord()
	} else {
		joinChord()

	}

}

func newChord() {
	successors := make([]string, Client.NumSuccessors)
	for i := range successors {
		successors[i] = Client.Address
	}
	fingerTable := make([]NodeInfo, 160)
	for i := range fingerTable {
		fingerTable[i] = NodeInfo{Address: Client.Address, Port: Client.Port}
	}
	Client.Predecessor = ""
}

func joinChord() {
	Client.Predecessor = ""
	args := FindSuccessorArgs{}
	reply := FindSuccessorReply{}
	ok := callNode("main.FindSuccessor", args, reply, Client.JoinAddress, Client.JoinPort)
	if ok {

	}

}

func takeCommand() {
	var command string
	fmt.Println("Welcome to the Chord server.\nType Lookup + file name to search for a file \n \n StoreFile + file name to store a file on the Chord ring  \n PrintState to print out the local state of the Chord client")
	fmt.Scanln(&command)
	if command == "Lookup" {
	}
	if command == "StoreFile" {
	}
	if command == "PrintState" {
	} else {
		fmt.Print("The provided command is invalid, please follow the guide providein the welcome message. ")
		takeCommand()
	}
}

func (n *Node) server(port int) {
	rpc.Register(n)
	rpc.HandleHTTP()
	portString := ":" + string(port)
	l, e := net.Listen("tcp", portString)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	////fmt.Println("listening on port 1234 with adress ", l.Addr())
	go http.Serve(l, nil)
}

func join(node *Node) {
	fmt.Printf("Joining an existing Chord ring at %s:%d...\n", node.JoinAddress, node.JoinPort)

	existingNode := fmt.Sprintf("%s:%d", node.JoinAddress, node.JoinPort)
	fmt.Println("Node already exist : ", existingNode)

	// successors := node.getSuccessors(existingNode)
	//Check for error

	// Get successor list from node
	// node.updateSuccessors(successors)
	// node.updateFingerTable()

	// TODO: Update successor list

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
