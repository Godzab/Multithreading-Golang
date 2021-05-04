package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	lock1 = sync.Mutex{}
	lock2 = sync.Mutex{}
)

func main() {
	go blueRobot()
	go redRobot()
	time.Sleep(20 * time.Second)
}

func blueRobot() {
	for {
		fmt.Println("Blue is acquiring the lock 1")
		lock1.Lock()
		fmt.Println("Blue is acquiring the lock 2")
		lock2.Lock()
		fmt.Println("Blue has acquired both locks")
		lock1.Unlock()
		lock2.Unlock()
		fmt.Println("Blue has released locks")
	}
}

func redRobot() {
	for {
		fmt.Println("Red is acquiring the lock 2")
		lock2.Lock()
		fmt.Println("Red is acquiring the lock 2")
		lock1.Lock()
		fmt.Println("Red has acquired both locks")
		lock2.Unlock()
		lock1.Unlock()
		fmt.Println("Red has released locks")
	}
}
