// Filename: netoworking2.go

package main

import (
	"fmt"
	"net"
	"time"
	"io"
	"bufio" //for buffered line reading
	"strings" //to trim whitespace
	"os" //for file writing
	"flag" //for commandline flags
)

func main() {
	//commandline flag for port configuration
	port := flag.String("port", "4000", "Port to listen on")
	flag.Parse()
	//define the target  host and port we want to connect to
	address := ":" + *port
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Printf("server listening on %s\n", address)
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

	//get the ip address for the file name
	clientIP := strings.Split(clientAddr, ":") [0]
	logFileName := clientIP + ".log"

	//open the file or create if one doesn't exist
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("[%s] Error creating log file for %s: %v\n", time.Now().Format(time.RFC3339), clientAddr, err)
		return
	}
	defer logFile.Close() // close file when its done

	reader := bufio.NewReader(conn) //buffered reader to read lines

	for {
		line, err := reader.ReadString('\n') //read until new line
		if err != nil {
			//handle client disconnection gracefully
			if err == io.EOF {
			fmt.Printf("[%s] Disconnected: %s\n", time.Now().Format(time.RFC3339), clientAddr) //log disconnection in UTC
		}else{
			fmt.Printf("[%s] Error reading from %s: %v\n", time.Now().Format(time.RFC3339), clientAddr, err)
		}
			return
		}

		clean := strings.TrimSpace(line) //trim whitespace

		//log message to the file
		logEntry := fmt.Sprintf("[%s] %s\n", time.Now().Format(time.RFC3339),clean)
		_, err = logFile.WriteString(logEntry)
		if err != nil {
			fmt.Printf("[%s] Error writing to log file for %s: %v\n", time.Now().Format(time.RFC3339), clientAddr, err)
		}


		_, err = conn.Write([]byte(clean + "\n")) // echo clean line
		if err != nil {
			fmt.Printf("[%s] Error writing to %s: %v\n", time.Now().Format(time.RFC3339), clientAddr, err)
			return
		}
	}
}
