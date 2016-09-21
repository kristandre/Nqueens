package solution

import (
	"fmt"
	"strings"
)

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

func (solutions *SolutionTree) ToString() string {
	s := ""
	if solutions.Queen >= 0 {
		s += fmt.Sprintf("%d ", solutions.Queen + 1)
	}
	if len(solutions.Children) > 0 {
		for _, child := range solutions.Children {
			s += child.ToString()
		}
	} else {
		s += fmt.Sprintln("")
	}
	return strings.TrimSpace(s)
}
