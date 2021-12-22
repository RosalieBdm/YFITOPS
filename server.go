package main

//run the server in a cmd with the command "go run server.go portnumber" (8000 works as a portnumber)
//in another cmd run the user with the commmand "go run user.go portnumber" (and obviously, the same portnumber as the server)

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var TM = [10][10]int{} // Truce matrice to handle the subscription process
// line : User with new sub to approve, collumn : waiting subscriber

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
		// Then he gets the number of his best match
		// Then we ask him if who he wants to follow

		go checkNewSubs(c, connectionNumber)
	}
}

func checkNewSubs(c net.Conn, num int) {
	fmt.Println("Entered checkNewSubs function")
	//fmt.Println(TM)
	connectionNumber := num
	waitingS, wsNum := WaitingSub(connectionNumber, TM) // WaitingSub is a function that checks if the user has any waiting subscribers (it returns a boolean and the waiting subscriber's number)
	if waitingS {
		fmt.Println("Entered waitingS loop")
		c.Write([]byte("You have a new subscribtion from user " + strconv.Itoa(wsNum) + ", ok to share your music (Y or N) ? \n"))
		reader := bufio.NewReader(c)
		answer, err := reader.ReadString('\n')
		answer = strings.TrimSuffix(answer, "\n")
		if err != nil {
			fmt.Println(err)
			return
		}
		switch answer {
		case "Y":
			// send data
			c.Write([]byte("ok, sending musics \n"))
			go sendData(wsNum, connectionNumber, c) // sendData is a function that sends the music from one client to another with their user number
		case "N":
			c.Write([]byte("ok, no sub \n"))
			TM[connectionNumber][wsNum] = 0
			go checkNewSubs(c, connectionNumber)
		default:
			c.Write([]byte("Please answer with Y or N ? \n"))
		}
	} else {
		c.Write([]byte("No new subscription \n"))
		go FindNewFriends(connectionNumber, c)
		//go handleSubscription(c, connectionNumber) // handleSubscription is a function that proposes the user to subscribe to other users
	}

}

func handleSubscription(connection net.Conn, userNum int) {
	fmt.Println("Entered handleSubscription function")
	// If a client wants to follow someone, that someone's number is saved so that he can be asked if he approves
	for {
		connection.Write([]byte("Who do you want to follow ? (enter user number) \n"))
		subBuf := bufio.NewReader(connection)
		subRequest, err := subBuf.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		subRequest = strings.TrimSuffix(subRequest, "\n")
		subNum, err := strconv.Atoi(subRequest) //convert to int
		if err != nil {
			fmt.Println(err)
			return
		}
		waitingSubscriber := userNum
		TM[subNum][waitingSubscriber] = 1 // We implement the truce matrice
	}
}

