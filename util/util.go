package util

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
)

func IndexOf(needle int, haystack []int) int {
	for i, item := range haystack {
		if item == needle {
			return i
		}
	}
	return -1
}

func SliceInArray(slice []int, arr [][]int) bool {
	if slice == nil || len(slice) == 0 {
		return false
	}
	for _, a := range arr {
		inArr := true
		for i, b := range slice {
			if len(a) < len(slice) {
				inArr = false
				continue
			}
			if a[i] != b {
				inArr = false
				continue
			}
		}

		if inArr {
			return true
		}
	}
	return false
}

func DeleteItem(item int, arr []int) []int {
	index := IndexOf(item, arr)
	return append(arr[:index], arr[index + 1:]...)
}

func EvaluateBoard(queens []int) int {
	crashCount := 0
	for i, queenRow := range queens {
		for j := i + 1; j < len(queens); j++ {
			queenNextRow := queens[j]
			k := j - i
			if queenRow + k == queenNextRow || queenRow - k == queenNextRow {
				crashCount++
			}
		}
	}
	return crashCount
}

func Rotate90(solution []int) []int {
	rotatedSolution := make([]int, len(solution))
	for i, n := range solution {
		rotatedSolution[len(solution) - n - 1] = i
	}
	return rotatedSolution
}

func QueensEqual(qa, qb []int) bool {
	for i, a := range qa {
		if qb[i] != a {
			return false
		}
	}
	return true
}

func GetAvailableRows(size int) []int {
	rows := make([]int, size)
	for i := 0; i < size; i++ {
		rows[i] = i
	}
	return rows
}


func PrepareBoard(queens []int) []int {
	availableRows := GetAvailableRows(len(queens))
	crashList := []int{}
	for i, queenRow := range queens {
		if IndexOf(queenRow, availableRows) >= 0 {
			DeleteItem(queenRow, availableRows)
		} else {
			crashList = append(crashList, i)
		}
	}
	for _, col := range crashList {
		queens[col] = availableRows[0]
		DeleteItem(availableRows[0], availableRows)
	}
	return queens
}

func GetQueensFromInput() []int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Insert Size: ")
	_, err := reader.ReadString('\n')

	fmt.Print("Insert Queens: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	queensStrList := strings.Split(strings.TrimSpace(text), " ")
	queens :=[]int{}

	for _, queenStr := range queensStrList {
		queen, _ := strconv.Atoi(queenStr)
		queens = append(queens, queen - 1)
	}

	return PrepareBoard(queens)
}

func PrintSolutionsToFile(solutions [][]int, filename string) {
	fo, _ := os.Create(filename)
	defer fo.Close()
	w := bufio.NewWriter(fo)
	w.WriteString("Number of solutions: " + strconv.Itoa(len(solutions)) + "\n\n")
	for _, queens := range solutions {
		for i, q := range queens {
			w.WriteString(strconv.Itoa(q + 1))
			if i < len(queens) - 1 {
				w.WriteString(" ")
			}
		}
		w.WriteString("\n")
	}
	w.Flush()
}
