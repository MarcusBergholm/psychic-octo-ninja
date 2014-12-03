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
	conn, err := net.Dial("tcp", ":3333")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected") // DEBUG
	// Creates a new input reader to get input from keyboard.
	reader := bufio.NewReader(os.Stdin)
	// Fetches the menu
	fetch(conn)
	//Sends and fetches to/from server.
	for {	
		send(conn, reader)
		fetch(conn)
	}
}

// Sends standard input to server.
func send(conn net.Conn, reader *bufio.Reader) {
	//fmt.Print("Enter text: ") // DEBUG
	// Reads keyboard input until newline.
    input, _ := reader.ReadString('\n')
    input = input[:len(input)-1]
	// Creates a new buffer based 	on input string.
	buf := bytes.NewBufferString(input)
	if buf.Len() == 10 || buf.Len() == 0 {
		buf.Write(make([]byte, 1))
	}
	// Take the next 10 bytes and then sent it.
	outputbuf := buf.Next(10)
	// If outputbuf has less then 10 exit the loop
	for len(outputbuf) == 10 || len(outputbuf) == 0  {
		conn.Write(outputbuf)
		outputbuf = buf.Next(10)
	}
	if len(outputbuf) == 10 {
		conn.Write(outputbuf)
		conn.Write(make([]byte, 1))
	} else {
		conn.Write(outputbuf)
	}
}

// Fetches data from server. 10 bytes at the time.
// Breaks the fetch if if encounter a newline.
func fetch(conn net.Conn) {
	var input string
	inputLen := 10
	for inputLen == 10 {
		// Create new array to store 10 bytes of input
		inputBuffer := make([]byte, 10)
		inputLen, _ = conn.Read(inputBuffer)
		// append the 10 bytes to the input string (return string) 
		input = input + string(inputBuffer[:inputLen])
		// If we encounter a newline (ASCII code for NL == 10)
		// break the loop.  
		//if inputBuffer[inputLen-1] == 10 {break}
	}
	fmt.Println(input)
}