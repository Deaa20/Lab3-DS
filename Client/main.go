package main

//make predecesoor and successor store address and port

import (
	"crypto/rand"
	"crypto/rsa"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
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
	if port < 0 || port > 65535 {
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

	publicKey, privateKey, err := GenerateKeyPair()
	if err != nil {
		fmt.Println("Error: Failed to generate key pair.")
	}

	CreateNode(address, port, joinAddress, joinPort, stabilizeInterval,
		fixFingersInterval, checkPredecessorInterval, numSuccessors, clientID, publicKey, privateKey)

	fmt.Print("waiting in the main \n")
	//firstTime := true
	for Node.Done() == false {
		//if firstTime == true {
		Stabilize()
		//FixFingers()
		//firstTime = false
		//}
		takeCommand(&Node)
		time.Sleep(time.Second)

	}

}

// Generate a new key pair for the node
func GenerateKeyPair() (*rsa.PublicKey, *rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	publicKey := &privateKey.PublicKey

	return publicKey, privateKey, nil
}

func takeCommand(node *NodeClient) {
	var command string
	var searchedFileName string
	var storedFileName string
	fmt.Println("------------------------------------------Welcome to the Chord server.------------------------------------------\n Type the following commands with respect to the case sensetivty \n" +
		" Lookup  to search for a file \n Store to store a file on the Chord ring\n PrintState to print out the local state of the Chord client \n" +
		"Leave for leaving the chord ring")

	fmt.Scanln(&command)

	if command == "Leave" {
		node.Status = 0
	} else if command == "Lookup" {
		fmt.Println("Provide the name of the file you are searching for")
		fmt.Scanln(&searchedFileName)
		searchedFileName = strings.TrimSpace(searchedFileName)
		Successor := Lookup(searchedFileName, node)
		fmt.Println("The file is stored on the  node with id : ", Successor.NodeID, "address and port :", Successor.Address, ":", Successor.Port)

	} else if command == "StoreFile" {
		fmt.Println("Provide the name of the file you are storing")
		fmt.Scanln(&storedFileName)
		storedFileName = strings.TrimSpace(storedFileName)
		clientStoreFile(storedFileName, node)
	} else if command == "PrintState" {
		PrintState(node)

	} else {
		fmt.Print("The provided command is invalid, please follow the guide providein the welcome message. ")
	}
}
