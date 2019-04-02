package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type client struct {
	index          int
	timeOfEntrance time.Time
	timeWaited float64
}

type mugs struct{
	availableMugs, mugsInUse int
}

func (mugs_ *mugs) breakMug() bool{
	if rand.Intn(100) > 80{
		profit-= 20
		return true
	}
	return false
}

func (mugs_ *mugs) getMug(){
	fmt.Print("Barista: Getting mug...\n")
	if mugs_.availableMugs == 0 {
		fmt.Printf("Barista: Oh I'mugs_ sorry, but our retarded bakery has only %d mugs...\n", mugs_.availableMugs+mugs_.mugsInUse)
	}
	for {
		if mugs_.availableMugs > 0 {
			if mugs_.breakMug(){
				fmt.Print("Mug broken :C\n")
				mugs_.availableMugs --
			}
			time.Sleep(1 * time.Second)
			mugs_.availableMugs --
			mugs_.mugsInUse ++
			fmt.Print("Barista: Oh... Here it is!!\n")
			break
		}
		time.Sleep(60*time.Millisecond)
	}
}

func createClient(index int) *client{
	return &client{index, time.Now(), 0}
}

func makeClientQueue(numberOfClients int) []*client{
	/*
	Create queue of clients of size 'numberOfClients'
	 */
	var out []*client
	for i:=0; i< numberOfClients; i++{
		out = append(out, createClient(i))
		fmt.Printf("Client %d entered bakery\n", i)
	}
	return out
}

func getCoffee(client_ chan *client, mugs_ *mugs, gotCoffee chan int){
	/*
	This function handles batista which is preparing coffee and client which is waiting.
	 */
	mu.Lock()
	profit += 8.50
	cli := <-client_
	fmt.Printf("***Barista is serving client %d***\n", cli.index)
	mugs_.getMug()
	time.Sleep(3 * time.Second)
	fmt.Printf("Client %d got coffe!\n", cli.index)
	//fmt.Printf(time.Now().String() + "\n")
	gotCoffee <- cli.index
	cli.timeWaited = time.Since(cli.timeOfEntrance).Seconds()
	fmt.Printf("Client %d waited %f seconds\n", cli.index, cli.timeWaited)
	fmt.Print("Cleaning express\n")
	time.Sleep(500 * time.Millisecond)
	mu.Unlock()
}

func clientIsEnjoyingCoffe(index int, finished chan int){
	/*
	This function handles client which is drinking his coffee.
	 */
	fmt.Printf("Client %d is enjoying his coffe\n", index)
	//fmt.Printf(time.Now().String() + "\n")

	//time.Sleep(time.Duration(rand.Intn(4)) * time.Second)
	time.Sleep(10*time.Second)
	fmt.Printf("Client %d finished his awesome coffe!\n", index)
	finished <- index
}

func enjoyCoffe(clients [] *client, got chan int, finished chan int){
	/*
	Receive through channel 'got' index of client which has got his coffee, and started drinking it.
	Then this function starts goroutine to handle their behaviour.

	This is intended to be run as goroutine.
	 */
	enjoyed := 0
	for enjoyed != len(clients){
		ind := <-got
		go clientIsEnjoyingCoffe(ind, finished)
		enjoyed++
	}
}

func clientIsReturningMug(index int, mugs_ *mugs, threadGroup *sync.WaitGroup){
	/*
	Receive index of client which is returning mug and handle it.

	This is intended to be run as goroutine.
	 */
	if mugs_.breakMug(){
		fmt.Printf("Client %d broken his mug\n", index)
		defer threadGroup.Done()
		return
	}
	fmt.Printf("Client %d is returning his mug!\n", index)
	time.Sleep(2*time.Second)
	mugs_.availableMugs ++
	mugs_.mugsInUse --
	fmt.Printf("Client %d returned his mug\n", index)
	defer threadGroup.Done()
}

func returnMug(clients [] *client, mugs_ *mugs, finished chan int, threadGroup *sync.WaitGroup){
	/*
	Receive through channel 'finished' index of client which finished his coffee and started to return mug and start
	goroutine to handle their behaviour.

	This is intended to be run as goroutine.
	 */
	returned := 0
	var retgr sync.WaitGroup
	for returned != len(clients){
		ind := <-finished
		retgr.Add(1)
		go clientIsReturningMug(ind, mugs_, &retgr)
		returned ++
	}
	retgr.Wait()
	defer threadGroup.Done()
}

func serveClients(clients [] *client, m *mugs){
	/*
	Main function to launch goroutines, which simulate real-life behaviour, for every client in clients
	 */
	var allClientsGroup sync.WaitGroup
	allClientsGroup.Add(1)
	c := make(chan *client, len(clients))
	gotCoffe := make(chan int, len(clients))
	finishedEnjoyingCoffe := make(chan int)
	for _, i := range clients{
		c <- i
		go getCoffee(c, m, gotCoffe)
	}
	go enjoyCoffe(clients, gotCoffe, finishedEnjoyingCoffe)
	go returnMug(clients, m, finishedEnjoyingCoffe, &allClientsGroup)
	allClientsGroup.Wait()
}

func getClientsAverageTime(clients []*client){
	/*
	Iterate over every client_ in clients and count average time of waiting for product
	 */
	var average float64
	for _, client_ := range clients{
		average += float64(client_.timeWaited)
	}
	average /= float64(len(clients))
	fmt.Printf("Average time of wainting: %f \n", average)
}

var mu = &sync.Mutex{}
var profit = 0.0
var lose = 0.0

func main(){
	clients := makeClientQueue(10)
	mugs := &mugs{15,0}
	serveClients(clients, mugs)
	getClientsAverageTime(clients)
	fmt.Print("Balance: ", profit - lose, "\n")
}