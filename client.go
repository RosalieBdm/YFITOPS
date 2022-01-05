package main

//run the server in a cmd with the command "go run server.go portnumber" (8000 works as a portnumber)
//in another cmd run the user with the commmand "go run user.go portnumber" (and obviously, the same portnumber as the server)

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// to handle connection, we need the host port
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port")
		return
	}
	CONNECT := ":" + arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		//Listens to the server, if it's a question, answers it
		MessageBuff := bufio.NewReader(c)
		Message, err := MessageBuff.ReadString('\n')
		fmt.Println(Message)

		if err != nil {
			fmt.Println(err)
			return
		}

		if checkIfQuestion(Message) {
			var txt string
			for {
				_, err := fmt.Scanln(&txt)
				if err != nil {
					fmt.Println("Veuillez rentrer une valeur valide svp.")
					continue
				}
				break
			}
			fmt.Fprintf(c, txt+"\n")
		}
	}

}

func checkIfQuestion(message string) bool { // checkIfQuestion is a function that returns true if a string has a "?" in it
	for i := 0; i < len(message); i++ {
		if string(message[i]) == "?" {
			return true
		}
	}
	return false
}
