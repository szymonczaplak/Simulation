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
	avalible_mugs, mugs_in_use int
}

func (m *mugs) break_mug() bool{
	if rand.Intn(100) > 80{
		profit-= 20
		return true
	}
	return false
}

func (m *mugs) get_mug(){
	fmt.Print("Barista: Getting mug...\n")
	if m.avalible_mugs == 0 {
		fmt.Printf("Barista: Oh I'm sorry, but our retarded bakery has only %d mugs...\n", m.avalible_mugs+m.mugs_in_use)
	}
	for {
		if m.avalible_mugs > 0 {
			if m.break_mug(){
				fmt.Print("Mug broken :C\n")
				m.avalible_mugs --
			}
			time.Sleep(1 * time.Second)
			m.avalible_mugs --
			m.mugs_in_use ++
			fmt.Print("Barista: Oh... Here it is!!\n")
			break
		}
		time.Sleep(60*time.Millisecond)
	}
}

func create_client(n int) *client{
	return &client{n, time.Now(), 0}
}

func make_client_queue(number_of_clients int) []*client{
	var out []*client
	for i:=0; i<number_of_clients; i++{
		out = append(out, create_client(i))
		fmt.Printf("Client %d entered bakery\n", i)
	}
	return out
}

func get_coffe(c chan *client, m *mugs, got_coffe chan int){
	mu.Lock()
	profit += 8.50
	cli := <-c
	fmt.Printf("***Barista is serving client %d***\n", cli.index)
	m.get_mug()
	time.Sleep(3 * time.Second)
	fmt.Printf("Client %d got coffe!\n", cli.index)
	//fmt.Printf(time.Now().String() + "\n")
	got_coffe <- cli.index
	cli.timeWaited = time.Since(cli.timeOfEntrance).Seconds()
	fmt.Printf("Client %d waited %f seconds\n", cli.index, cli.timeWaited)
	fmt.Print("Cleaning express\n")
	time.Sleep(500 * time.Millisecond)
	mu.Unlock()
}

func client_is_enjoying_coffe(index int, finished chan int){
	fmt.Printf("Client %d is enjoying his coffe\n", index)
	//fmt.Printf(time.Now().String() + "\n")

	//time.Sleep(time.Duration(rand.Intn(4)) * time.Second)
	time.Sleep(10*time.Second)
	fmt.Printf("Client %d finished his awesome coffe!\n", index)
	finished <- index
}

func enjoy_coffe(clients [] *client, got chan int, finished chan int){
	enjoyed := 0
	for enjoyed != len(clients){
		ind := <-got
		go client_is_enjoying_coffe(ind, finished)
		enjoyed++
	}
}

func client_is_returning_mug(ind int, m *mugs, group *sync.WaitGroup){
	if m.break_mug(){
		fmt.Printf("Client %d broken his mug\n", ind)
		defer group.Done()
		return
	}
	fmt.Printf("Client %d is returning his mug!\n", ind)
	time.Sleep(2*time.Second)
	m.avalible_mugs ++
	m.mugs_in_use --
	fmt.Printf("Client %d returned his mug\n", ind)
	defer group.Done()
}

func return_mug(clients [] *client, m *mugs, finished chan int, group *sync.WaitGroup){
	returned := 0
	var retgr sync.WaitGroup
	for returned != len(clients){
		ind := <-finished
		retgr.Add(1)
		go client_is_returning_mug(ind, m, &retgr)
		returned ++
	}
	retgr.Wait()
	defer group.Done()
}

func serve_clients(clients [] *client, m *mugs){
	var all_clients_group sync.WaitGroup
	all_clients_group.Add(1)
	c := make(chan *client, len(clients))
	got_coffe := make(chan int, len(clients))
	finished_enjoying_coffe := make(chan int)
	for _, i := range clients{
		c <- i
		go get_coffe(c, m, got_coffe)
	}
	go enjoy_coffe(clients, got_coffe, finished_enjoying_coffe)
	go return_mug(clients, m, finished_enjoying_coffe, &all_clients_group)
	all_clients_group.Wait()
}

func get_clients_average_time(clients []*client){
	var average float64
	for _,i := range clients{
		average += float64(i.timeWaited)
	}
	average /= float64(len(clients))
	fmt.Printf("Average time of wainting: %f \n", average)
}

var mu = &sync.Mutex{}
var profit = 0.0
var lose = 0.0

func main(){
	clients := make_client_queue(10)
	mugs := &mugs{15,0}
	serve_clients(clients, mugs)
	get_clients_average_time(clients)
	fmt.Print("Balance: ", profit - lose, "\n")
}