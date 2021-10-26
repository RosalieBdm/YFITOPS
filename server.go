package main
//run the server in a cmd with the command "go run server.go portnumber" (8000 works as a portnumber)
//in another cmd run the user with the commmand "go run user.go portnumber" (and obviously, the same portnumber as the server)

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"list"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	CONNECT := ":" + arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(c, text+"\n")

		message, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Print("->: " + message)
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
	}
	//MUSIC := 

}
