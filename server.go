package main

//run the server in a cmd with the command "go run server.go portnumber" (8000 works as a portnumber)
//in another cmd run the user with the commmand "go run user.go portnumber" (and obviously, the same portnumber as the server)

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	//defer l.Close()

	NBCO := 0
	for {
		c, err := l.Accept()
		NBCO += 1
		fmt.Println(NBCO)
		d, err := l.Accept()
		NBCO += 1
		fmt.Println(NBCO)

		if err != nil {
			fmt.Println(err)
			return
		}

		go exchange(c, d, NBCO)
	}
}

func exchange(connection1 net.Conn, connection2 net.Conn, connum int) {
	defer connection1.Close()
	c := bufio.NewReader(connection1)
	defer connection2.Close()
	d := bufio.NewReader(connection2)

	for {
		netData1, err1 := c.ReadString('\n')
		netData2, err2 := d.ReadString('\n')

		if err1 != nil {
			fmt.Println(err1)
			return
		}
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		if strings.TrimSpace(string(netData1)) == "STOP" {
			fmt.Println("Exiting TCP server!")
			return
		}
		if strings.TrimSpace(string(netData2)) == "STOP" {
			fmt.Println("Exiting TCP server!")
			return
		}

		//fmt.Print("-> ", string(netData1))
		//fmt.Print("-> ", string(netData2))
		/*t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"*/
		connection1.Write([]byte(netData2))
		connection2.Write([]byte(netData1)) //[]byte(myTime)
		//fmt.Println("j'en suis la")
		break
	}
}
