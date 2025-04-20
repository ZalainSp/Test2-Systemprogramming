// Filename: netoworking2.go

package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	//define the target  host and port we want to connect to
	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("server listeing on :4000")
	// our program runs an infinite loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		//handle each client connection using a go routine
		go handleConnection(conn)
	}

}
func handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String() //get client IP
	fmt.Printf("[%s] Connected: %s\n", time.Now().Format(time.RFC3339), clientAddr) //log address and timestamp in UTC

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("[%s] Disconnected: %s\n", time.Now().Format(time.RFC3339), clientAddr) //log disconnection in UTC
			return
		}
		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Println("error writing to client:", err)
		}
	}
}
