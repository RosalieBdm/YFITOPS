package main

//run the server in a cmd with the command "go run server.go portnumber" (8000 works as a portnumber)
//in another cmd run the user with the commmand "go run user.go portnumber" (and obviously, the same portnumber as the server)

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

var WS []int      // Users with new subs
var wSUBNum []int // Users waiting for their subscription to be approved

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	// We listen to the port
	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	// To keep an eye on the number of connection
	NBCO := 0

	for {
		// If a client calls this port, we accept
		c, err := l.Accept()
		NBCO += 1
		fmt.Print(NBCO)
		fmt.Println(" users connected")

		if err != nil {
			fmt.Println(err)
			return
		}

		// Getting the users number
		c.Write([]byte("What's your user number ? \n"))
		connectionNumberString, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("user number : " + connectionNumberString)
		connectionNumberString = strings.TrimSuffix(connectionNumberString, "\n")
		connectionNumber, err := strconv.Atoi(connectionNumberString)
		if err != nil {
			fmt.Println(err)
			return
		}

		// First we check if the connected user has new subscribers, and ask if he's willing to share his music
		// Then we ask him if who he wants to follow

		go checkNewSubs(c, connectionNumber)
	}
}

func checkNewSubs(c net.Conn, num int) {
	fmt.Println("Entered checkNewSubs function")
	// The WS tab contains the number of all the users who've got new subscribers
	connectionNumber := num
	if NumInWS(connectionNumber, WS) {
		fmt.Println("entered NumInWS loop")
		numIndex := getNumIndex(connectionNumber, WS)
		c.Write([]byte("You have a new subscriber, ok to share your music (Y or N) ? \n"))
		reader := bufio.NewReader(c)
		answer, err := reader.ReadString('\n')
		answer = strings.TrimSuffix(answer, "\n")
		if err != nil {
			fmt.Println(err)
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		switch answer {
		case "Y":
			// send data
			c.Write([]byte("ok, sending musics \n"))
			go sendData(wSUBNum[numIndex], connectionNumber, c)
		case "N":
			c.Write([]byte("ok, no sub \n"))
			go handleSubscription(c, connectionNumber)
			//return
		default:
			c.Write([]byte("Please answer with Y or N ? \n"))
		}
	} else {
		c.Write([]byte("No new subscription \n"))
		go handleSubscription(c, connectionNumber)
	}

}

func handleSubscription(connection net.Conn, userNum int) {
	fmt.Println("entered handleSubscription function")
	//defer connection.Close()
	// If a client wants to follow someone, that someone's number is saved so that he can be asked if he approves
	// To do : implement a truce matrice
	for {
		connection.Write([]byte("Who do you want to follow ? (enter user number) \n"))
		fmt.Println("message sent")
		subBuf := bufio.NewReader(connection)
		subRequest, err := subBuf.ReadString('\n')
		subRequest = strings.TrimSuffix(subRequest, "\n")
		if err != nil {
			fmt.Println(err)
			return
		}
		subRequest = strings.TrimSuffix(subRequest, "\n")
		subNum, err := strconv.Atoi(subRequest)
		if err != nil {
			fmt.Println(err)
			return
		}
		WS = append(WS, subNum)
		waitingSubscriber := userNum
		wSUBNum = append(wSUBNum, waitingSubscriber)

		fmt.Println(WS)
		fmt.Println(wSUBNum)
	}
}

func sendData(Waitinguser int, Connecteduser int, c net.Conn) {
	// We open our data base, a txt file in which 1 user = 1 line = 1 list of liked musics
	fmt.Println("entered sendData function")
	file, err := os.Open("utilisateurs.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	message, _, err := ReadLine(file, Connecteduser-1)
	fmt.Println(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	WriteLine(message, Waitinguser)

	if err != nil {
		fmt.Println(err)
		return
	}

	go handleSubscription(c, Connecteduser)

}

func NumInWS(n int, list []int) bool {
	// Checks if a number is in a list
	for _, b := range list {
		if b == n {
			return true
		}
	}
	return false
}

func getNumIndex(num int, list []int) int {
	for k, v := range list {
		if num == v {
			return k
		}
	}
	return 0 //not found.
}

func ReadLine(r io.Reader, lineNum int) (line string, lastLine int, err error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			return sc.Text(), lastLine, sc.Err()
		}
	}
	return line, lastLine, nil //io.EOF
}

func WriteLine(data string, lineNum int) {
	input, err := ioutil.ReadFile("utilisateurs.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	lines := strings.Split(string(input), "\n")
	lines[lineNum] = strings.TrimSuffix(lines[lineNum], "\n") + " " + data
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile("utilisateurs.txt", []byte(output), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}
