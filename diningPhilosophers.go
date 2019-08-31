// Implement the dining philosopher’s problem with the following constraints/modifications.

//     There should be 5 philosophers sharing chopsticks, with one chopstick between each adjacent pair of philosophers.
//     Each philosopher should eat only 3 times (not in an infinite loop as we did in lecture)
//     The philosophers pick up the chopsticks in any order, not lowest-numbered first (which we did in lecture).
//     In order to eat, a philosopher must get permission from a host which executes in its own goroutine.
//     The host allows no more than 2 philosophers to eat concurrently.
//     Each philosopher is numbered, 1 through 5.
//     When a philosopher starts eating (after it has obtained necessary locks) it prints “starting to eat <number>” on a line by itself, where <number> is the number of the philosopher.
//     When a philosopher finishes eating (before it has released its locks) it prints “finishing eating <number>” on a line by itself, where <number> is the number of the philosopher.
package main
import(
       "fmt"
       "sync"
       "math/rand"
       )
type ChopS struct{sync.Mutex}

type Philo struct {
  leftCS, rightCS *ChopS
  id int
  // status string
}

func host(philos *[]*Philo, host_wg *sync.WaitGroup) { 
  var wg sync.WaitGroup
  counter := 0
  defer host_wg.Done()
  for i := 0; i<5; i++ {
    wg.Add(1)
    counter++
    go (*philos)[i].eat(&wg)
    if (counter == 2 || i == 4) {
      wg.Wait()
      counter = 0
    } 
  }
}

func (p *Philo) eat(wg *sync.WaitGroup) {
    defer wg.Done()
    for i:=0; i<3; i++ {
        p.leftCS.Lock()
        p.rightCS.Lock()
        fmt.Printf("starting to eat %d\n", p.id)
        fmt.Printf("finishing eating %d\n", p.id)
        p.leftCS.Unlock()
        p.rightCS.Unlock() 
    }
}

func main() {
  // Init the array of philosophers
  CSticks := make([]*ChopS, 5)
  for i := 0; i < 5; i++ {
     CSticks[i] = new(ChopS)
  }
  philos := make([]*Philo, 5)
  var host_wg sync.WaitGroup
  for i := 0; i < 5; i++ {
      // chopstics index gen
     var chopsticsL, chopsticsR int
     chopsticsL, chopsticsR = func() (int, int) {
        lc := rand.Intn(5)
        rc := rand.Intn(5)
        for lc == rc {
          lc = rand.Intn(5)
        }
        return lc, rc
      }()
     philos[i] = &Philo{CSticks[chopsticsL], CSticks[chopsticsR], i+1}
  }
  host_wg.Add(1)
  go host(&philos, &host_wg)
  host_wg.Wait()
  fmt.Println("Press any key to exit")
  fmt.Scanln()
}