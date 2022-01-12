package main

//run the server in a cmd with the command "go run server.go portnumber" (8000 works as a portnumber)
//in another cmd run the user with the commmand "go run user.go portnumber" (and obviously, the same portnumber as the server)

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"fslock"
)

var TM = [10000][10000]int{} // Truce matrice to save the subscriptions (subscription = sub request + approval)
// line : user that's being followed, collumn : following user

var TMSP = [10000][10000]int{} // Truce matrice to handle the subscription process (subscription process = sub request waiting for an approval)
// line : user with new sub to approve, collumn : waiting subscriber

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

		// Getting the users number (this identification will be used all along our program)
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
		// Then we ask him who he wants to follow

		go checkNewSubs(c, connectionNumber)
	}
}

func checkNewSubs(c net.Conn, num int) { // to understand better how the function works, please see the joined file "checkNewSubs_explanation"
	fmt.Println("Entered checkNewSubs function")
	connectionNumber := num
	waitingS, wsNum := WaitingSub(connectionNumber) // WaitingSub is a function that checks if the user has any waiting subscribers (it returns a boolean and the waiting subscriber's number)
	if waitingS {
		fmt.Println("Entered waitingS loop")
		c.Write([]byte("You have a new subscribtion from user " + strconv.Itoa(wsNum) + ", ok to share your music (Y or N) ? \n")) // We ask for the approval of the client to share his musics with other users
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
			// cancel the subscription request
			c.Write([]byte("ok, no sub \n"))
			TMSP[connectionNumber][wsNum] = 0
			go checkNewSubs(c, connectionNumber)
		default:
			c.Write([]byte("Please answer with Y or N ? \n"))
		}
	} else {
		// when there is no subscription request, we move on
		c.Write([]byte("No new subscription \n"))
		go FindNewFriends(connectionNumber, c)
	}

}

func handleSubscription(connection net.Conn, userNum int) { // // If a client wants to follow someone, that someone's number is saved so that he can be asked if he approves when he connects to the server
	fmt.Println("Entered handleSubscription function")
	for {
		connection.Write([]byte("Who do you want to follow ? (enter user number) \n"))
		subBuf := bufio.NewReader(connection)
		subRequest, err := subBuf.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		subRequest = strings.TrimSuffix(subRequest, "\n")
		subNum, err := strconv.Atoi(subRequest) // the number of the user our client wants to follow
		if err != nil {
			fmt.Println(err)
			return
		}
		// a client cannot follow himself
		if subNum == userNum {
			connection.Write([]byte("You can't follow yourself \n"))
			go handleSubscription(connection, userNum)
		} else {
			waitingSubscriber := userNum
			TMSP[subNum][waitingSubscriber] = 1 // We implement the truce matrice for subscription process (TMSP)
		}
	}
}

