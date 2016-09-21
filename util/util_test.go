package util

import (
	"testing"
	"reflect"
)

var testPopulation = [][]int{
	{0, 1, 2, 3, 4, 5, 6, 7},
	{0, 1, 6, 2, 4, 2, 6, 7},
	{7, 1, 2, 4, 3, 5, 6, 0},
	{7, 6, 5, 4, 3, 2, 1, 0},
}

func TestIndexOf(t *testing.T) {
	var queens = []int{0, 1, 2, 3, 4, 5, 6, 7}
	for i := 0; i < len(queens); i++ {
		runIndexOfTest(i, i, queens, t)
	}
	runIndexOfTest(8, -1, queens, t)
	runIndexOfTest(-1, -1, queens, t)
}

func TestSliceInArray(t *testing.T) {
	var queens = []int{0, 1, 2, 3, 4, 5, 6, 7}
	if !SliceInArray(queens, testPopulation) {
		t.Fatal("1expected TRUE, got FALSE")
	}
	if !SliceInArray([]int{7, 6, 5, 4, 3, 2, 1, 0}, testPopulation) {
		t.Fatal("2expected TRUE, got FALSE")
	}
	if SliceInArray([]int{1, 2, 3}, testPopulation) {
		t.Fatal("3expected FALSE, got TRUE")
	}
	if SliceInArray([]int{1, 2, 3, 4, 5, 6, 7, 8}, testPopulation) {
		t.Fatal("4expected FALSE, got TRUE")
	}
	if SliceInArray([]int{}, testPopulation) {
		t.Fatal("5expected FALSE, got TRUE")
	}
}

func TestDeleteItem(t *testing.T) {
	var queens = []int{0, 1, 2, 3, 4, 5, 6, 7}
	queenToDelete := 3
	deletedQueens := DeleteItem(queenToDelete, queens)
	if len(queens) != 8 {
		t.Fatal("Manipulated input!")
	}
	if len(deletedQueens) != 7 {
		t.Fatal("Exptected length 7, but got something else")
	}
	for _, q := range deletedQueens {
		if q == queenToDelete {
			t.Fatalf("queen %d is still in list", queenToDelete)
		}
	}
}

func TestEvaluateBoard(t *testing.T) {
	runEvaluateBoardTest([]int{0, 1, 2, 3, 4, 5, 6, 7}, 28, t)
	runEvaluateBoardTest([]int{7, 6, 5, 4, 3, 2, 1, 0}, 28, t)
	runEvaluateBoardTest([]int{0, 2, 4, 6, 3, 5, 7, 1}, 2, t)
	runEvaluateBoardTest([]int{1, 3, 5, 7, 2, 0, 6, 4}, 0, t)
}

func TestRotate90(t *testing.T) {
	rotatedQueens := Rotate90([]int{0, 1, 2, 3, 4, 5, 6, 7})
	expected := []int{7, 6, 5, 4, 3, 2, 1, 0}
	if !reflect.DeepEqual(rotatedQueens, expected) {
		t.Fatalf("expected %v, got %v", expected, rotatedQueens)
	}

	rotatedQueens = Rotate90([]int{1, 3, 5, 7, 2, 0, 6, 4})
	expected = []int{3, 6, 2, 7, 1, 4, 0, 5}
	if !reflect.DeepEqual(rotatedQueens, expected) {
		t.Fatalf("expected %v, got %v", expected, rotatedQueens)
	}
}

func TestQueensEqual(t *testing.T) {
	q1 := []int{3, 6, 2, 7, 1, 4, 0, 5}
	q2 := []int{3, 6, 2, 7, 1, 4, 0, 5}
	equal := QueensEqual(q1, q2)
	if !equal {
		t.Fatalf("expected TRUE, but got FALSE, %v, %v", q1, q2)
	}

	q3 := []int{6, 3, 2, 7, 1, 4, 0, 5}
	equal = QueensEqual(q1, q3)
	if equal {
		t.Fatalf("expected FALSE, but got TRUE, %v, %v", q1, q2)
	}
}

func TestGetAvailableRows(t *testing.T) {
	input := 8
	expected := []int{0, 1, 2, 3, 4, 5, 6, 7}
	actual := GetAvailableRows(input)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func runIndexOfTest(input, expectedOutput int, queens []int, t *testing.T) {
	index := IndexOf(input, queens)
	if index != expectedOutput {
		t.Fatalf("expected %d, got %d.", expectedOutput, index)
	}
}

func runEvaluateBoardTest(input []int, expectedOutput int, t *testing.T) {
	crashCount := EvaluateBoard(input)
	if crashCount != expectedOutput {
		t.Fatalf("expected %d crashes, but got %d crashes", expectedOutput, crashCount)
	}
}