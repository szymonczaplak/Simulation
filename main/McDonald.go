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
		go make_order(cli, finishedOrder)
		go wait_for_order(finishedOrder)
	}
	go check_all_done(clients, &allClientsGroup)
	allClientsGroup.Wait()
}

func check_all_done(clients [] *client, threadGroup *sync.WaitGroup)  {
	done :=0
	for done == 0 {
		if servedClients == len(clients) {
			done = 1
		}
	}
	defer threadGroup.Done()
}

func wait_for_order(finishedOrder chan *client){
	done := 0
	cli := <-finishedOrder
	for done == 0{
		fmt.Printf("Client %d is still waiting for his meal\n", cli.index)
		mu.Lock()
		if count(HAMBURGER, cli.order) > ready.readyHamburger {
			mu.Unlock()
			time.Sleep(2)
			continue
		}
		if count(CHEESEBURGER, cli.order) > ready.readyCheeseburger {
			mu.Unlock()
			time.Sleep(2)
			continue
		}
		if count(CHICKENBURGER, cli.order) > ready.readyChickenburger {
			mu.Unlock()
			time.Sleep(2)
			continue
		}
		if count(FRIES, cli.order) > ready.readyFries {
			mu.Unlock()
			time.Sleep(2)
			continue
		}
		if count(COLA, cli.order) > ready.readyCola {
			mu.Unlock()
			time.Sleep(2)
			continue
		}
		if count(NUGGETS, cli.order) > ready.readyNuggets {
			mu.Unlock()
			time.Sleep(2)
			continue
		}
		if tray.available >0{
			tray.available--
			tray.inUse++
		}
		mu.Unlock()
		cli.timeWaited = time.Since(cli.timeOfEntrance)
		done = 1
	}
	fmt.Printf("Client %d got his meal\n", cli.index)
	go callClient(cli)
}

func callClient(cli *client){
	for worker.caschiersAvalible <0{
	}
	worker.caschiersAvalible --
	worker.caschiersInUse ++
	fmt.Printf("Client %d is going to get his meal\n", cli.index)
	time.Sleep(3)
	go enjoyMeal(cli)
}

func enjoyMeal(cli *client){
	fmt.Printf("Client %d is enjoying his meal\n", cli.index)
	time.Sleep(10)
	fmt.Printf("Client %d finished his meal\n", cli.index)
	go clientIsReturningTray(cli)
}

func clientIsReturningTray(cli *client)  {
	fmt.Printf("Client %d is returning his tray\n", cli.index)
	time.Sleep(3)
	fmt.Printf("Client %d returned his tray\n", cli.index)
	tray.available++
	tray.available--
	servedClients++
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

func make_order(cli *client, finishedOrder chan *client) {
	choice := rand.Intn(1)+1
	if choice == 0{ //checkout
		for worker.caschiersAvalible<0{
			}
			worker.caschiersAvalible-=1
			worker.caschiersInUse+=1
			time.Sleep(5*time.Second)
			go zrob_zarcie(cli)
			}

	if choice == 1{ //checkout
		for worker.caschiersAvalible<0{
			}
			worker.caschiersAvalible-=1
			worker.caschiersInUse+=1
			time.Sleep(5*time.Second)
			go zrob_zarcie(cli)
			}
	finishedOrder <- cli
}


func zrob_zarcie(cli *client){
	for _,ord := range(cli.order){
			switch ord {
			case HAMBURGER:
				go doHamburger()
			case CHEESEBURGER:
				go doCheeseburger()
			case CHICKENBURGER:
				go doChickenburger()
			case FRIES:
				go doFries()
			case COLA:
				doCola()
			case NUGGETS:
				go doNuggets()
			}
		}

}

func doHamburger(){ //dodac while
	for ready.readyHamburger<0{
	}
	worker.kitchenWorkersInUse+=1;
	worker.kitchenWorkersAvalible-=1;
	time.Sleep(10 * time.Second) //robi hamburgera
	ready.readyHamburger+=1
}

func doCheeseburger(){ //dodac while
	for ready.readyCheeseburger<0 {
	}
	worker.kitchenWorkersInUse+=1;
	worker.kitchenWorkersAvalible-=1;
	time.Sleep(10 * time.Second) //robi hamburgera
	ready.readyCheeseburger+=1
}

func doChickenburger(){ //dodac while
	for ready.readyChickenburger<0 {
	}
	worker.kitchenWorkersInUse+=1;
	worker.kitchenWorkersAvalible-=1;
	time.Sleep(10 * time.Second) //robi hamburgera
	ready.readyChickenburger+=1

}

func doFries(){ //dodac while
	for ready.readyFries<0 {
	}
	worker.kitchenWorkersInUse+=1;
	worker.kitchenWorkersAvalible-=1;
	time.Sleep(10 * time.Second) //robi hamburgera
	ready.readyFries+=1

}

func doNuggets(){ //dodac while
	for ready.readyNuggets<0 {
	}
	worker.kitchenWorkersInUse+=1;
	worker.kitchenWorkersAvalible-=1;
	time.Sleep(10 * time.Second) //robi hamburgera
	ready.readyNuggets+=1

}

func doCola(){ //dodac while
	for ready.readyCola<0 {
	}
	worker.kitchenWorkersInUse+=1;
	worker.kitchenWorkersAvalible-=1;
	time.Sleep(10 * time.Second) //robi hamburgera
	ready.readyCola+=1
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
var worker = workers{3,0,3,0}
var servedClients  = 0

func main(){
	clients := make_client_queue(10)
	serve_clients(clients)
	get_clients_average_time(clients)
	fmt.Print("Balance: ", profit - lose, "\n")
}