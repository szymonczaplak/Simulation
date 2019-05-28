package main

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"math/rand"
	"sync"
	"time"
)

const (
	HAMBURGER     = 1
	CHEESEBURGER  = 2
	CHICKENBURGER = 3
	FRIES         = 4
	COLA          = 5
	NUGGETS       = 6

	TIME_INSTANT        = 10 * time.Millisecond
	TIME_GETTING_MEAL   = 1000 * time.Millisecond
	TIME_EATING_MEAL    = 60000 * time.Millisecond
	TIME_RETURNING_TRAY = 2000 * time.Millisecond
	TIME_MAKING_ORDER   = 1500 * time.Millisecond

	TIME_BURGER        = 10000 * time.Millisecond
	TIME_CHEESEBURGER  = 11200 * time.Millisecond
	TIME_CHICKENBURGER = 15000 * time.Millisecond //dolozyc 2 miejsce
	TIME_FRIES         = 15000 * time.Millisecond
	TIME_COLA          = 1000 * time.Millisecond
	TIME_NUGGETSY      = 20000 * time.Millisecond
)

type client struct {
	index          int
	timeOfEntrance time.Time
	timeWaited     time.Duration
	order          []int
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
	for _, i := range order {
		if i == what {
			amount++
		}
	}
	return amount
}

func serve_clients(clients []*client) {
	var allClientsGroup sync.WaitGroup
	allClientsGroup.Add(1)

	finishedOrder := make(chan *client, len(clients))
	for _, cli := range clients {
		fmt.Printf("Client %d has order", cli.index)
		fmt.Print(cli.order)
		fmt.Print("\n")
		go make_order(cli, finishedOrder)
		go wait_for_order(finishedOrder)
	}
	go check_all_done(clients, &allClientsGroup)
	allClientsGroup.Wait()
}

func check_all_done(clients []*client, threadGroup *sync.WaitGroup) {
	done := 0
	for done == 0 {
		if servedClients >= len(clients) {
			done = 1
		}
		time.Sleep(TIME_INSTANT)
	}
	defer threadGroup.Done()
}

