package main

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"sync"
	"time"
)

const (
	Hamburger     = 1
	CHEESEBURGER  = 2
	CHICKENBURGER = 3
	FRIES         = 4
	COLA          = 5
	NUGGETS       = 6

	TimeInstant       = 10 * time.Millisecond
	TimeGettingMeal   = 1000 * time.Millisecond
	TimeEatingMeal    = 6000 * time.Millisecond
	TimeReturningTray = 2000 * time.Millisecond
	TimeMakingOrder   = 1500 * time.Millisecond

	TimeBurger        = 10000 * time.Millisecond
	TimeCheeseburger  = 11200 * time.Millisecond
	TimeChickenburger = 15000 * time.Millisecond //dolozyc 2 miejsce
	TimeFries         = 15000 * time.Millisecond
	TimeCola          = 1000 * time.Millisecond
	TimeNuggetsy      = 20000 * time.Millisecond
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

func serveClients(clients []*client) {
	var allClientsGroup sync.WaitGroup
	allClientsGroup.Add(1)

	finishedOrder := make(chan *client, len(clients)+3)
	for _, cli := range clients {
		fmt.Printf("Client %d has order", cli.index)
		fmt.Print(cli.order)
		fmt.Print("\n")
		go makeOrder(cli, finishedOrder)
		go waitForOrder(finishedOrder)
	}
	go checkAllDone(clients, &allClientsGroup)
	allClientsGroup.Wait()
}

func checkAllDone(clients []*client, threadGroup *sync.WaitGroup) {
	done := 0
	for done == 0 {
		if servedClients >= len(clients) {
			done = 1
		}
		time.Sleep(TimeInstant)
	}
	defer threadGroup.Done()
}

func collectOrder(cli *client) {
	for _, product := range cli.order {
		switch product {
		case Hamburger:
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

func waitForOrder(finishedOrder chan *client) {
	done := 0
	cli := <-finishedOrder
	for done == 0 {
		mu.Lock()
		if count(Hamburger, cli.order) > ready.readyHamburger {
			mu.Unlock()
			time.Sleep(TimeInstant)
			//			fmt.Printf("Client %d is still waiting for his meal\n", cli.index)
			continue
		}
		if count(CHEESEBURGER, cli.order) > ready.readyCheeseburger {
			mu.Unlock()
			time.Sleep(TimeInstant)
			//			fmt.Printf("Client %d is still waiting for his meal\n", cli.index)
			continue
		}
		if count(CHICKENBURGER, cli.order) > ready.readyChickenburger {
			mu.Unlock()
			time.Sleep(TimeInstant)
			//			fmt.Printf("Client %d is still waiting for his meal\n", cli.index)
			continue
		}
		if count(FRIES, cli.order) > ready.readyFries {
			mu.Unlock()
			time.Sleep(TimeInstant)
			//			fmt.Printf("Client %d is still waiting for his meal\n", cli.index)
			continue
		}
		if count(COLA, cli.order) > ready.readyCola {
			mu.Unlock()
			time.Sleep(TimeInstant)
			//			fmt.Printf("Client %d is still waiting for his meal\n", cli.index)
			continue
		}
		if count(NUGGETS, cli.order) > ready.readyNuggets {
			mu.Unlock()
			time.Sleep(TimeInstant)
			//			fmt.Printf("Client %d is still waiting for his meal\n", cli.index)
			continue
		}
		if tray.available < 0 {
			mu.Unlock()
			time.Sleep(TimeInstant)
			//			fmt.Printf("Client %d is still waiting for his meal\n", cli.index)
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
		time.Sleep(TimeInstant)
	}
	mu.Lock()
	worker.caschiersAvalible--
	worker.caschiersInUse++
	mu.Unlock()
	fmt.Printf("Client %d is going to get his meal\n", cli.index)
	time.Sleep(TimeGettingMeal)
	mu.Lock()
	worker.caschiersAvalible++
	worker.caschiersInUse--
	mu.Unlock()
	enjoyMeal(cli)
}

func enjoyMeal(cli *client) {
	fmt.Printf("Client %d is enjoying his meal\n", cli.index)
	time.Sleep(TimeEatingMeal)
	fmt.Printf("Client %d finished his meal\n", cli.index)
	clientIsReturningTray(cli)
}

func clientIsReturningTray(cli *client) {
	fmt.Printf("Client %d is returning his tray\n", cli.index)
	time.Sleep(TimeReturningTray)
	fmt.Printf("Client %d returned his tray\n", cli.index)
	tray.available++
	tray.inUse--
	mu.Lock()
	servedClients++
	mu.Unlock()
	fmt.Printf("Client %d has just been served\n", cli.index)
}

func createClient(clientId int) *client {
	var order []int
	numberOfElements := 4
	for i := 0; i < numberOfElements; i++ {
		order = append(order, i+1)
	}
	return &client{clientId, time.Now(), 0, order}
}

func makeClientQueue(number_of_clients int) []*client {
	var out []*client
	for i := 0; i < number_of_clients; i++ {
		out = append(out, createClient(i))
		fmt.Printf("Client %d entered bakery\n", i)
	}
	return out
}

func getClientsAverageTime(clients []*client) float64 {
	var average float64
	for _, i := range clients {
		fmt.Printf("Client %d waited %d\n", i.index, i.timeWaited)
		average += i.timeWaited.Seconds()
	}
	out := average / float64(len(clients))
	fmt.Printf("Average time of wainting: %f \n", average)
	return out
}

func makeOrder(cli *client, finishedOrder chan *client) {
	fmt.Printf("Client %d is making an order\n", cli.index)
	for worker.caschiersAvalible < 0 {
		//fmt.Print("There is no cashiers\n")
		time.Sleep(TimeInstant)
	}
	mu.Lock()
	worker.caschiersAvalible -= 1
	worker.caschiersInUse += 1
	mu.Unlock()
	time.Sleep(TimeMakingOrder)
	go makeMeal(cli)
	mu.Lock()
	finishedOrder <- cli
	mu.Unlock()
	mu.Lock()
	worker.caschiersAvalible += 1
	worker.caschiersInUse -= 1
	mu.Unlock()
}

func makeMeal(cli *client) {
	for _, ord := range cli.order {
		switch ord {
		case Hamburger:
			go prepareMeal(TimeBurger, &ready.readyHamburger, &fryingPans)
		case CHEESEBURGER:
			go prepareMeal(TimeCheeseburger, &ready.readyCheeseburger, &fryingPans)
		case CHICKENBURGER:
			go prepareMeal(TimeChickenburger, &ready.readyChickenburger, &fryingPans)
		case FRIES:
			go prepareMeal(TimeFries, &ready.readyFries, &friesMaker)
		case COLA:
			go prepareMeal(TimeCola, &ready.readyCola, &colaDispenser)
		case NUGGETS:
			go prepareMeal(TimeNuggetsy, &ready.readyNuggets, &fryingPans)
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
	dev.available++
	dev.inUse--
	mu.Unlock()
}

func makeSimulation(clients []*client) (float64, float64) {
	startTime := time.Now()
	serveClients(clients)
	duration := time.Since(startTime).Seconds()
	fmt.Printf("Time of dealing with 10 clients %f second \n", duration)
	average := getClientsAverageTime(clients)
	fmt.Print("Balance: ", profit-lose, "\n")
	return duration, average
}

func plotResults(simTime []float64, parameter []int, name string) {
	xys := make(plotter.XYs, len(simTime))
	for ind, v := range simTime {
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

func resetClientsTimer(clients []*client) {
	for _, cli := range clients {
		cli.timeOfEntrance = time.Now()
	}
}

func setOrders(clients []*client) {
	for i, cli := range clients {
		cli.order = orders[i]
	}
}

func resetAllVariables(change int) {
	fryingPans = device{8, 0}
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
var fryingPans = device{8, 0}

//var checkouts = device{8, 0
//var selfCheckouts = device{8, 0}
var friesMaker = device{8, 0}
var colaDispenser = device{8, 0}
var tray = device{10, 0}
var ready = readyMeals{0, 0, 0, 0, 0, 0}
var worker = workers{10, 0, 10, 0}
var servedClients = 0
var orders = [][]int{{1, 2, 3}, {4, 5}, {5}, {4, 2, 1, 2}, {1}, {4, 2, 6}, {6, 5}, {4, 1, 6, 3}, {3, 5, 1, 1}, {5, 4, 6, 3}}

func main() {
	var simulationTime []float64
	//var parameter []int
	var avrgCustomerServiceTime []float64
	clients := makeClientQueue(10)

	//--------------Simulation with initial values -----------------
	resetClientsTimer(clients)
	setOrders(clients)
	duration, averageDuration := makeSimulation(clients)
	simulationTime = append(simulationTime, duration)
	avrgCustomerServiceTime = append(avrgCustomerServiceTime, averageDuration)

	fmt.Print("------------------\n")
	fmt.Print(simulationTime)
	fmt.Print("\n")
	fmt.Print(avrgCustomerServiceTime)
	//plot_results(avrg_customer_service_time, parameter, "average_customer_waiting_time_trays.png")
	//plot_results(simulation_time, parameter, "simulation_time_trays.png")
}