func sendData(Waitinguser int, Connecteduser int, c net.Conn) { // TODO : use a tab instead of strings et not send the musics the user already have
	// We open our data base, a txt file in which 1 user = 1 line = 1 list of musics
	fmt.Println("Entered sendData function")
	file, err := os.Open("DataBase.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	sc := bufio.NewScanner(file)
	message, err := ReadLine(sc, Connecteduser+1) // ReadLine is a function that returns a specific line from a file
	//fmt.Println(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	WriteLine(message, Waitinguser) // WriteLine is a function that append some text at the end of a specific line
	if err != nil {
		fmt.Println(err)
		return
	}

	TM[Connecteduser][Waitinguser] = 0 // We update the truce matrice

	go checkNewSubs(c, Connecteduser)

}

func FindNewFriends(cnum int, c net.Conn) /*(pnum int, bestCorr int)*/ {
	fmt.Println("Entered FindNewFriends function")
	file, err := os.Open("DataBase.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	sc := bufio.NewScanner(file)
	//Reader := bufio.NewReader(file)

	// get the music from the connected user and put it in a tab
	musics, err := ReadLine(sc, cnum)
	if err != nil {
		fmt.Println(err)
		return
	}
	musicList := strings.Split(musics, ";")

	fmt.Println("Got the music list")
	fmt.Println(musicList)

	// get the number of users
	/*lineCount, err := getNumberOfLine(sc)
	if err != nil {
		fmt.Println(err)
		return
	}

	LineCount, err := lineCountfunction(Reader)
	if err != nil {
		fmt.Println(err)
		return
	} */

	CompareComp := [10]int{} // a tab that will stock the compability rate with every user

	//fmt.Println("line count with scan : " + strconv.Itoa(lineCount))
	//fmt.Println("line count with read : " + strconv.Itoa(LineCount))

	LineCount := 10

	for i := 1; i < LineCount+1; i++ { //on balaye les lignes (les users)
		fmt.Println("goroutine loop " + strconv.Itoa(i))
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			log.Fatal(err)
		}
		compareTwoLines(sc, musicList, i, CompareComp) // a goroutine that compares the client's music list to every user's simultaneously
	}

	fmt.Println("Finished the goroutine loop")

	CompareComp[cnum] = 0 // rajouter un if dans le for, pas de compareTwoLine si c'est cnum ou si deja 1 dans TM
	bestCorr := MaxArray(CompareComp)
	pnum := indexOf(bestCorr, CompareComp)

	fmt.Println("Got the correlation tab")
	fmt.Println(bestCorr)
	fmt.Println(pnum)

	c.Write([]byte("Your best match is user " + strconv.Itoa(pnum) + "\n"))

	go handleSubscription(c, cnum)
	//return pnum, bestCorr
}

func compareTwoLines(sc *bufio.Scanner, userMusics []string, line int, compareC [10]int) {
	fmt.Println("Entered compareTwoLines function")
	mToCompare, err := ReadLine(sc, line) // get the musics from a line
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(mToCompare)
	mToCompareList := strings.Split(mToCompare, ";") // put them in a tab
	fmt.Println(mToCompareList)
	comp := 0
	for j := 0; j < len(userMusics); j++ {
		for k := 0; k < len(mToCompareList); k++ {
			if userMusics[j] == mToCompareList[k] {
				comp += 1
			}
		}
	}
	comp = comp / len(userMusics)
	compareC[line] = int(comp)
	fmt.Println(compareC)
}

func WaitingSub(n int, list [10][10]int) (sub bool, subNum int) {
	for i := 0; i < 6; i++ {
		if list[n][i] == 1 {
			return true, i
		}
	}
	return false, -1 // No new sub
}

func ReadLine(sc *bufio.Scanner, lineNum int) (line string, err error) {
	fmt.Println("Entered ReadLine function")
	//sc := bufio.NewScanner(f)
	sc.Text()
	currentLine := 0
	for sc.Scan() {
		currentLine++
		if currentLine == lineNum {
			line = sc.Text()
			err = sc.Err()
			return
			//return sc.Text(), sc.Err()
		}
	}
	return line, nil //io.EOF
}

/*func getNumberOfLine(sc *bufio.Scanner) (lineCount int, err error) {
	lineCount = 0
	for sc.Scan() {
		if true {
			lineCount++
			return
		}

	}
	//readerDone := make(chan struct{}, 1)
	//defer close(readerDone)
	return lineCount, sc.Err()
}*/

func lineCountfunction(r io.Reader) (n int, err error) {
	var count int
	const lineBreak = '\n'

	buf := make([]byte, bufio.MaxScanTokenSize)

	for {
		bufferSize, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}

		var buffPosition int
		for {
			i := bytes.IndexByte(buf[buffPosition:], lineBreak)
			if i == -1 || bufferSize == buffPosition {
				break
			}
			buffPosition += i + 1
			count++
		}
		if err == io.EOF {
			break
		}
	}

	return count, nil
}

func WriteLine(data string, lineNum int) {
	fmt.Println("Entered WriteLine function")
	input, err := ioutil.ReadFile("util")
	if err != nil {
		fmt.Println(err)
		return
	}
	lines := strings.Split(string(input), "\n")
	lines[lineNum] = strings.TrimSuffix(lines[lineNum], "\n") + " " + strings.TrimSuffix(data, "\n")
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile("util", []byte(output), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func MaxArray(array [10]int) int {
	var max int = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
	}
	return max
}

func indexOf(element int, data [10]int) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}
