package main

import (
	"math/rand"
	"fmt"
	"time"
	"math"
	"sync"
	"../util"
	"../solution"
	"strconv"
)

const(
	runtime = 60
	workerCount = 10
	startTemp = 0.05
	dt = 0.00001
)

var (
	size = 0
	solutions = solution.SolutionTree{Queen: -1}

	wg sync.WaitGroup
	mutex sync.Mutex

)

const printQueens = true

func getRandomIndex() int {
	return rand.Intn(size)
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
	queensCopy := make([]int, size)
	copy(queensCopy, queens)
	queen1 := queensCopy[queenIndex1]
	queen2 := queensCopy[queenIndex2]
	queensCopy[queenIndex1] = queen2
	queensCopy[queenIndex2] = queen1

	return queensCopy
}

func sa_search(startQueens []int) []int {
	currentQueens := startQueens
	currentEval := util.EvaluateBoard(startQueens)
	temp := startTemp
	for temp > 0 {
		if currentEval == 0 {
			return currentQueens
		}
		neighbour := generateNeighbour(currentQueens)
		neighbourEval := util.EvaluateBoard(neighbour)
		if printQueens {
			fmt.Println("Current   board(", currentEval, ")", currentQueens)
			fmt.Println("Neighbour board(", neighbourEval, ")", neighbour)
			fmt.Println("------------------------------")
		}
		if neighbourEval == 0 {
			if printQueens {
				fmt.Println("Solution board(", neighbourEval, ")", neighbour)
			}
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
	startTime := time.Now()
	for time.Since(startTime) / 1000000000 < runtime {
		queens := sa_search(queens)
		addSolution(queens)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	queens := util.GetQueensFromInput()
	size = len(queens)

	start := time.Now()

	if printQueens {
		sa_search(queens)
	} else {
		wg.Add(workerCount)
		for i := 0; i < workerCount; i++ {
			go findSolutions(queens)
		}
		wg.Wait()
	}

	elapsed := time.Since(start)
	if !printQueens {
		fmt.Println(size, "queens")
		fmt.Println(solutions.Size, "solutions")
	}
	fmt.Println("Execution time:", elapsed)

	util.PrintSolutionsToFile(solutions.GetSolutions(), "results/SA" + strconv.Itoa(size) + ".txt")

}