func collectOrder(cli *client) {
	for _, product := range cli.order {
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

func wait_for_order(finishedOrder chan *client) {
	done := 0
	cli := <-finishedOrder
	for done == 0 {
		mu.Lock()
		if count(HAMBURGER, cli.order) > ready.readyHamburger {
			mu.Unlock()
			time.Sleep(TIME_INSTANT)
			fmt.Printf("Client %d is still waiting for his meal\n", cli.index)
			continue
		}
		if count(CHEESEBURGER, cli.order) > ready.readyCheeseburger {
			mu.Unlock()
			time.Sleep(TIME_INSTANT)
			fmt.Printf("Client %d is still waiting for his meal\n", cli.index)
			continue
		}
		if count(CHICKENBURGER, cli.order) > ready.readyChickenburger {
			mu.Unlock()
			time.Sleep(TIME_INSTANT)
			fmt.Printf("Client %d is still waiting for his meal\n", cli.index)
			continue
		}
		if count(FRIES, cli.order) > ready.readyFries {
			mu.Unlock()
			time.Sleep(TIME_INSTANT)
			fmt.Printf("Client %d is still waiting for his meal\n", cli.index)
			continue
		}
		if count(COLA, cli.order) > ready.readyCola {
			mu.Unlock()
			time.Sleep(TIME_INSTANT)
			fmt.Printf("Client %d is still waiting for his meal\n", cli.index)
			continue
		}
		if count(NUGGETS, cli.order) > ready.readyNuggets {
			mu.Unlock()
			time.Sleep(TIME_INSTANT)
			fmt.Printf("Client %d is still waiting for his meal\n", cli.index)
			continue
		}
		if tray.available < 0 {
			mu.Unlock()
			time.Sleep(TIME_INSTANT)
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

func callClient(cli *client) {
	for worker.caschiersAvalible == 0 {
		//fmt.Print("There is no cashiers\n")
		time.Sleep(TIME_INSTANT)
	}
	mu.Lock()
	worker.caschiersAvalible--
	worker.caschiersInUse++
	mu.Unlock()
	fmt.Printf("Client %d is going to get his meal\n", cli.index)
	time.Sleep(TIME_GETTING_MEAL)
	mu.Lock()
	worker.caschiersAvalible++
	worker.caschiersInUse--
	mu.Unlock()
	enjoyMeal(cli)
}

func enjoyMeal(cli *client) {
	fmt.Printf("Client %d is enjoying his meal\n", cli.index)
	time.Sleep(TIME_EATING_MEAL)
	fmt.Printf("Client %d finished his meal\n", cli.index)
	clientIsReturningTray(cli)
}

func clientIsReturningTray(cli *client) {
	fmt.Printf("Client %d is returning his tray\n", cli.index)
	time.Sleep(TIME_RETURNING_TRAY)
	fmt.Printf("Client %d returned his tray\n", cli.index)
	tray.available++
	tray.inUse--
	mu.Lock()
	servedClients++
	mu.Unlock()
	fmt.Printf("Client %d has just been served\n", cli.index)
}

func create_client(clientId int) *client {
	var order []int
	number_of_elements := rand.Intn(6)
	for i := 0; i < number_of_elements; i++ {
		order = append(order, rand.Intn(5)+1)
	}
	return &client{clientId, time.Now(), 0, order}
}

func make_client_queue(number_of_clients int) []*client {
	var out []*client
	for i := 0; i < number_of_clients; i++ {
		out = append(out, create_client(i))
		fmt.Printf("Client %d entered bakery\n", i)
	}
	return out
}

func get_clients_average_time(clients []*client) float64 {
	var average float64
	for _, i := range clients {
		fmt.Printf("Client %d waited %d\n", i.index, i.timeWaited)
		average += i.timeWaited.Seconds()
	}
	out := average / float64(len(clients))
	fmt.Printf("Average time of wainting: %f \n", average)
	return out
}

func make_order(cli *client, finishedOrder chan *client) {
	fmt.Printf("Client %d is making an order\n", cli.index)
	choice := rand.Intn(1) + 1
	if choice == 0 { //checkout
		for worker.caschiersAvalible < 0 {
			//fmt.Print("There is no cashiers\n")
			time.Sleep(TIME_INSTANT)
		}
		mu.Lock()
		worker.caschiersAvalible -= 1
		worker.caschiersInUse += 1
		mu.Unlock()
		time.Sleep(TIME_MAKING_ORDER)
		go zrob_zarcie(cli)
	}

	if choice == 1 { //checkout
		for worker.caschiersAvalible < 0 {
			//fmt.Print("There is no cashiers\n")
			time.Sleep(TIME_INSTANT)
		}
		mu.Lock()
		worker.caschiersAvalible -= 1
		worker.caschiersInUse += 1
		mu.Unlock()
		time.Sleep(TIME_MAKING_ORDER)
		go zrob_zarcie(cli)
	}
	finishedOrder <- cli

	mu.Lock()
	worker.caschiersAvalible += 1
	worker.caschiersInUse -= 1
	mu.Unlock()
}

func zrob_zarcie(cli *client) {
	for _, ord := range cli.order {
		switch ord {
		case HAMBURGER:
			go prepareMeal(TIME_BURGER, &ready.readyHamburger, &frying_pans)
		case CHEESEBURGER:
			go prepareMeal(TIME_CHEESEBURGER, &ready.readyCheeseburger, &frying_pans)
		case CHICKENBURGER:
			go prepareMeal(TIME_CHICKENBURGER, &ready.readyChickenburger, &frying_pans)
		case FRIES:
			go prepareMeal(TIME_FRIES, &ready.readyFries, &friesMaker)
		case COLA:
			go prepareMeal(TIME_COLA, &ready.readyCola, &colaDispenser)
		case NUGGETS:
			go prepareMeal(TIME_NUGGETSY, &ready.readyNuggets, &frying_pans)
		}
	}
}

func prepareMeal(duration time.Duration, what *int, dev *device) {
	fmt.Printf("Preparing something, now is : %d \n", *what)
	for worker.kitchenWorkersAvalible == 0 || dev.available == 0 {
		//fmt.Print("There is no kitchen workers\n")
		time.Sleep(20)
	}
	mu.Lock()
	worker.kitchenWorkersInUse += 1
	worker.kitchenWorkersAvalible -= 1
	dev.available--
	dev.inUse++
	mu.Unlock()
	time.Sleep(duration) //robi hamburgera
	mu.Lock()
	*what += 1
	mu.Unlock()
	mu.Lock()
	worker.kitchenWorkersInUse -= 1
	worker.kitchenWorkersAvalible += 1
	dev.available--
	dev.inUse++
	mu.Unlock()
}

func make_simulation(clients []*client) (float64, float64) {
	start_time := time.Now()
	serve_clients(clients)
	duration := time.Since(start_time).Seconds()
	fmt.Printf("Time of dealing with 10 clients %f second \n", duration)
	average := get_clients_average_time(clients)
	fmt.Print("Balance: ", profit-lose, "\n")
	return duration, average
}

func plot_results(sim_time []float64, parameter []int, name string) {
	xys := make(plotter.XYs, len(sim_time))
	for ind, v := range sim_time {
		xys[ind].X = float64(parameter[ind])
		xys[ind].Y = v
	}
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Wykres"
	p.X.Label.Text = "number of chefs"
	p.Y.Label.Text = "time [s]"

	err = plotutil.AddLinePoints(p,
		"time of serving 10 clients", xys)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, name); err != nil {
		panic(err)
	}
}

func reset_clients_timer(clients []*client) {
	for _, cli := range clients {
		cli.timeOfEntrance = time.Now()
	}
}

func reset_all_variables(change int) {
	frying_pans = device{8, 0}
	friesMaker = device{8, 0}
	colaDispenser = device{8, 0}
	tray = device{change, 0}
	ready = readyMeals{0, 0, 0, 0, 0, 0}
	worker = workers{10, 0, 10, 0}
	servedClients = 0
}

var mu = &sync.Mutex{}
var profit = 0.0
var lose = 0.0
var frying_pans = device{8, 0}
var checkouts = device{8, 0}
var selfCheckouts = device{8, 0}
var friesMaker = device{8, 0}
var colaDispenser = device{8, 0}
var tray = device{4, 0}
var ready = readyMeals{0, 0, 0, 0, 0, 0}
var worker = workers{10, 0, 10, 0}
var servedClients = 0

func main() {
	var simulation_time []float64
	var parameter []int
	var avrg_customer_service_time []float64
	clients := make_client_queue(10)

	//--------------Simulation with initial values -----------------
	parameter = append(parameter, 4)
	reset_clients_timer(clients)
	st, acst := make_simulation(clients)
	simulation_time = append(simulation_time, st)
	avrg_customer_service_time = append(avrg_customer_service_time, acst)

	time.Sleep(1 * time.Second)

	parameter = append(parameter, 6)
	reset_all_variables(6)
	servedClients = 0
	reset_clients_timer(clients)
	st, acst = make_simulation(clients)
	simulation_time = append(simulation_time, st)
	avrg_customer_service_time = append(avrg_customer_service_time, acst)

	time.Sleep(1 * time.Second)

	parameter = append(parameter, 7)
	reset_all_variables(7)
	servedClients = 0
	reset_clients_timer(clients)
	st, acst = make_simulation(clients)
	simulation_time = append(simulation_time, st)
	avrg_customer_service_time = append(avrg_customer_service_time, acst)

	time.Sleep(1 * time.Second)

	parameter = append(parameter, 8)
	reset_all_variables(8)
	servedClients = 0
	reset_clients_timer(clients)
	st, acst = make_simulation(clients)
	simulation_time = append(simulation_time, st)
	avrg_customer_service_time = append(avrg_customer_service_time, acst)

	time.Sleep(1 * time.Second)

	parameter = append(parameter, 10)
	reset_all_variables(10)
	servedClients = 0
	reset_clients_timer(clients)
	st, acst = make_simulation(clients)
	simulation_time = append(simulation_time, st)
	avrg_customer_service_time = append(avrg_customer_service_time, acst)

	time.Sleep(1 * time.Second)

	parameter = append(parameter, 15)
	reset_all_variables(15)
	servedClients = 0
	reset_clients_timer(clients)
	st, acst = make_simulation(clients)
	simulation_time = append(simulation_time, st)
	avrg_customer_service_time = append(avrg_customer_service_time, acst)

	time.Sleep(1 * time.Second)

	parameter = append(parameter, 20)
	reset_all_variables(20)
	servedClients = 0
	reset_clients_timer(clients)
	st, acst = make_simulation(clients)
	simulation_time = append(simulation_time, st)
	avrg_customer_service_time = append(avrg_customer_service_time, acst)

	fmt.Print("---------------------\n")
	fmt.Print(simulation_time)
	fmt.Print("\n")
	fmt.Print(avrg_customer_service_time)
	plot_results(avrg_customer_service_time, parameter, "average_customer_waiting_time_trays.png")
	plot_results(simulation_time, parameter, "simulation_time_trays.png")
}
