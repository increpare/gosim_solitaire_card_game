package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	//"runtime/debug"
	"time"
)

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

var cards = [][]int{
	{0, 0, 0, 0}, //0
	{1, 0, 0, 0}, //1
	{0, 1, 0, 0}, //2
	{0, 0, 1, 0}, //3
	{1, 1, 0, 0}, //4
	{1, 0, 1, 0}, //5
	{0, 1, 1, 0}, //6
	{0, 0, 1, 1}, //7
	{1, 1, 1, 0}, //8
	{1, 0, 1, 1}, //9
	{0, 1, 1, 1}, //10
	{1, 1, 1, 1}, //11
}

var compatibility [][]int

func compatible(x []int, y []int) int {
	if (x[0] <= y[0] && x[1] <= y[1] && x[2] <= y[2] && x[3] <= y[3]) ||
		(x[0] >= y[0] && x[1] >= y[1] && x[2] >= y[2] && x[3] >= y[3]) {
		return 1
	} else {
		return 0
	}
}

func genCompatibilityMatrix() {
	var l = len(cards)
	compatibility = make([][]int, l)
	for i := 0; i < l; i++ {
		compatibility[i] = make([]int, l)
	}

	for i := 0; i < l; i++ {
		compatibility[i][i] = 1
		for j := i + 1; j < l; j++ {
			compatibility[i][j] = compatible(cards[i], cards[j])
			compatibility[j][i] = compatibility[i][j]
		}
	}
}

type State struct {
	top      int
	left     int
	right    int
	playable []int
}

func genStartingDecks() []State {
	l := len(cards)
	result := make([]State, l*(l-1)/2)
	counter := 0
	for i := 0; i < l; i++ {
		for j := i + 1; j < l; j++ {
			var s State
			s.top = -1
			s.left = i
			s.right = j
			s.playable = make([]int, l-2)
			for k := 0; k < l-2; k++ {
				s.playable[k] = k
				if k >= j-1 {
					s.playable[k] += 2
				} else if k >= i {
					s.playable[k]++
				}
			}
			result[counter] = s
			counter++
		}
	}
	return result
}

func dupstate(s State) State {
	var result State
	result.top = s.top
	result.left = s.left
	result.right = s.right
	result.playable = make([]int, len(s.playable))
	copy(result.playable, s.playable)
	return result
}

func isSolid(s State) bool {
	if s.left >= 0 && s.right >= 0 {
		return true
	}

	if len(s.playable) == 0 {
		return true
	}

	return false
}

func solidify(s State) []State {
	if len(s.playable) == 0 {
		var result = make([]State, 1)
		result[0] = s
		return result
	}

	if s.left == -1 && s.right == -1 {
		return genStartingDecks()
	}

	var result []State = make([]State, len(s.playable))

	for i := 0; i < len(s.playable); i++ {
		var t State
		toplay := s.playable[i]
		t.top = s.top
		t.playable = make([]int, len(s.playable)-1)
		for j := 0; j < i; j++ {
			t.playable[j] = s.playable[j]
		}
		for j := i + 1; j < len(s.playable); j++ {
			t.playable[j-1] = s.playable[j]
		}
		if s.left == -1 {
			t.left = toplay
			t.right = s.right
		} else {
			t.left = s.left
			t.right = toplay
		}
		result[i] = t

	}
	//fmt.Printf("solidify:\n%+v\nto:\n%+v\n", s, result)
	return result
}

func initState() State {
	var result State
	result.top = -1
	result.left = -1
	result.right = -1
	result.playable = make([]int, len(cards))
	for i := 0; i < len(cards); i++ {
		result.playable[i] = i
	}
	return result
}

func canPlayLeft(s State) bool {
	if s.top == -1 {
		return true
	}

	if s.left == -1 {
		return false
	}

	if compatibility[s.left][s.top] == 1 {
		return true
	} else {
		return false
	}
}

func canPlayRight(s State) bool {
	if s.top == -1 {
		return true
	}

	if s.right == -1 {
		return false
	}

	if compatibility[s.right][s.top] == 1 {
		return true
	} else {
		return false
	}
}

func playLeft(s State) State {
	s.top = s.left
	s.left = -1
	return s
}

func playRight(s State) State {
	s.top = s.right
	s.right = -1
	return s
}

func E(s State) float64 {
	if s.left == -1 && s.right == -1 {
		return 1
	}

	var leftScore float64 = 0
	if canPlayLeft(s) {
		var t = playLeft(s)
		var possibilities []State = solidify(t)
		var l = len(possibilities)
		//fmt.Printf("%v -> %v\n", len(s.playable), len(possibilities[0].playable))
		for i := 0; i < l; i++ {
			var p = possibilities[i]
			leftScore += E(p) / float64(l)
		}
	}

	var rightScore float64 = 0
	if canPlayRight(s) {
		var t = playRight(s)
		var possibilities []State = solidify(t)
		var l = len(possibilities)
		for i := 0; i < l; i++ {
			var p = possibilities[i]
			rightScore += E(p) / float64(l)
		}
	}

	return math.Max(leftScore, rightScore)
}

func doSim() {
	var init = initState()
	var startingStates = solidify(init)
	var probs []float64 = make([]float64, len(startingStates))

	for i := 0; i < len(startingStates); i++ {
		probs[i] = E(startingStates[i])
		var left = startingStates[i].left
		var right = startingStates[i].right
		fmt.Printf("(%v,%v) -> %.1f%%\n", left, right, 100*probs[i])
	}

	var average float64 = 0
	for i := 0; i < len(probs); i++ {
		average += probs[i] / float64(len(probs))
	}
	fmt.Printf("total average = %.1f%%\n", 100*average)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	rand.Seed(time.Now().UTC().UnixNano())
	genCompatibilityMatrix()
	doSim()
}
