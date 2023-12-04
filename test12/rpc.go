package main

import (
	"os"
	"strconv"
)

// type ExampleArgs struct {
// 	X int
// }

// type ExampleReply struct {
// 	Y int
// }

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/test12"
	s += strconv.Itoa(os.Getuid())
	return s
}