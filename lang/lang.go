package lang


// Used in server.go to represent 
// different languages to the client
type Lang struct {
	LangType string
	WelcomeTxt string
	Login string
	Menu string
	Error string
	PassPromt string
	AllreadyActiveUser string
	WrongPass string
	WrongCardnr string
	Withdraw string
	TwoDigitCode string
	SuccessfulWithdraw string
	Deposit string
	SuccessfulDeposit string
	Balance string
	Banner string
}

// Used languages in server.
// Server uses a pointer to these Langs.
var EngLang Lang
var SweLang Lang

func InitLang() {
	// a struct with english as language
	EngLang = Lang{"ENG",
				"Welcome to Bank\n",
				"To login enter your card nr:",
				"(1) Withdraw\n(2) Deposit\n(3) Balance\n(4) Change language\n(5) Log off",
				"Incorrect input\n",
				"Please enter your password:",
				"User allready active, press any key to try login again",
				"Wrong password, press any key to try again",
				"Wrong cardnumber, press any key to try again",
				"Enter the amount you wish to withdraw",
				"Please enter your two digit code to confirm your withdraw",
				"Withdraw succesful. Your new balance is: ",
				"Enter the amount you wish to deposit",
				"Deposit succesful. Your new balance is: ",
				"Your balance is: ",
				"The best Bank that exists!"}

	// a struct with swedish as language
	SweLang = Lang{"SWE",
				"Välkommen till Bank\n",
				"För att logga in slå in ditt kortnummer:",
				"(1) Uttag\n(2) Insättning\n(3) Saldo\n(4) Byt språk\n(5) Logga ut",
				"Felaktig indata",
				"Slå in din kod:",
				"Användare redan aktiv, tryck på valfri tanget för att försöka igen",
				"Felaktigt lösenord, tryck på valfri tanget för att försöka igen",
				"Felaktigt kortnummer, tryck på valfri tanget för att försöka igen",
				"Ange summan du önskar att ta ut:",
				"Ange din tvåsiffriga för att genomföra uttaget:",
				"Uttag lyckades. Ditt nya saldo är: ",
				"Ange den summa du önskar sätta in:",
				"Insättning lyckades. Ditt nya saldo är: ",
				"Ditt saldo är: ",
				"Bästa Banken som existerar!"}
}