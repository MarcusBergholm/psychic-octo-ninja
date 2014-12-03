package main

import (
    "fmt"
    "net"
    "os"
    "bytes"
    "./users"
    "strconv"
)

const (
    CONN_HOST = "localhost"
    CONN_PORT = "4444"
    CONN_TYPE = "tcp"
    TEMP_MENU = "Welcome to Bank\n(1) Withdraw\n(2) Deposit\n(3) Balance\n(4) Change language\n(5) Log off"
    TEMP_MENU_SWE = "Välkommen till Bank\n(1) Uttag\n(2) Insättning\n(3) Saldo"
    TEMP_LOGIN = "Welcome to Bank\nTo login enter your card nr: "
)

var accounts []users.User

func main() {
	// Initializes the valid bank accounts
	accounts = users.InitPeople()

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
	user := login(conn)
	for {
		  // Read the incoming connection 
		  s, err := fetch(conn)
		  if err != nil {
		  	logOff(user)
		  	return
		  }
		  // Send a response back to connection
		  switch s {
		  case "1": withdraw(conn, user)
		  case "2": deposit(conn, user)
		  case "3": balance(conn, user)
		  case "4": changeLang(conn)
		  case "5": logOff(user)
		  			user = login(conn)
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
	// Take the next 10 bytes and then sent it.
	outputbuf := buf.Next(10)
	// If outputbuf has less then 10 exit the loop
	for len(outputbuf) == 10 {
		conn.Write(outputbuf)
		outputbuf = buf.Next(10)
	}
	if len(outputbuf) == 0 {
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
		// If user disconnected, return with error
		if err != nil {
	   		return "", err
	   	}
		// append the 10 bytes to the input string (return string) 
		input = input + string(inputBuffer[:inputLen])
	}
	return input, nil
}

// Sends login screen to client and promts it to enter cardnr
// and password. If user enterd a valid cardnr and password 
// the user logs on and will see the Banks menu
func login(conn net.Conn) (user users.User){
	for {
		send(conn, "Please enter your cardnumber:")
		cardnr, _ := fetch(conn)
		cardnrInt, _ := strconv.Atoi(cardnr)
		user, succes := loginRequest(conn, cardnrInt)
		if succes {
			fmt.Println(user.GetName() + " has connected")
			send(conn, TEMP_MENU)
			return user
		}
	}
}

func loginRequest(conn net.Conn, cardnr int) (u users.User, succes bool) {
	for i, user := range accounts {
		if user.GetCardnumber() == cardnr {
			for {
				send(conn, "Please enter your password:")
				pass, _ := fetch(conn)
				passInt, _ := strconv.Atoi(pass)
				if user.GetPassword() == passInt {
					if user.GetStatus() {
						send(conn, "User allready active, press any key to try login again")
						fetch(conn)
						return u, false
					}
					user.SetStatus(true)
					accounts[i] = user
					return user, true
				} else {
					send(conn, "Wrong password, press any key to try again")
					fetch(conn)
				}
			}
		}
	}
	send(conn, "Wrong cardnumber, press any key to try again")
	fetch(conn)
	return u, false
}


// Log off the active user.
func logOff(user users.User) {
	for index, oldUser := range accounts {
		if user.GetName() == oldUser.GetName() {
			fmt.Println(user.GetName() + " has disconnected")
			user.SetStatus(false)
			accounts[index] = user
		}
	}
}

func withdraw(conn net.Conn, user users.User) {
	send(conn, "Enter the amount you wish to withdraw")
	withdrawAmount, _ := fetch(conn)
	withdrawAmountInt, _ := strconv.Atoi(withdrawAmount)
	send(conn, "Please enter your two diget code to confirm your withdraw")
	code, _ := fetch(conn)
	codeInt, _ := strconv.Atoi(code)
	if user.GetTwoDigitCode(codeInt) {
		user.Withdraw(withdrawAmountInt)
		send(conn, "You succesfully took out: " + withdrawAmount + "\n " + TEMP_MENU)
	} else {
		send(conn, "Wrong code\n" + TEMP_MENU)
	}
}
func deposit(conn net.Conn, user users.User) {
	send(conn, "Enter the amount you wish to deposit")
	depositAmount, _ := fetch(conn)
	depositAmountInt, _ := strconv.Atoi(depositAmount)
	user.Deposit(depositAmountInt)
	send(conn, "You succesfully entered: " + depositAmount + "\n" + TEMP_MENU)
}
func balance(conn net.Conn, user users.User) {
	send(conn, "Your balance is: "+ strconv.Itoa(user.GetBalance()) +"\n" + TEMP_MENU)
}
func changeLang(conn net.Conn) {}











