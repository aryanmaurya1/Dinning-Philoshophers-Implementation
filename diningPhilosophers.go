package main

import (
	"fmt"
	"sync"
)

const maxEat = 3
const maxPermissions = 5 * maxEat
const maxSimultaneousEat = 2

var permissionMutex sync.Mutex

type chopStick struct {
	sync.Mutex
}

type philosopher struct {
	left  *chopStick
	right *chopStick
	id    int
}

func (p *philosopher) eat() {
	fmt.Printf("starting to eat %d\n", p.id+1)
	fmt.Printf("finishing eating %d\n", p.id+1)
}

func (p *philosopher) askPermission(c chan int, per chan bool) {
	c <- p.id
	per <- true
}

func (p *philosopher) done(d chan bool) {
	d <- true
}

func (p *philosopher) Lock() {
	p.left.Lock()
	p.right.Lock()
}

func (p *philosopher) Unlock() {
	p.left.Unlock()
	p.right.Unlock()
}

func (p *philosopher) run(c chan int, per chan bool, done chan bool, wg *sync.WaitGroup) {
	var i int
	for {
		permissionMutex.Lock()
		p.askPermission(c, per)
		status := <-per
		<-c
		if status {
			p.Lock()
			permissionMutex.Unlock()
			p.eat()
			i++
			p.Unlock()
			p.done(done)
		} else {
			permissionMutex.Unlock()
		}
		if i == maxEat {
			fmt.Println("\nPhilosopher : ", p.id+1, " done 	\n")
			break
		}
	}
	wg.Done()
}

func host(c chan int, per chan bool, d chan bool) {
	var permissionTable = make([]int, 5, 5)
	var currentPGranted int
	var totalPGranted int
	for i := 0; i < maxPermissions; i++ {
		status := <-per
		id := <-c
		if status && (currentPGranted < maxSimultaneousEat) && (permissionTable[id] < maxEat) {
			permissionTable[id]++
			currentPGranted++
			totalPGranted++
			per <- true
			c <- id
		} else {
			per <- false
			c <- id
		}
		var done bool
		select {
		case done = <-d:
			if done {
				currentPGranted--
			}
		default:

		}
	}
}

func main() {

	var wg sync.WaitGroup
	var cstics = make([]*chopStick, 5, 5)
	var philo = make([]*philosopher, 5, 5)
	for i := 0; i < 5; i++ {
		cstics[i] = new(chopStick)
	}
	for i := 0; i < 5; i++ {
		philo[i] = &philosopher{cstics[i], cstics[(i+1)%5], i}
	}

	var per = make(chan bool, 100)
	var c = make(chan int, 100)
	var d = make(chan bool, 100)

	wg.Add(5)
	go host(c, per, d)
	for i := 0; i < 5; i++ {
		go philo[i].run(c, per, d, &wg)
	}
	wg.Wait()
}
