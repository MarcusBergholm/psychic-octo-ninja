package main

import (
    "fmt"
    "net"
    "os"
    "bytes"
    "./users"
    "strconv"
    "./lang"
)

const (
    CONN_HOST = "localhost"
    CONN_PORT = "4444"
    CONN_TYPE = "tcp"
)

var accounts []users.User

func main() {
	// Initializes the valid bank accounts
	accounts = users.InitPeople()
	lang.InitLang()

    // Listen for incoming connections.
    l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
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
	lang := &lang.EngLang
	user := login(conn, lang)
	for {
		  // Read the incoming connection 
		  s, err := fetch(conn)
		  if err != nil {
		  	logOff(user)
		  	return
		  }
		  // Send a response back to connection
		  switch s {
		  case "1": withdraw(conn, &user, lang)
		  case "2": deposit(conn, &user,lang)
		  case "3": balance(conn, &user, lang)
		  case "4": lang = changeLang(conn, lang)
		  case "5": logOff(user)
		  			user = login(conn, lang)
		  // ADMIN OPTIONS ONLY! USED TO CHANGE BANNER! (CASE 17 && 18)
		  case "17": if user.GetName() == "Admin" {
		  				changeBannerEng(conn, lang)
		  			 } else {
		  			 	send(conn, lang.Error + lang.WelcomeTxt + lang.Banner + "\n" + lang.Menu)
		  			 }
		  case "18": if user.GetName() == "Admin" {
		  				changeBannerSwe(conn, lang)
		  			 } else {
		  			 	send(conn, lang.Error + lang.WelcomeTxt + lang.Banner + "\n" + lang.Menu)
		  			 }
		  default: send(conn, lang.Error + lang.WelcomeTxt + lang.Banner + "\n" + lang.Menu)
		  }
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
	// If buffer is empty send an empty byte to
	// client. This is used to tell when transmit is done. 
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
func login(conn net.Conn, lang *lang.Lang) (user users.User){
	for {
		// Send login screen to client
		send(conn, lang.Login)
		// Fetches cardnr
		cardnr, _ := fetch(conn)
		cardnrInt, _ := strconv.Atoi(cardnr)
		// Tries to login with the fetched cardnr
		user, succes := loginRequest(conn, cardnrInt, lang)
		// If login succeeds the server will block any other connection
		// to login with the cardnr until the connection logs out
		if succes {
			fmt.Println(user.GetName() + " has connected")
			send(conn, lang.WelcomeTxt + lang.Banner + "\n" + lang.Menu)
			return user
		}
	}
}

// Tries to login to an accound based on the cardnr
func loginRequest(conn net.Conn, cardnr int, lang *lang.Lang) (u users.User, succes bool) {
	// Check if the cardnr exists in the bank
	for i, user := range accounts {
		//If it exists server tries to log in
		if user.GetCardnumber() == cardnr {
			for {
				// Promt the user to enter its password
				send(conn, lang.PassPromt)
				// Fetches the password
				pass, _ := fetch(conn)
				passInt, _ := strconv.Atoi(pass)
				// Check if the password is correct
				if user.GetPassword() == passInt {
					// Check if the account is allready active
					// if so retun false
					if user.GetStatus() {
						send(conn, lang.AllreadyActiveUser)
						fetch(conn)
						return u, false
					}
					// Login success, sets the account to active
					// and return the user
					user.SetStatus(true)
					accounts[i] = user
					return user, true
				} else {
					// Tells the client he entered wrong password
					// and he needs to try again
					send(conn, lang.WrongPass)
					fetch(conn)
				}
			}
		}
	}
	// If no cardnr matches we tell the client 
	// that the cardnr was incorrect
	send(conn, lang.WrongCardnr)
	fetch(conn)
	return u, false
}


// Log off the active user.
// Now other clients can log in on the account
func logOff(user users.User) {
	// Loops thro the accounts to find the active one.
	for index, oldUser := range accounts {
		if user.GetName() == oldUser.GetName() {
			//Sets the account to inactive
			fmt.Println(user.GetName() + " has disconnected")
			user.SetStatus(false)
			accounts[index] = user
		}
	}
}

// The withdraw function in the bank
// Promts user to enter the amount he wishes to withdraw
// and promts a two digit code for verifcation on the withdraw
func withdraw(conn net.Conn, user *users.User, lang *lang.Lang) {
	// Sends the client the withdraw interface
	send(conn, lang.Withdraw)
	// Fetches the amount the user wishes to withdraw
	withdrawAmount, _ := fetch(conn)
	withdrawAmountInt, _ := strconv.Atoi(withdrawAmount)
	// Promts user to enter his two digit code
	send(conn, lang.TwoDigitCode)
	// Fetches the two digit code
	code, _ := fetch(conn)
	codeInt, _ := strconv.Atoi(code)
	// Check if the two digit code is valid
	if user.GetTwoDigitCode(codeInt) {
		// Preform the withdraw
		user.Withdraw(withdrawAmountInt)
		// Tells the user that the withdraw was successful
		send(conn, lang.SuccessfulWithdraw + strconv.Itoa(user.GetBalance()) + "\n " + lang.WelcomeTxt + lang.Banner + "\n" + lang.Menu)
	} else {
		// Tells the user that the withdraw was unsuccessful
		send(conn, lang.Error + lang.WelcomeTxt + lang.Banner + "\n" +lang.Menu)
	}
}

// The deposit function in the bank.
// Promts the user to enter the amount he wishes to deposit
func deposit(conn net.Conn, user *users.User, lang *lang.Lang) {
	// Sends the client the deposit interface
	send(conn, lang.Deposit)
	// Fetches the deposit amount
	depositAmount, _ := fetch(conn)
	depositAmountInt, err := strconv.Atoi(depositAmount)
	if err != nil {
		// tells the user the deposit was unsuccessful
		send(conn, lang.Error + lang.WelcomeTxt + lang.Banner + "\n" + lang.Menu)
	} else {
		// Preform the withdraw
		user.Deposit(depositAmountInt)
		// Tell the user the deposit was successful
		send(conn, lang.SuccessfulDeposit + strconv.Itoa(user.GetBalance()) + "\n" + lang.WelcomeTxt + lang.Banner + "\n" + lang.Menu)
	}
}

// The balanace function in the bank.
// Tells the client the balance of the active account
func balance(conn net.Conn, user *users.User, lang *lang.Lang) {
	// Sends the balance to client
	send(conn, lang.Balance + strconv.Itoa(user.GetBalance()) +"\n" + lang.WelcomeTxt + lang.Banner + "\n" + lang.Menu)
}

// The change language in the bank.
// Changes between swedish and english based on the 
// current language selected
func changeLang(conn net.Conn, currentLang *lang.Lang) (l *lang.Lang){
	// If swedish is acitve change to english
	if currentLang.LangType == "SWE" {
		l = &lang.EngLang
		// Sends the new menu to client
		send(conn, l.WelcomeTxt + l.Banner + "\n" + l.Menu)
	// If english is active change to swedish
	} else if currentLang.LangType == "ENG" {
		l = &lang.SweLang
		// Sends the new menu to client
		send(conn, l.WelcomeTxt + l.Banner + "\n" + l.Menu)
	}
	// return the new language
	return l
}


// ADMIN FUNCTION ONLY.
// Changes the banner in the swedish language
// Max 80 characters
func changeBannerSwe(conn net.Conn, lang *lang.Lang) {
	// Check if the current language is swedish
	// If not promt the admin to change language
	// before changing the swedish banner
	if lang.LangType == "SWE" {
		// Promts the admin to change the banner (Max 80 characters) 
		send(conn, "Ändra banner (Max 80 tecken):")
		// Fetches the new banner
		banner, _ := fetch(conn)
		if len(banner) > 80 {
			// Tell the admin the banner is to long
			send(conn, "För lång banner!\n" + lang.Banner + "\n" + lang.Menu)
		} else {
			// Set the new banner for swedish
			lang.Banner  = banner
			// Tell the admin the banner was successfully changed
			send(conn, "Bytet lyckades!\n" + lang.Banner + "\n" + lang.Menu)
		}
	} else {
		// Tell the admin to change language to swedish before changing it
		send(conn, "Change to swedish before changing the banner!\n" + lang.WelcomeTxt + lang.Banner + "\n" + lang.Menu)
	}
}

// ADMIN FUNCTIIN ONLY.
// Changes the banner in the english language
// Max 80 characters
func changeBannerEng(conn net.Conn, lang *lang.Lang) {
	// Check if the current language is english
	// If not promt the admin to change language
	// before changing the english banner
	if lang.LangType == "ENG" {
		// Promts the admin to change the banner (Max 80 characters)
		send(conn, "Change banner (Max 80 characters):")
		// Fetches the new banner
		banner, _ := fetch(conn)
		if len(banner) > 80 {
			// Tell the admin the banner is to long
			send(conn, "Banner is to long!\n" + lang.Banner + "\n" + lang.Menu)
		} else {
			// Set the new banner for english
			lang.Banner = banner
			// Tell the admin the banner was successfully changed
			send(conn, "Banner successfully changed!\n" + lang.Banner + "\n" + lang.Menu)
		}
	} else {
		// Tell the admin to change language to english before changing it
		send(conn, "Ändra till engelska innan du byter banner!\n" + lang.WelcomeTxt + lang.Banner + "\n" + lang.Menu)
	}
}