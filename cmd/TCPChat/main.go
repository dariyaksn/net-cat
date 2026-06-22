package main

import (
	"fmt"
	"netcat/netcat"
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(1)
	}
	port := ":8989"

	if len(os.Args) == 2 {
		port = ":" + strings.TrimSpace(os.Args[1])
	}

	srv := netcat.NewServer(10, "General")
	fmt.Println("Starting server at", port)
	if err := srv.Start(port); err != nil {
		fmt.Println("Server error:", err)
	}
}
