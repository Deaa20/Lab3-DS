package main

//make predecesoor and successor store address and port

import (
	"flag"
	"fmt"
	"os"
)

type Key string
type NodeAddress string

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

		
		fmt.Printf("success")


	//takeCommand()

	/* 	if joinAddress == "" {
	   		create(&Client)
	   	} else {
	   		join(&Client)
	   	}
	*/
}

func CreateNode(address string, port int, joinAddress string, joinPort int, stabilizeInterval int,
	fixFingersInterval int, checkPredecessorInterval int, numSuccessors int, clientID string) NodeClient {

	var node NodeClient

	Node = NodeClient{
		Address:                  address,
		Port:                     port,
		JoinAddress:              joinAddress,
		StabilizeInterval:        stabilizeInterval,
		FixFingersInterval:       fixFingersInterval,
		CheckPredecessorInterval: checkPredecessorInterval,
		NumSuccessors:            numSuccessors,
		ClientID:                 clientID,
	}

	if Node.JoinAddress == "" {
		node = NewChord()
	} else {
		JoinChord()

	}
	return node

}

func takeCommand() {
	var command string
	fmt.Println("------------------------------------------\nWelcome to the Chord server.\n------------------------------------------\nType Lookup + file name to search for a file\nStoreFile + file name to store a file on the Chord ring\nPrintState to print out the local state of the Chord client")
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



func join(node *NodeClient) {
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
