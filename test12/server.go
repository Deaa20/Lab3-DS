package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Coordinator struct {
	nReduce              int
	files                []string
	completeMapTasks     int
	completedReduceTasks int
	// Your definitions here.
}

// Your code here -- RPC handlers for the worker to call.

// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	fmt.Println("Server listening on :1234...")
	http.Serve(l, nil)
}

func main() {
	fmt.Println("run server.go first")
	c := Coordinator{}
	c.server()
	select {} // Block indefinitely to keep the server running
}
