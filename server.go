package main

import (
    "fmt"
    "net"
    "os"
    "bytes"
)

const (
    CONN_HOST = "localhost"
    CONN_PORT = "3333"
    CONN_TYPE = "tcp"
    TEMP_MENU = "Welcome to Bank\n(1) Withdraw\n(2) Deposit\n(3) Balance\n(4) Change language\n(5) Log off"
    TEMP_MENU_SWE = "Välkommen till Bank\n(1) Uttag\n(2) Insättning\n(3) Saldo"
    TEMP_LOGIN = "Welcome to Bank\nTo login enter your card nr:"
)

func main() {
    // Listen for incoming connections.
    l, err := net.Listen(CONN_TYPE,":"+CONN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
    for {
	    // Listen for an incoming connection.
	    conn, err := l.Accept()

	    if err != nil {
	        fmt.Println("Error accepting: ", err.Error())
	        os.Exit(1)
	    }
	    // Handle connections in a new goroutine.
	    go handleConnection(conn)
	}
}

// Handles incoming requests for the connection conn.
func handleConnection(conn net.Conn) {
	fmt.Println("User connected") // DEBUG for server
	login(conn)
	for {
		  // Read the incoming connection 
		  s, err := fetch(conn)
		  if err != nil {return}
		  // Send a response back to connection
		  switch s {
		  case "1": withdraw(conn)
		  case "2": deposit(conn)
		  case "3": balance(conn)
		  case "4": changeLang(conn)
		  case "5": err = logOff(conn)
		  default: send(conn, "Incorrect input\n" + TEMP_MENU)
		  }
		  if err != nil {return}
	}
}

// Sends the string input to the client with only
// 10 bytes at the time.
func send(conn net.Conn, output string) {
	// Creates a new buffer based on input string.
	buf := bytes.NewBufferString(output)
	if buf.Len() == 10 {
		buf.Write(make([]byte, 1))
	}
	// Take the next 10 bytes and then sent it.
	outputbuf := buf.Next(10)
	// If outputbuf has less then 10 exit the loop
	for len(outputbuf) == 10 {
		conn.Write(outputbuf)
		outputbuf = buf.Next(10)
	}

	if len(outputbuf) == 10 || len(outputbuf) == 0 {
		conn.Write(outputbuf)
		conn.Write(make([]byte, 1))
	} else {
		conn.Write(outputbuf)
	}
}

// Fetch what the client is requesting.
// Breaks the featch by encountering a newline in the input.
func fetch(conn net.Conn) (input string, err error) {
	inputLen := 10
	for inputLen == 10 {
		// Create new array to store 10 bytes of input
		inputBuffer := make([]byte, 10)
		inputLen, err = conn.Read(inputBuffer)
		// If we do not have any connection with the client 
		// disconnect the client
		if err != nil {
			fmt.Println("User disconnected") // Only prints in server.
		    return "", err
		}
		// append the 10 bytes to the input string (return string) 
		input = input + string(inputBuffer[:inputLen])
		// If we encounter a newline (ASCII code for NL == 10)
		// break the loop.  
		//if inputBuffer[inputLen-1] == 10 {break}
	}
	return input, nil
}

// Sends login screen to client and promts it to enter cardnr
// and password. If user enterd a valid cardnr and password 
// the user logs on and will see the Banks menu
func login(conn net.Conn) {
	send(conn, TEMP_LOGIN)
	cardNum, _ := fetch(conn)
	fmt.Println(cardNum)
	send(conn, "Pleace enter your password:")
	password, _ := fetch(conn)
	fmt.Println(password)
	send(conn, TEMP_MENU)
}

func withdraw(conn net.Conn) {
	send(conn, "Enter the amount you wish to withdraw")
	withdrawAmount, _ := fetch(conn)
	fmt.Println(withdrawAmount)
	send(conn, "You succesfully took out: " + withdrawAmount + "\n " + TEMP_MENU)
}
func deposit(conn net.Conn) {
	send(conn, "Enter the amount you wish to deposit")
	depositAmount, _ := fetch(conn)
	fmt.Println(depositAmount)
	send(conn, "You succesfully entered: " + depositAmount + "\n" + TEMP_MENU)
}
func balance(conn net.Conn) {
	send(conn, "You have no balance\n" + TEMP_MENU)
}
func changeLang(conn net.Conn) {}
func logOff(conn net.Conn) (err error) {
	return err
}

















