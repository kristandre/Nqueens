package main

import (
  "math/rand"
  "time"
  "sync"
  "../util"
  "../solution"
  "fmt"
  "math"
  "strings"
  "sort"
  "strconv"
)

type Pair struct {
  x, y int
}

func (p *Pair) getKey() string {
  return fmt.Sprintf("%s,%s", math.Min(float64(p.x), float64(p.y)), math.Max(float64(p.x), float64(p.y)))
}

const (
  runTime = 10
  workerCount = 10
  tabuNum = 10
  seldomSwapProb = 0.0075
)

var (
  size = 0
  solutions = solution.SolutionTree{Queen: -1}

  wg sync.WaitGroup
  mutex sync.Mutex
)

func addSolution(s []int) {
  mutex.Lock()
  queens := s
  solutions.AddQueens(queens)
  for i := 0; i < 3; i++ {
    queens = util.Rotate90(queens)
    solutions.AddQueens(queens)
  }
  mutex.Unlock()
}

func initializeTabuList(size int) map[string]int {
  tabuMap := make(map[string]int)
  for i := 0; i < size - 1; i++ {
    for j := i + 1; j < size; j++ {
      pair := Pair{x:i, y:j}
      tabuMap[pair.getKey()] = 0
    }
  }
  return tabuMap
}

func decTabuList(tabuList map[string]int) {
  for key, tabuNum := range tabuList {
    if tabuNum > 0 {
      tabuList[key]--
    }
  }
}

func createSwapList(queens []int, row int) map[int][]Pair {
  swapMap := make(map[int][]Pair)
  for _, queenRow := range queens {
    if queenRow == row {
      continue
    }
    pair := Pair{x:queenRow, y:row}
    swapVal := evaluateSwap(queens, pair)
    _, inMap := swapMap[swapVal]
    if !inMap {
      swapMap[swapVal] = []Pair{}
    }
    swapMap[swapVal] = append(swapMap[swapVal], pair)
  }
  return swapMap
}

func evaluateSwap(queens []int, pair Pair) int {
  crashesBefore := util.EvaluateBoard(queens)
  queensCopy := make([]int, size)
  copy(queensCopy, queens)
  queensCopy = swap(queensCopy, pair)
  crashesAfter := util.EvaluateBoard(queensCopy)
  if crashesAfter == 0 {
    addSolution(queensCopy)
  }
  return crashesAfter - crashesBefore
}

func getFewestSwapsPair(longTermMem, tabuList map[string]int) Pair {
  fewestSwapsCount := 999999999999
  bestPair := Pair{x:-1}
  for key, swapCount := range longTermMem {
    keys := strings.Split(key, ",")
    x, _ := strconv.Atoi(keys[0])
    y, _ := strconv.Atoi(keys[1])
    pair := Pair{x:x, y:y}
    if bestPair.x == -1 {
      bestPair = pair
      continue
    }
    if tabuList[pair.getKey()] > 0 {
      continue
    }
    if swapCount < fewestSwapsCount {
      fewestSwapsCount = swapCount
      bestPair = pair
    }
  }
  return bestPair
}

func swap(queens []int, pair Pair) []int {
  queen1Index := util.IndexOf(pair.x, queens)
  queen2Index := util.IndexOf(pair.y, queens)
  queens[queen1Index] = pair.y
  queens[queen2Index] = pair.x
  return queens
}

func tabuSearch(queens []int) {
  defer wg.Done()
  tabuList := initializeTabuList(size)
  longTermMem := initializeTabuList(size)
  startTime := time.Now()
  nextQueen := 0
  tabuQueens := queens
  for time.Since(startTime) / 1000000000 < runTime {
    if rand.Float32() < seldomSwapProb {
      // Swap a seldom selected pair
      fewestPair := getFewestSwapsPair(longTermMem, tabuList)
      tabuQueens = swap(tabuQueens, fewestPair)
      decTabuList(tabuList)
      tabuList[fewestPair.getKey()] = tabuNum
      longTermMem[fewestPair.getKey()]++
      continue
    }
    swapList := createSwapList(tabuQueens, nextQueen)
    values := getSwapMapKeys(swapList)
    for _, val := range values {
      queenFound := false
      pairs := swapList[val]
      for _, pair := range pairs {
        if tabuList[pair.getKey()] > 0 {
          continue
        }
        q := swap(tabuQueens, pair)
        tabuQueens = q
        decTabuList(tabuList)
        tabuList[pair.getKey()] = tabuNum
        longTermMem[pair.getKey()]++
        nextQueen = pair.x
        queenFound = true
        break
      }
      if queenFound {
        break
      }
    }
  }
}

func getSwapMapKeys(m map[int][]Pair) []int {
  keys := []int{}
  for key := range m {
    keys = append(keys, key)
  }
  sort.Ints(keys)
  return keys
}

func main() {
  start := time.Now()
  rand.Seed(time.Now().UnixNano())
  queens := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29}
  size = len(queens)
  queens = util.PrepareBoard(queens)
  wg.Add(workerCount)
  for i := 0; i < workerCount; i++ {
    q := make([]int, len(queens))
    copy(q, queens)
    go tabuSearch(q)
  }
  wg.Wait()
  elapsed := time.Since(start)
  fmt.Println(size, "queens")
  fmt.Println(solutions.Size, "solutions")
  fmt.Println("Execution time:", elapsed)
}
