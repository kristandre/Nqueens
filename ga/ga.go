package main

import (
	"math/rand"
	"sort"
	"fmt"
	"time"
	"../util"
	"sync"
	"../solution"
)

const (
	mutationProb float32 = 0.001
	runtime = 30
	tournamentSize int = 20
	populationSize int = 50
	workerCount = 10
)

var (
	solutions = solution.SolutionTree{Queen: -1}

	size int = 0
	availableRows = []int{}
	wg sync.WaitGroup
	mutex sync.Mutex
)

type Board struct {
	queens []int
	eval   int
}
type Population []Board

func (b *Board) Evaluate() int {
	crashCount := util.EvaluateBoard(b.queens)
	b.eval = crashCount
	if crashCount == 0 {
		addSolution(b.queens)
	}
	return crashCount
}

func (b *Board) Mutate() {
	randomIndex1 := rand.Intn(size)
	randomIndex2 := rand.Intn(size)
	for randomIndex2 == randomIndex1 {
		randomIndex2 = rand.Intn(size)
	}
	b.queens[randomIndex1], b.queens[randomIndex2] = b.queens[randomIndex2], b.queens[randomIndex1]
	queens := b.queens
	for i, q := range queens {
		if q + randomIndex1 > size - 1 {
			queens[i] = (q + randomIndex1) - size
		} else {
			queens[i] = q + randomIndex1
		}
	}
}

func (b Population) Len() int {
	return len(b)
}
func (b Population) Less(i, j int) bool {
	return b[i].eval < b[j].eval
}
func (b Population) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
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

func tournamentSelection(k int, population Population) Population {
	gladiators := ShufflePopulation(k, population)
	sort.Sort(gladiators)
	mom := gladiators[0]
	dad := gladiators[1]
	if util.QueensEqual(mom.queens, dad.queens) {
		mom.queens = rand.Perm(size)
		mom.Evaluate()
	}
	worst := gladiators[2:]
	children := make(Population, len(worst))
	for _, glad := range worst {
		child := tournamentCrossover(mom.queens, dad.queens)
		if child.eval != 0 && rand.Float32() < mutationProb {
			child.Mutate()
			child.Evaluate()
		}
		children = append(children, child)
		population[PopulationIndexOf(glad, population)] = child
	}
	return children
}

func tournamentCrossover(mom, dad []int) Board {
	childQueens := make([]int, len(mom))
	ar := make([]int, size)
	copy(ar, availableRows)

	for i := 0; i < size; i++ {
		if mom[i] == dad[i] {
			ar = util.DeleteItem(dad[i], ar)
		}
	}

	for i := 0; i < size; i++ {
		if mom[i] == dad[i] {
			childQueens[i] = dad[i]
		} else {
			randomIndex := rand.Intn(len(ar))
			childQueens[i] = ar[randomIndex]
			ar = append(ar[:randomIndex], ar[randomIndex + 1:]...)
		}
	}

	child := Board{queens:childQueens, eval: 999999}
	child.Evaluate()
	return child
}

func findSolutions(queens []int) {
	defer wg.Done()
	startTime := time.Now()
	population := createMutatedPopulation(queens)
	for time.Since(startTime) / 1000000000 < runtime {
		tournamentSelection(tournamentSize, population)
	}
}

func createMutatedPopulation(queens []int) Population {
	mutatedPopulation := make(Population, populationSize)
	for i := 0; i < populationSize; i++ {
		queensCopy := make([]int, size)
		copy(queensCopy, queens)
		board := Board{queens: queensCopy}
		board.Mutate()
		board.Evaluate()
		if board.eval == 0 {
			addSolution(board.queens)
		}
		mutatedPopulation[i] = board
	}
	return mutatedPopulation
}

func main() {
	rand.Seed(time.Now().UnixNano())
	queens := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29}
	size = len(queens)
	availableRows = util.GetAvailableRows(size)
	start := time.Now()
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go findSolutions(queens)
	}
	wg.Wait()
	//fmt.Println("TOURNAMENT_SIZE", TOURNAMENT_SIZE, "Run", j + 1, ":", solutions.Size, "solutions found")
	elapsed := time.Since(start)
	fmt.Println(size, "queens")
	fmt.Println(solutions.Size, "solutions")
	fmt.Println("Execution time:", elapsed)

}

func PopulationIndexOf(board Board, population Population) int {
	for i, b := range population {
		if util.QueensEqual(b.queens, board.queens) {
			return i
		}
	}
	return -1
}

func ShufflePopulation(k int, population Population) Population {
	dest := make(Population, k)
	indices := []int{}
	for len(indices) < k {
		r := rand.Intn(len(population))
		if util.IndexOf(r, indices) < 0 {
			indices = append(indices, r)
		}
	}
	for i, v := range indices {
		dest[i] = population[v]
	}
	return dest
}
