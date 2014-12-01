package main

import (
    "fmt"
    "net"
    "bufio"
    "log"
    "os"
)


func main() {
	conn, err := net.Dial("tcp", ":3333")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected")
	reader := bufio.NewReader(os.Stdin)
	for {	
		fmt.Print("Enter text: ")
	    text, _ := reader.ReadString('\n')

		var a int
		buff := make([]byte, 1024)
		conn.Write([]byte(text))
		a, err = conn.Read(buff)
		if err != nil {
			log.Fatal(err)
		}
		s := string(buff[:a])
		fmt.Println(s)
	}
}