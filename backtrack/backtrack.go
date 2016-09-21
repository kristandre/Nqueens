package main

import (
	"fmt"
	"time"
	"sync"
	"github.com/kristandre/Nqueens/util"
	"github.com/kristandre/Nqueens/solution"
)

func isValidPosition(queens []int, row, col int) bool {
	k := 1
	for i := col -1; i >= 0; i-- {
		queenRow := queens[i]
		if queenRow == row + k || queenRow == row - k{
			return false
		}
		k++
	}
	return true
}

var solutions = solution.SolutionTree{Queen:-1}
var mutex sync.Mutex


func addSolution(solution []int) {
	if (len(solution) <= 0) {
		return
	}
	mutex.Lock()
	solutions.AddQueens(solution)
	mutex.Unlock()
}

func backtrackOneSolution(availableRows, queens []int, col, start int) []int {
	size := len(queens)
	for _,row := range availableRows {
		if isValidPosition(queens, row, col){
			queensCopy := make([]int, size)
			availableRowsCopy := make([]int, len(availableRows))
			copy(queensCopy, queens)
			copy(availableRowsCopy, availableRows)

			availableRowsCopy = util.DeleteItem(row, availableRowsCopy)
			queensCopy[col] = row

			if col == size - 1 {
				return queensCopy
			} else {
				solutionQueens := backtrackOneSolution(availableRowsCopy, queensCopy, col + 1, start)
				if len(solutionQueens) > 0 {
					return solutionQueens
				}
			}
		}
	}
	return []int{} // No solution found
}


func backtrackAllSolutionsMT(availableRows, queens []int, col, start int, wg *sync.WaitGroup) {
	size := len(queens)
	for _,row := range availableRows {
		if isValidPosition(queens, row, col){
			queensCopy := make([]int, size)
			availableRowsCopy := make([]int, len(availableRows))
			copy(queensCopy, queens)
			copy(availableRowsCopy, availableRows)

			availableRowsCopy = util.DeleteItem(row, availableRowsCopy)
			queensCopy[col] = row

			if col == size - 1 {
				addSolution(queensCopy)
			} else {
				backtrackAllSolutionsMT(availableRowsCopy, queensCopy, col + 1, start, wg)
			}
		}
	}
	if col == start {
		wg.Done()
	}
}

func mtbt(availableRows, queens []int, col int){
	N := len(queens)
	var wg sync.WaitGroup
	for _,row := range availableRows {
		if !isValidPosition(queens, row, col){
			continue
		}
		queensCopy := make([]int, N)
		availableRowsCopy := make([]int, len(availableRows))
		copy(queensCopy, queens)
		copy(availableRowsCopy, availableRows)
		availableRowsCopy = util.DeleteItem(row, availableRowsCopy)
		queensCopy[col] = row
		wg.Add(1)
		go backtrackAllSolutionsMT(availableRowsCopy, queensCopy, col + 1, col + 1, &wg)
	}
	wg.Wait()
}

func main(){
	queens := []int{1, 3, 5, 7, -1, -1, -1, -1}
	//queens := []int{0, 6, 1, 5, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}
	availableRows := getAvailableRows(queens)
	firstEmptyColumn := getFirstEmptyColumn(queens)
	start := time.Now()
	//mtbt(availableRows, queens, firstEmptyColumn)
	oneSolution := backtrackOneSolution(availableRows, queens, firstEmptyColumn, firstEmptyColumn)
	elapsed := time.Since(start)
	fmt.Println(len(queens), "queens")
	fmt.Println("solution", oneSolution)
	//fmt.Println(solutions.Size, "solutions")
	fmt.Println("Execution time:", elapsed)
}

func getAvailableRows(queens []int) []int {
	N := len(queens)
	availableRows := make([]int, N)
	for i := 0; i < N; i++ {
		availableRows[i] = i
	}
	for _,queenRow := range queens {
		if queenRow >= 0 {
			availableRows = util.DeleteItem(queenRow, availableRows)
		}
	}
	return availableRows
}

func getFirstEmptyColumn(queens []int) int {
	for i,queen := range queens {
		if queen < 0 {
			return i
		}
	}
	return -1
}
