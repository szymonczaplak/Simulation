package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const (
	HAMBURGER  = 1
	CHEESEBURGER  = 2
	CHICKENBURGER  = 3
	FRIES = 4
	COLA = 5
	NUGGETS = 6
	)

type client struct {
	index          int
	timeOfEntrance time.Time
	timeWaited time.Duration
	order []int
}

type device struct {
	available, inUse int
}

type workers struct {
	kitchenWorkersAvalible, kitchenWorkersInUse, caschiersAvalible, caschiersInUse int
}

type readyMeals struct {
	readyHamburger, readyCheeseburger, readyChickenburger, readyFries, readyCola, readyNuggets int
}

func count(what int, order []int) int {
	amount := 0
	for _, i := range order{
		if i == what{
			amount ++
		}
	}
	return amount
}

func serve_clients(clients [] *client){
	var allClientsGroup sync.WaitGroup
	allClientsGroup.Add(1)

	finishedOrder := make(chan *client, len(clients))
	for _, cli := range clients{
		fmt.Printf("Client %d has order", cli.index)
		fmt.Print(cli.order)
		fmt.Print("\n")
		go make_order(cli, finishedOrder)
		go wait_for_order(finishedOrder)
	}
	go check_all_done(clients, &allClientsGroup)
	allClientsGroup.Wait()
}

func check_all_done(clients [] *client, threadGroup *sync.WaitGroup)  {
	done :=0
	for done == 0 {
		if servedClients >= len(clients) {
			done = 1
		}
		fmt.Print(strconv.Itoa(servedClients))
		time.Sleep(2 * time.Second)
	}
	defer threadGroup.Done()
}

func collectOrder(cli *client){
	//fmt.Print("IMPLEMANTUJ collectOrder")
	for _, product := range cli.order{
		switch product {
		case HAMBURGER:
			ready.readyHamburger--
		case CHEESEBURGER:
			ready.readyCheeseburger--
		case CHICKENBURGER:
			ready.readyChickenburger--
		case FRIES:
			ready.readyFries--
		case COLA:
			ready.readyCola--
		case NUGGETS:
			ready.readyNuggets--
			}
	}
}

func wait_for_order(finishedOrder chan *client){
	done := 0
	cli := <-finishedOrder
	for done == 0{
		mu.Lock()
		if count(HAMBURGER, cli.order) > ready.readyHamburger {
			mu.Unlock()
			time.Sleep(2*time.Second)
			fmt.Printf("Client %d is still waiting for his meal\n", cli.index)
			continue
		}
		if count(CHEESEBURGER, cli.order) > ready.readyCheeseburger {
			mu.Unlock()
			time.Sleep(2*time.Second)
			fmt.Printf("Client %d is still waiting for his meal\n", cli.index)
			continue
		}
		if count(CHICKENBURGER, cli.order) > ready.readyChickenburger {
			mu.Unlock()
			time.Sleep(2*time.Second)
			fmt.Printf("Client %d is still waiting for his meal\n", cli.index)
			continue
		}
		if count(FRIES, cli.order) > ready.readyFries {
			mu.Unlock()
			time.Sleep(2*time.Second)
			fmt.Printf("Client %d is still waiting for his meal\n", cli.index)
			continue
		}
		if count(COLA, cli.order) > ready.readyCola {
			mu.Unlock()
			time.Sleep(2*time.Second)
			fmt.Printf("Client %d is still waiting for his meal\n", cli.index)
			continue
		}
		if count(NUGGETS, cli.order) > ready.readyNuggets {
			mu.Unlock()
			time.Sleep(2*time.Second)
			fmt.Printf("Client %d is still waiting for his meal\n", cli.index)
			continue
		}
		if tray.available < 0 {
			mu.Unlock()
			time.Sleep(2*time.Second)
			fmt.Printf("Client %d is still waiting for his meal\n", cli.index)
			continue
		}
		done = 1
	}
	collectOrder(cli)
	tray.available--
	tray.inUse++
	mu.Unlock()
	cli.timeWaited = time.Since(cli.timeOfEntrance)
	fmt.Printf("%d client's order is ready\n", cli.index)
	callClient(cli)
}

func callClient(cli *client){
	for worker.caschiersAvalible ==0{
		fmt.Print("There is no cashiers\n")
		time.Sleep(1 * time.Second)
	}
	mu.Lock()
	worker.caschiersAvalible --
	worker.caschiersInUse ++
	mu.Unlock()
	fmt.Printf("Client %d is going to get his meal\n", cli.index)
	time.Sleep(3)
	mu.Lock()
	worker.caschiersAvalible ++
	worker.caschiersInUse --
	mu.Unlock()
	enjoyMeal(cli)
}