func sendData(Waitinguser int, Connecteduser int, c net.Conn) { // TODO : use a tab instead of strings et not send the musics the user already have
	// We open our data base, a txt file in which 1 user = 1 line = 1 list of musics
	fmt.Println("Entered sendData function")
	file, err := os.Open("Util.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	sc := bufio.NewScanner(file)
	message, err := ReadLine(sc, Connecteduser) // ReadLine is a function that returns a specific line from a file (see details below)
	fmt.Println(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	WriteLine(message, Waitinguser-1) // WriteLine is a function that appends some text at the end of a specific line (see details below)

	//We use a lock to protect our data
	lock := fslock.New("DataBase.txt")
	lockErr := lock.TryLock()
	if lockErr != nil {
		fmt.Println("failed to acquire lock > " + lockErr.Error())
		return
	}

	fmt.Println("Aquired the lock")
	WriteLine(message, Waitinguser) // WriteLine is a function that appends some text at the end of a specific line (see details below)

	// release the lock
	lock.Unlock()
	fmt.Println("Released the lock")
	if err != nil {
		fmt.Println(err)
		return
	}
	

	TMSP[Connecteduser][Waitinguser] = 0   // We update the truce matrice for subscription process
	TM[Connecteduser-1][Waitinguser-1] = 1 // also the truce matrice for subscription

	go checkNewSubs(c, Connecteduser)

}

func FindNewFriends(cnum int, c net.Conn) { // A function that returns the best match of the client, meaning the user with the more music in common
	fmt.Println("Entered FindNewFriends function")
	c.Write([]byte("Searching for your best musical match...\n"))
	file, err := os.Open("DataBase.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	sc := bufio.NewScanner(file)

	// get the music from the client's line in our database
	musics, err := ReadLine(sc, cnum)
	if err != nil {
		fmt.Println(err)
		return
	}
	musicList := strings.Split(musics, ";") // We put it in a tab in which one case = one music

	fmt.Println("User's music list : ")
	fmt.Println(musicList)

	defer file.Close()

	var CompareCompArray [10000]float64                                        // a tab that will stock the compability rate with every user
	var CompareCompSlice []float64 = CompareCompArray[0:len(CompareCompArray)] // We use a slice because arrays cannot be modified in a function --> IMPROVAL IDEA : Pointer to array

	LineCount := 10000

	for i := 1; i < LineCount+1; i++ { // For every user in the data base, we want to compare their music list with our client's
		// could (should) be improved : we open and close our file in every loop to reset the scanner, a reader may be more appropriate ?
		file, err := os.Open("DataBase.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		sc := bufio.NewScanner(file)
		if i != (cnum) && TM[i-1][cnum-1] == 0 { // we don't compare the user with himself or with the users he already follows
			go compareTwoLines(sc, musicList, i, CompareCompSlice) // a goroutine that compares the client's music list to every user's simultaneously (see details below)
		}
		defer file.Close()
	}

	bestCorr := MaxSlice(CompareCompSlice)      // We get the best correlation score
	pnum := indexOf(bestCorr, CompareCompSlice) // and the user with this best correlation score
	c.Write([]byte("Your best match is user " + strconv.Itoa(pnum+1) + " with a compability rate of " + strings.TrimSuffix(strconv.FormatFloat(bestCorr*100, 'E', -1, 64), "E+00") + "% \n"))

	go handleSubscription(c, cnum)
}

func compareTwoLines(sc *bufio.Scanner, userMusics []string, line int, compareC []float64) { // A funnction that calculates the correlation score between two slice
	mToCompare, err := ReadLine(sc, line) // We get the musics from a line
	if err != nil {
		fmt.Println(err)
		return
	}
	mToCompareList := strings.Split(mToCompare, ";") // We put them in a slice
	comp := 0.0
	// we compare every bow of a slice with every ones of another
	for j := 0; j < len(userMusics); j++ {
		for k := 0; k < len(mToCompareList); k++ {
			if userMusics[j] == mToCompareList[k] {
				comp++
			}
		}
	}
	comp = comp / float64(len(userMusics)) // We divide by the number of musics our client has : that's our correlation score
	compareC[line-1] = comp                // We store this value in the slice, the index of the value is the number of the user
}

func WaitingSub(n int) (sub bool, subNum int) { // A function that returns the user number of a waiting subscriber if there is any
	for i := 0; i < 6; i++ {
		if TMSP[n][i] == 1 {
			return true, i
		}
	}
	return false, -1 // No new sub
}

func ReadLine(sc *bufio.Scanner, lineNum int) (line string, err error) { // A function that uses a scanner to read a .txt file and return the line corresponding to the argument lineNum
	currentLine := 0
	for sc.Scan() {
		currentLine++
		if currentLine == lineNum {
			line = sc.Text()
			err = sc.Err()
			return
		}
	}
	return line, nil
}

func lineCountfunction(r io.Reader) (n int, err error) { // A fonction that counts the number of line of a file using a reader
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

func WriteLine(data string, lineNum int) { // A function that writes at the end of a specific line : we get the text, and split it line by line, in order to put every line in the boxes a tab. Then we add some text in the box of a specific line, et rebuild the text, to write it in the file
	fmt.Println("Entered WriteLine function")
	input, err := ioutil.ReadFile("Util.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	lines := strings.Split(string(input), "\n")
	lines[lineNum] = strings.Trim(lines[lineNum], "\n") + strings.Trim(" ", "\n") + strings.Trim(data, "\n")
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile("Util.txt", []byte(output), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func MaxSlice(slice []float64) float64 { // A function that returns the maximum value of a slice
	var max float64 = slice[0]
	for _, value := range slice {
		if max < value {
			max = value
		}
	}
	return max
}

func indexOf(element float64, data []float64) int { // A function that returns the index of a specific value in a slice
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}
