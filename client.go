package main

import (
    "fmt"
    "net"
    "bufio"
    "log"
    "os"
    "bytes"
)

// Main function for client. Connects to the server thro
// dial. If connected to server, client will loop send and 
// fetch to/from server.
func main() {
	// Dials the server.
	conn, err := net.Dial("tcp", ":4444")
	// If client fails to connect to server 
	// use fatal error
	if err != nil {
		log.Fatal(err)
	}
	// Creates a new input reader to get input from keyboard.
	reader := bufio.NewReader(os.Stdin)
	// Fetches Login from server
	fetch(conn)
	//Sends and fetches to/from server.
	for {	
		send(conn, reader)
		fetch(conn)
	}
}

// Sends standard input to server.
func send(conn net.Conn, reader *bufio.Reader) {
	// Reads keyboard input until newline.
    input, _ := reader.ReadString('\n')
    // Remove \n from input
    input = input[:len(input)-1]
	// Creates a new buffer based on input string.
	buf := bytes.NewBufferString(input)
	// Take the next 10 bytes and then sent it to server.
	outputbuf := buf.Next(10)
	for len(outputbuf) == 10 {
		conn.Write(outputbuf)
		outputbuf = buf.Next(10)
	}
	// If the last package was exactly 10 bytes 
	// we send a unsignd byte to let the server know
	// we are done sending. If the last package is less
	// then 10 bytes then we know we are sending unsignd 
	// at the end of the outputbuf
	if len(outputbuf) == 0 {
		conn.Write(make([]byte, 1))
	} else {
		conn.Write(outputbuf)
	}
}

// Fetches data from server. 10 bytes at the time.
// Breaks the fetch if if encounter a newline.
func fetch(conn net.Conn) {
	// Used to store the input from server
	var input string
	inputLen := 10 
	// If server sends a packages that is less then 10 bytes
	// then we know the server is done sending
	for inputLen == 10 {
		// Create new array to store 10 bytes of input
		inputBuffer := make([]byte, 10)
		inputLen, _ = conn.Read(inputBuffer)
		// append the 10 bytes to the input string (return string) 
		input = input + string(inputBuffer[:inputLen])	
	}
	// Prints the input from the server
	fmt.Println(input)
}