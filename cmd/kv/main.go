package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const SocketPath = "/tmp/kvstore.sock"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: No command provided, or too few arguments")
		os.Exit(1)
	}	 

	conn, err := net.Dial("unix", SocketPath)
	if err != nil {
		fmt.Println("ERROR: cannot connect to daemon:", err)
		os.Exit(1)
	}
	defer conn.Close()

	cmd := strings.Join(os.Args[1:], " ") + "\n"
	_, err = conn.Write([]byte(cmd))
	if err != nil {
		fmt.Println("ERROR: write failed:", err)
		os.Exit(1)
	}

	resp, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("ERROR: read failed:", err)
		os.Exit(1)
	}
	fmt.Println(strings.TrimSpace(resp))
}
