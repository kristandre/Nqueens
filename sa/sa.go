package main

import (
	"math/rand"
	"fmt"
	"time"
	"math"
	"sync"
	"../util"
	"github.com/pkg/profile"
	"../solution"
)

var (
	SIZE = 0
	counter = 0
	timeInFunc = 0
	solutions = solution.SolutionTree{Queen: -1}

	SOLUTION_GOAL = 300000
	WORKER_COUNT = 10
	wg sync.WaitGroup
	mutex sync.Mutex
)

func getRandomIndex() int {
	return rand.Intn(SIZE)
}

func generateNeighbour(queens []int) []int {
	randomIndex1 := getRandomIndex()
	randomIndex2 := getRandomIndex()
	for randomIndex2 == randomIndex1 {
		randomIndex2 = getRandomIndex()
	}

	return swap(queens, randomIndex1, randomIndex2)
}

func swap(queens []int, queenIndex1, queenIndex2 int) []int {
	queensCopy := make([]int, SIZE)
	copy(queensCopy, queens)
	queen1 := queensCopy[queenIndex1]
	queen2 := queensCopy[queenIndex2]
	queensCopy[queenIndex1] = queen2
	queensCopy[queenIndex2] = queen1

	return queensCopy
}

func sa_search(startQueens []int, temp, dt float32) []int {
	currentQueens := startQueens
	currentEval := util.EvaluateBoard(startQueens)
	for temp > 0 {
		counter++
		if currentEval == 0 {
			return currentQueens
		}
		neighbour := generateNeighbour(currentQueens)
		neighbourEval := util.EvaluateBoard(neighbour)
		if neighbourEval == 0 {
			return neighbour
		}
		if neighbourEval < currentEval {
			currentQueens = neighbour
			currentEval = neighbourEval
		} else {
			q := float64(currentEval - neighbourEval) / float64(currentEval)
			p := math.Min(1.0, math.Exp(q / float64(temp)))
			x := rand.Float64()
			if p > x {
				currentQueens = neighbour
				currentEval = neighbourEval
			}
		}
		temp -= dt
	}
	return []int{}
}

func addSolution(solution []int) {
	if (len(solution) <= 0) {
		return
	}
	mutex.Lock()
	solutions.AddQueens(solution)
	for i := 0; i < 3; i++ {
		solution = util.Rotate90(solution)
		solutions.AddQueens(solution)
	}
	mutex.Unlock()
}

func findSolutions(queens []int) {
	defer wg.Done()
	for solutions.Size < SOLUTION_GOAL {
		queens := sa_search(queens, 0.05, 0.000001)
		addSolution(queens)
	}
}

func main() {
	start := time.Now()
	defer profile.Start().Stop()
	rand.Seed(time.Now().UnixNano())
	queens := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29}
	SIZE = len(queens)

	wg.Add(WORKER_COUNT)
	for i := 0; i < WORKER_COUNT; i++ {
		go findSolutions(queens)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Println(SIZE, "queens")
	fmt.Println(solutions.Size, "solutions")
	//fmt.Println(solutions.ToString())
	fmt.Println(float32(float32(timeInFunc) / 1000000.0), "in func")
	fmt.Println("Execution time:", elapsed)
	fmt.Println(counter, "rounds in sa")
	//fmt.Println(solutions)
}
