package solution

import (
	//"fmt"
	//"strings"
)
//import "fmt"

type SolutionTree struct {
	Children map[int]SolutionTree
	Queen    int
	Size     int
}

func (s *SolutionTree) AddQueens(queens []int) bool {
	if len(s.Children) == 0 {
		s.Children = make(map[int]SolutionTree)
	}
	nextQueen := queens[0]
	n := len(queens)
	isNewSolution := false
	_, ok := s.Children[nextQueen]
	if ok && n > 1 {
		nextSolution := s.Children[nextQueen]
		if nextSolution.AddQueens(queens[1:]) {
			s.Size = s.Size + 1
			return true
		} else {
			return false
		}
	}else if !ok {

		child := SolutionTree{Queen: nextQueen}
		if n > 1 {
			isNewSolution = child.AddQueens(queens[1:])
			if isNewSolution {
				s.Size = s.Size + 1
			}
		} else if n == 1 {
			isNewSolution = true
		}
		s.Children[nextQueen] = child
	}
	return isNewSolution
}

func (solutions *SolutionTree) GetSolutions() [][]int {
	allSolutions := [][]int{}
	var me int
	me = solutions.Queen
	if len(solutions.Children) > 0 {
		for _, child := range solutions.Children {
			childSolutions := child.GetSolutions()
			for _, childSolution := range childSolutions {
				x := []int{}
				if (me >= 0) {
					x = append(x, me)
				}
				x = append(x, childSolution...)
				allSolutions = append(allSolutions, x)
			}
		}
	} else {
		h := []int{solutions.Queen}
		allSolutions = append(allSolutions, h)
	}
	return allSolutions
}
