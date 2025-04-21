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

const maxMessageLength = 1024; //max messages (length in bytes)

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

	//create inavtivity timer for 30 seconds
	inactivityTimer := time.NewTimer(30 * time.Second)

	//goroutine to handle the timeout
	go func(){
		<-inactivityTimer.C
		fmt.Printf("[%s] Inactivity timeout, disconnecting %s\n", time.Now().Format(time.RFC3339), clientAddr)
		conn.Close()
	}()

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

		//resets the timer when there is activity
		if !inactivityTimer.Stop() {
			<-inactivityTimer.C
		}
		inactivityTimer.Reset(30 * time.Second)

		clean := strings.TrimSpace(line) //trim whitespace

		//hand the server presonality mode responses
		var response string
		switch {
		case clean == "/time":
			//respond with the current server time
			response = time.Now().Format(time.RFC3339)
		case clean == "/quit":
			//close the connection
			response = "Goodbye!"
			conn.Write([]byte(response + "\n"))
			fmt.Printf("[%s] Disconnected: %s\n", time.Now().Format(time.RFC3339), clientAddr) // log disconnection
			return
		case strings.HasPrefix(clean, "/echo "):
			//echo the message after the /echo command
			response = clean[6:] //remove the /echo  part
		case clean == "hello":
			response = "Hi there!"
		case clean == "bye":
			response = "Goodbye!"
			conn.Write([]byte(response + "\n"))
			fmt.Printf("[%s] Disconnected: %s\n", time.Now().Format(time.RFC3339), clientAddr) // log disconnection
			return
		case clean == "":
			response = "Say something..."
		default:
			//check message length and truncate
			if len(clean) > maxMessageLength {
				clean = clean[:maxMessageLength] //truncate the message
				_, err := conn.Write([]byte("Message too long and was truncated.\n"))
				if err != nil {
					fmt.Printf("[%s] Error writing to %s: %v\n", time.Now().Format(time.RFC3339), clientAddr, err)
					return
				}
			}
			//the response should reflect the truncated message
			response = clean
		}

	//log message to the file
		logEntry := fmt.Sprintf("[%s] %s\n", time.Now().Format(time.RFC3339), clean)
		_, err = logFile.WriteString(logEntry)
		if err != nil {
			fmt.Printf("[%s] Error writing to log file for %s: %v\n", time.Now().Format(time.RFC3339), clientAddr, err)
		}

		///echo the response back to the client
		_, err = conn.Write([]byte(response + "\n"))
		if err != nil {
			fmt.Printf("[%s] Error writing to %s: %v\n", time.Now().Format(time.RFC3339), clientAddr, err)
			return
		}
	}
}
