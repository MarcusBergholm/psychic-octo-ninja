package users

// Represent an account in the bank.
type User struct {
		name string
		cardnumber int
    	password int
    	balance int
    	status bool
    	codes []int
    	usedCodes []int
	}

// Creates three new accounts in the bank and return
// an array consisting of the accounts
func InitPeople() (users []User){
	// Two digit code for withdraws in the bank.
	codes := []int{11,13,15,17,19,21,23,25,27,29,31,33,35,37,39,41,43,45,
					47,49,51,53,55,57,59,61,63,65,67,69,71,73,75,77,79,81,
					83,85,87,89,91,93,95,97,99}

	// Init user 1
	User1 := User{"Spongebob Squarepants", 123456789, 12345, 10000, false, codes, []int{}}

	// Init user 2
	User2 := User{"Mr.Drep", 987654321, 54321, 999, false, codes, []int{}}

	// Init user 2
	User3 := User{"Mrs.Derpina", 123123123, 123123, 123123, false, codes, []int{}}

	Admin := User{"Admin", 1337, 1337, 0, false, []int{}, []int{}}

	return []User{User1, User2, User3, Admin}
}

/* Setters for this user */	

// Set a username for this user
func (username *User) SetName(name string) {
    username.name = name
}

// Set a cardnumber for this user
func (username *User) SetCardnumber(cardnumber int) {
    username.cardnumber = cardnumber
}

// Set a new password for this user
func (username *User) SetPassword(password int) {
    username.password = password
}

// Adjust balance for this user
func (username *User) SetBalance(balance int) {
    username.balance = balance
}

// Set login status for this user
func (username *User) SetStatus(status bool) {
    username.status = status
}

// Withdraw money from this user
func (username *User) Withdraw(amount int) {
    username.balance = username.balance - amount
}

// Deposit money to this user
func (username *User) Deposit(amount int) {
    username.balance = username.balance + amount
}


/* Getters for this user */

// Get this users password
func (username *User) GetName() string {
    return username.name
}

// Get this users password
func (username *User) GetCardnumber() int {
    return username.cardnumber
}

// Get this users password
func (username *User) GetPassword() int {
    return username.password
}

// Get this users balance
func (username *User) GetBalance() int {
    return username.balance
}

// Get this users password
func (username *User) GetStatus() bool {
    return username.status
}
// Check if the code matches any of the two digit codes in user.
// If code is allready used return false
func (username *User) GetTwoDigitCode(code int) bool {
	for _, twoDigitCode := range username.codes {
		if code == twoDigitCode {
			for _, usedCode := range username.usedCodes {
				if code == usedCode {
					return false
				}
			}
			username.usedCodes = append(username.usedCodes, code)
			return true
		}
	}
	return false
}