# Test2-Systemprogramming

How to run:
1 - Run the server (default port 4000): go run main.go
1.1 - Run the server on custom port: go run main.go --port (your desired port) eg. go run main.go --port 3000
2 - connect using TCP client: nc localhost 4000 (using netcat on default port 4000)
2.1 - connect using TCP client on custom port: nc localhost (the port you entered when running) eg. nc localhost 3000

Features:
-Input whitespace trimming
-Command-based protocol (/time, /quit, /echo)
-Per-client logging to IP based .log files
-Message overflow protection
-Graceful disconnect handling
-Custom responses (hello, bye, empty input)

Implementing the log all messages feature was the most educationally enriching feature helping me understand the file I/O more in go. It showed me the 
importance of features that back up bebugging in these type of applications.
The inactivity feature was the most researched upon feature. I did not expect it to be so simple but it took a while
to understand fully, I also had to make sure it worked properly due to how go reads inputs and closed the connection safely using a goroutine.

Demonstration video
https://youtu.be/a01TEXdaKns 
