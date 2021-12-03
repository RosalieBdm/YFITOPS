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
	// to handle connection, we need the host port and the user number
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port and user number.")
		return
	}
	CONNECT := ":" + arguments[1]
	//NUMBER, err := strconv.Atoi(arguments[2])
	/*if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}*/
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

func checkIfQuestion(message string) bool {
	//Question := false
	for i := 0; i < len(message); i++ {
		if string(message[i]) == "?" {
			return true
		}
	}
	return false
}

/*scanner := bufio.NewScanner(os.Stdin)
scanner.Scan()
num := scanner.Text()
//NUMBER, err := strconv.Atoi(num)
fmt.Fprintf(c, num)
time.Sleep(8 * time.Second)
if err != nil {
	panic(err)
}*/

/*if Question == "ok, no sub \n" || Question == {
	continue
}*/

// First, the user sends his own number so the server can check if he's got any waiting sub
//fmt.Printf("What's your user number ? \n")
//fmt.Fprintf(c, string(NUMBER))

/*
	// We open our data base, a txt file in which 1 user = 1 line = 1 list of liked musics
	MUSIC, err := ioutil.ReadFile("Users.txt")
	file, err := os.OpenFile("Users.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	defer file.Close() // on ferme automatiquement à la fin de notre programme

	// When connected, a user will enter the number of another user he wants to follow
	// That way when he subscribes to him, he gets his musics in his line of the data base

	//fmt.Print("Who do you want to follow ? Enter user number \n")


		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		subscription := scanner.Text()
		subscriptionNumber, err := strconv.Atoi(subscription)
		fmt.Fprintf(c, string(subscriptionNumber)+"\n")
		if err != nil {
			panic(err)
		}
		// A function to handle the subscription
		// We need to send the sub request to the server
		//subRequest := []int{NUMBER, subscriptionNumber}
		fmt.Fprintf(c, subscription) // bonus : also send the NUMBER
		// When the other user is connected, he gets the information that he has a sub waiting
		// He says yes or no
		// if yes, they exchange their musics
		//message, _, err := ReadLine(file, subscriptionNumber) //we get the music from the subscribed
		// if no, we send a message to inform the first user

		//fmt.Print("->: " + message)
		//fmt.Fprintf(c, string(message)+"\n")
		//_, err = file.WriteString(string("\n" + message + "\n")) // bientôt remplacé par une fonction

		if err != nil {
			panic(err)
		}

		if strings.TrimSpace(string(MUSIC)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}

		if MUSIC != nil {
			fmt.Println(err)
		} */

/*
func ReadLine(r io.Reader, lineNum int) (line string, lastLine int, err error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			return sc.Text(), lastLine, sc.Err()
		}
	}
	return line, lastLine, io.EOF
}
*/
// We need a function to write in the end of a line
