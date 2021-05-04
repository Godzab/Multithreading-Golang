package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const matrixsize = 250

var (
	rwLock    = sync.RWMutex{}
	cond      = sync.NewCond(rwLock.RLocker())
	waitGroup = sync.WaitGroup{}
	matrixA   = [matrixsize][matrixsize]int{}
	matrixB   = [matrixsize][matrixsize]int{}
	result    = [matrixsize][matrixsize]int{}
)

func generateRandomMatrix(matrix *[matrixsize][matrixsize]int) {
	for row := 0; row < matrixsize; row++ {
		for col := 0; col < matrixsize; col++ {
			matrix[row][col] += rand.Intn(10) - 5
		}
	}
}

func workOutRow(row int) {
	rwLock.RLock()
	for {
		waitGroup.Done()
		cond.Wait()
		for col := 0; col < matrixsize; col++ {
			for i := 0; i < matrixsize; i++ {
				result[row][col] += matrixA[row][i] * matrixB[i][col]
			}
		}
	}

}

func main() {
	fmt.Println("....WORKING...")
	waitGroup.Add(matrixsize)
	for row := 0; row < matrixsize; row++ {
		go workOutRow(row)
	}
	start := time.Now()
	for i := 0; i < 100; i++ {
		waitGroup.Wait()
		rwLock.Lock()
		generateRandomMatrix(&matrixA)
		generateRandomMatrix(&matrixB)
		waitGroup.Add(matrixsize)
		rwLock.Unlock()
		cond.Broadcast()
	}
	elapsed := time.Since(start)
	fmt.Println("....DONE...")
	fmt.Println("Processing took:", elapsed)
}
