package main

import (
	"math/rand"
	"sort"
	"fmt"
	"time"
	"../util"
	"github.com/pkg/profile"
	"sync"
	"../solution"
)

const (
	MUTATION_PROB float32 = 0.001
	RUNTIME = 6
)

var (
	solutions = solution.SolutionTree{Queen: -1}
	TOURNAMENT_SIZE int = 20
	POPULATION_SIZE int = 44

	SIZE int = 0
	availableRows = []int{}
	WORKER_COUNT = 10
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
	randomIndex1 := rand.Intn(SIZE)
	randomIndex2 := rand.Intn(SIZE)
	for randomIndex2 == randomIndex1 {
		randomIndex2 = rand.Intn(SIZE)
	}
	b.queens[randomIndex1], b.queens[randomIndex2] = b.queens[randomIndex2], b.queens[randomIndex1]
	queens := b.queens
	for i, q := range queens {
		if q + randomIndex1 > SIZE - 1 {
			queens[i] = (q + randomIndex1) - SIZE
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
		mom.queens = rand.Perm(SIZE)
		mom.Evaluate()
	}
	worst := gladiators[2:]
	children := make(Population, len(worst))
	for _, glad := range worst {
		child := tournamentCrossover(mom.queens, dad.queens)
		if child.eval != 0 && rand.Float32() < MUTATION_PROB {
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
	ar := make([]int, SIZE)
	copy(ar, availableRows)

	for i := 0; i < SIZE; i++ {
		if mom[i] == dad[i] {
			ar = util.DeleteItem(dad[i], ar)
		}
	}

	for i := 0; i < SIZE; i++ {
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
	for time.Since(startTime) / 1000000000 < RUNTIME {
		tournamentSelection(TOURNAMENT_SIZE, population)
	}
}

func createMutatedPopulation(queens []int) Population {
	mutatedPopulation := make(Population, POPULATION_SIZE)
	for i := 0; i < POPULATION_SIZE; i++ {
		queensCopy := make([]int, SIZE)
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
	defer profile.Start().Stop()
	rand.Seed(time.Now().UnixNano())
	//queens := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29}
	//SIZE = len(queens)
	SIZE = 20
	availableRows = util.GetAvailableRows(SIZE)
	//printPop()
	SAMPLE_SIZE := 10
	start := time.Now()
	for k := 0; k < 10; k++ {
		POPULATION_SIZE += 1
		totalSolutions := 0
		for j := 0; j < SAMPLE_SIZE; j++ {
			queens := rand.Perm(SIZE)
			wg.Add(WORKER_COUNT)
			for i := 0; i < WORKER_COUNT; i++ {
				go findSolutions(queens)
			}
			wg.Wait()
			//fmt.Println("TOURNAMENT_SIZE", TOURNAMENT_SIZE, "Run", j + 1, ":", solutions.Size, "solutions found")
			totalSolutions += solutions.Size
			solutions = solution.SolutionTree{Queen:-1}
		}
		fmt.Println("POPULATION_SIZE", POPULATION_SIZE, "Avg solutions:", totalSolutions / SAMPLE_SIZE)
	}
	elapsed := time.Since(start)
	fmt.Println(SIZE, "queens")
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

//func printPop() {
//	for _, b := range population {
//		fmt.Println(b.queens)
//	}
//}