func enjoyMeal(cli *client){
	fmt.Printf("Client %d is enjoying his meal\n", cli.index)
	time.Sleep(10)
	fmt.Printf("Client %d finished his meal\n", cli.index)
	clientIsReturningTray(cli)
}

func clientIsReturningTray(cli *client)  {
	fmt.Printf("Client %d is returning his tray\n", cli.index)
	time.Sleep(3)
	fmt.Printf("Client %d returned his tray\n", cli.index)
	tray.available++
	tray.inUse--
	mu.Lock()
	servedClients++
	mu.Unlock()
	fmt.Printf("Client %d has just been served\n", cli.index)
}

func create_client(clientId int) *client{
	var order []int
	number_of_elements := rand.Intn(6)
	for i :=0 ; i<number_of_elements; i++{
		order = append(order, rand.Intn(5)+1)
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
		fmt.Printf("Client %d waited %d\n", i.index, i.timeWaited)
		average += float64(i.timeWaited)
	}
	average /= float64(len(clients))
	fmt.Printf("Average time of wainting: %f \n", average)
}

func make_order(cli *client, finishedOrder chan *client) {
	fmt.Printf("Client %d is making an order\n", cli.index)
	choice := rand.Intn(1)+1
	if choice == 0{ //checkout
		for worker.caschiersAvalible<0{
			fmt.Print("There is no cashiers\n")
			time.Sleep(1 * time.Second)
			}
			mu.Lock()
			worker.caschiersAvalible-=1
			worker.caschiersInUse+=1
			mu.Unlock()
			time.Sleep(5*time.Second)
			go zrob_zarcie(cli)
			}

	if choice == 1{ //checkout
		for worker.caschiersAvalible<0{
			fmt.Print("There is no cashiers\n")
			time.Sleep(1 * time.Second)
			}
			mu.Lock()
			worker.caschiersAvalible-=1
			worker.caschiersInUse+=1
			mu.Unlock()
			time.Sleep(5*time.Second)
			go zrob_zarcie(cli)
			}
	finishedOrder <- cli
	
	mu.Lock()
	worker.caschiersAvalible+=1
	worker.caschiersInUse-=1
	mu.Unlock()
}

func zrob_zarcie(cli *client){
	for _,ord := range(cli.order){
			switch ord {
			case HAMBURGER:
				go prepareMeal(time.Second * 3, &ready.readyHamburger)
			case CHEESEBURGER:
				go prepareMeal(time.Second * 3, &ready.readyCheeseburger)
			case CHICKENBURGER:
				go prepareMeal(time.Second * 3, &ready.readyChickenburger)
			case FRIES:
				go prepareMeal(time.Second * 3, &ready.readyFries)
			case COLA:
				go prepareMeal(time.Second * 3, &ready.readyCola)
			case NUGGETS:
				go prepareMeal(time.Second * 3, &ready.readyNuggets)
			}
		}
}

func prepareMeal(duration time.Duration, what *int){
	fmt.Printf("Preparing something, now is : %d \n", *what)
	mu.Lock()
	worker.kitchenWorkersInUse+=1
	worker.kitchenWorkersAvalible-=1
	mu.Unlock()
	time.Sleep(duration) //robi hamburgera
	mu.Lock()
	*what+=1
	mu.Unlock()
	mu.Lock()
	worker.kitchenWorkersInUse-=1
	worker.kitchenWorkersAvalible+=1
	mu.Unlock()
}

var mu = &sync.Mutex{}
var profit = 0.0
var lose = 0.0
var frying_pans = device{3,0}
var checkouts = device{3,0}
var selfCheckouts = device{3,0}
var friesMaker = device{3,0}
var tray = device{10,0}
var ready = readyMeals{0,0,0,0,0,0}
var worker = workers{5,0,5,0}
var servedClients  = 0

func main(){
	clients := make_client_queue(10)
	start_time := time.Now()
	serve_clients(clients)
	duration := time.Since(start_time)
	seconds := duration / time.Second
	fmt.Printf("Time of dealing with 10 clients %d second \n", seconds)
	get_clients_average_time(clients)
	fmt.Print("Balance: ", profit - lose, "\n")
}