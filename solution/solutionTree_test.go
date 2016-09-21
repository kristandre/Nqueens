package solution

import (
	"testing"
	"../util"
	"fmt"
)

func TestSolutionTree_AddQueens(t *testing.T) {
	solution := SolutionTree{Queen:-1}
	s := []int{3, 6, 2, 7, 1, 4, 0, 5}
	solution.AddQueens(s)

	nextSolution := solution
	for _, queen := range s {
		if len(nextSolution.Children) != 1 {
			t.Fatal("no children found, expected 1")
		}
		if nextSolution.Children[queen].Queen != queen {
			t.Fatalf("expected queen %d, got queen %d", queen, solution.Children[queen].Queen)
		}
		nextSolution = nextSolution.Children[queen]
	}
}

func TestSolutionTree_AddQueens2 (t *testing.T) {
	solution := SolutionTree{Queen:-1}
	s := []int{3, 6, 2, 7, 1, 4, 0, 5}
	solution.AddQueens(s)
	s2 := util.Rotate90(s)
	solution.AddQueens(s2)
	s3 := util.Rotate90(s2)
	solution.AddQueens(s3)
	fmt.Println(s, s2, s3)
	if len(solution.Children) != 2 {
		t.Fatalf("exptected 2 children, got %d", len(solution.Children))
	}
	if len(solution.Children[3].Children) != 2 {
		t.Fatalf("exptected 2 children, got %d", len(solution.Children))
	}

	if solution.Size != 3 {
		t.Fatalf("exptected size 3, got size %d", solution.Size)
	}
}

func TestSolutionTree_AddQueens3(t *testing.T) {
	solution := SolutionTree{Queen:-1}
	s := []int{3, 6, 2, 7, 1, 4, 0, 5}
	s2 := []int{3, 6, 2, 7, 1, 4, 0, 5}
	solution.AddQueens(s)
	solution.AddQueens(s2)

	nextSolution := solution
	for _, queen := range s {
		if len(nextSolution.Children) != 1 {
			t.Fatal("no children found, expected 1")
		}
		if nextSolution.Children[queen].Queen != queen {
			t.Fatalf("expected queen %d, got queen %d", queen, solution.Children[queen].Queen)
		}
		nextSolution = nextSolution.Children[queen]
	}
}

func TestSolutionTree_AddQueens4(t *testing.T) {
	solution := SolutionTree{Queen:-1}
	s := []int{3, 6, 2, 7, 1, 4, 0, 5}
	s2 := []int{3, 6, 2, 7, 1, 4, 0, 5}
	s3 := []int{3, 6, 2, 7, 1, 4, 5, 0}
	solution.AddQueens(s)
	solution.AddQueens(s2)
	solution.AddQueens(s3)

	if len(solution.Children) != 1 {
		t.Fatalf("exptected 1 children, got %d", len(solution.Children))
	}

	if solution.Size != 2 {
		t.Fatalf("exptected size 2, got size %d", solution.Size)
	}
}

func TestSolutionTree_ToString(t *testing.T) {
	solution := SolutionTree{Queen:-1}
	s := []int{3, 6, 2, 7, 1, 4, 0, 5}
	solution.AddQueens(s)
	expected := "4 7 3 8 2 5 1 6"
	actual := solution.ToString()
	if actual != expected {
		t.Fatalf("expected %s, got %s", expected, actual)
	}
}

func TestSolutionTree_ToString2(t *testing.T) {
	solution := SolutionTree{Queen:-1}
	s := []int{3, 6, 2, 7, 1, 4, 0, 5}
	s4 := []int{3, 6, 2, 7, 1, 4, 5, 0}
	solution.AddQueens(s)
	solution.AddQueens(s4)
	s2 := util.Rotate90(s)
	solution.AddQueens(s2)
	s3 := util.Rotate90(s2)
	solution.AddQueens(s3)
	expected := `4 7 3 8 2 5 1 6
	3 2 8 6 1 3 5 7
	3 8 4 7 1 6 2 4`
	actual := solution.ToString()
	if actual != expected {
		t.Fatalf("expected %s, got %s", expected, actual)
	}
}
