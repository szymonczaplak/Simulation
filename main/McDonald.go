package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	HAMBURGER  = 1
	CHEESEBURGER  = 2
	CHICKENBURGER  = 3
	FRYTKI = 4
	COLA = 5
	NUGGETSY = 6
	)

type client struct {
	index          int
	timeOfEntrance time.Time
	timeWaited float64
	order []int
}

type device struct {
	available, inUse int
}

type workers struct {
	kitchenWorkersAvalible, kitchenWorkersInUse, caschiersAvalible, caschiersInUse int
}

type readyMeals struct {
	readyHamburger, readyCheeseburger, reedyChickenburger, readyFries, readyCola, readyNuggets int
}

func serve_clients(clients [] *client){
	var allClientsGroup sync.WaitGroup
	allClientsGroup.Add(1)

	finishedOrder := make(chan *client, len(clients))
	for _, cli := range clients{
		go make_order(cli, finishedOrder)
		go wait_for_order(finishedOrder)
	}
	go check_all_done(&allClientsGroup)
	allClientsGroup.Wait()
}



func create_client(clientId int) *client{
	var order []int
	number_of_elements := rand.Intn(6)
	for i :=0 ; i<number_of_elements; i++{
		order = append(order, rand.Intn(7)+1)
	}
	return &client{clientId, time.Now(), 0, order}
}

func make_client_queue(number_of_clients int) []*client{
	var out []*client
	for i:=0; i<number_of_clients; i++{
		out = append(out, create_client(i))
		fmt.Printf("Client %d entered bakery\n", i)
	}
	return out
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
var frying_pans = device{3,0}
var checkouts = device{3,0}
var selfCheckouts = device{3,0}
var friesMaker = device{3,0}
var tray = device{10,0}
var ready = readyMeals{5,5,5,5,5,5}


func main(){
	clients := make_client_queue(10)
	get_clients_average_time(clients)
	fmt.Print("Balance: ", profit - lose, "\n")
}